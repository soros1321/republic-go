package identity_test

import (
	. "github.com/republicprotocol/go-identity"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MultiAddresses with support for Republic Protocol", func() {

	Context("after importing the identity package", func() {

		It("there should be a protocol called republic", func() {
			Ω(ProtocolWithName("republic").Name).Should(Equal("republic"))
		})

		Specify("the Republic Protocol code should be defined with the correct constants", func() {
			Ω(ProtocolWithCode(RepublicCode).Name).Should(Equal("republic"))
			Ω(ProtocolWithCode(RepublicCode).Code).Should(Equal(RepublicCode))
		})
	})
})
