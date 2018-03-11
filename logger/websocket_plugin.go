package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/republicprotocol/go-do"
)

type WebSocketPlugin struct {
	do.GuardedObject

	Srv          *http.Server
	Host         string `json:"host"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	registration string `json:"registration"`

	info  chan interface{}
	error chan Message
	warn  chan interface{}
	usage chan Usage
}

func NewWebSocketPlugin(host, port, username, password string) Plugin {
	plugin := &WebSocketPlugin{
		GuardedObject: do.NewGuardedObject(),
		Host:          host,
		Port:          port,
		Username:      username,
		Password:      password,
		info:          make(chan interface{}, 1),
		error:         make(chan Message, 1),
		warn:          make(chan interface{}, 1),
		usage:         make(chan Usage, 1),
	}
	return plugin
}

func (plugin *WebSocketPlugin) logHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer c.Close()

	go func() {
		for {
			request := &struct {
				Name string `json: "name"`
			}{}
			err := c.ReadJSON(request)
			if err != nil {
				log.Println(err)
				return
			}
			if request.Name == TagRegister {
				if plugin.registration != "" {
					registration := new(Registration)
					err = json.Unmarshal([]byte(plugin.registration), registration)
					if err != nil {
						log.Println(err)
						return
					}
					err := c.WriteJSON(registration)
					if err != nil {
						return
					}
				}
			}
		}
	}()

	// Broadcast errors
	for {
		select {
		case u := <-plugin.usage:
			c.WriteJSON(u)
		case e := <-plugin.error:
			c.WriteJSON(e)
		case i := <-plugin.info:
			c.WriteJSON(i)
		case warning := <-plugin.warn:
			c.WriteJSON(warning)
		default:
			break
		}
	}

	//todo : how to close this
}

func (plugin *WebSocketPlugin) Start() error {
	plugin.Srv = &http.Server{
		Addr: plugin.Host + ":" + plugin.Port,
	}
	http.HandleFunc("/logs", plugin.logHandler)
	go func() {
		log.Println(fmt.Sprintf("WebSocket logger listening on %s:%s", plugin.Host, plugin.Port))
		err := plugin.Srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (plugin *WebSocketPlugin) Stop() error {
	return plugin.Srv.Shutdown(nil)
}

type Message struct {
	Time    string
	Type    string
	Message string
}

func (plugin *WebSocketPlugin) Info(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	event := Event{
		Type: "event",
		Time: time.Now(),
		Data: EventData{
			Tag:     tag,
			Level:   "INFO",
			Message: message,
		},
	}
	if tag == TagRegister {
		plugin.registration = message
		return nil
	}
	if len(plugin.info) == 1 {
		<-plugin.info
	}
	plugin.info <- event

	return nil
}

func (plugin *WebSocketPlugin) Error(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	msg := Message{
		time.Now().Format("2006/01/02 15:04:05 "), tag, message,
	}
	if len(plugin.error) == 1 {
		<-plugin.error
	}
	plugin.error <- msg
	return nil
}

func (plugin *WebSocketPlugin) Warn(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	event := Event{
		Type: "event",
		Time: time.Now(),
		Data: EventData{
			Tag:     tag,
			Level:   "WARN",
			Message: message,
		},
	}
	if len(plugin.warn) == 1 {
		<-plugin.warn
	}
	plugin.warn <- event

	return nil
}

func (plugin *WebSocketPlugin) Usage(cpu float32, memory, network int32) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	usage := Usage{
		Type: "usage",
		Time: time.Now(),
		Data: UsageData{
			Cpu:     cpu,
			Memory:  memory,
			network: network,
		},
	}
	if len(plugin.usage) == 1 {
		<-plugin.usage
	}
	plugin.usage <- usage
	return nil
}