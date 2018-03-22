package shamir

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"

	"github.com/republicprotocol/republic-go/stackint"
)

// A Share struct represents some share of a secret after the secret has been
// encoded.
type Share struct {
	Key   int64
	Value stackint.Int1024
	// ValueBig *big.Int
}

// Shares are a slice of Share structs.
type Shares []Share

// Split a secret into Shares. N represents the number of Shares that the
// secret will be split into, and K represents the number of Share required to
// reconstruct the secret. Prime is used to define the finite field from which
// secrets can be selected. A slice of Shares, or an error, is returned.
func Split(n int64, k int64, prime, secret *stackint.Int1024) (Shares, error) {
	// Validate the encoding by checking that N is greater than K, and that the
	// secret is within the finite field.
	if n < k {
		return nil, NewNKError(n, k)
	}
	if prime.LessThanOrEqual(secret) {
		return nil, NewFiniteFieldError(secret)
	}

	mont := stackint.PrimeM

	// Generate K polynomial coefficients, where the first coefficient is the
	// secret.
	one := stackint.One()
	max := prime.Sub(&one)
	coefficients := make([]*stackint.MontInt, k)
	secretM := mont.ToMont(secret)
	coefficients[0] = &secretM

	for i := int64(1); i < k; i++ {
		coefficient, err := stackint.Random(rand.Reader, &max)
		if err != nil {
			return nil, err
		}
		coefficentM := mont.ToMont(&coefficient)
		coefficients[i] = &coefficentM
	}

	// Create N shares.
	shares := make(Shares, n)
	for x := int64(1); x <= n; x++ {

		// accum := coefficients[0]
		accumM := (*coefficients[0]).MontClone()
		// base := x
		// base := stackint.FromUint64(uint64(x))
		baseM := mont.FromUint64(uint64(x))
		// expMod := base % prime
		expM := baseM.MontClone()
		// expMod := exp.Mod(prime)

		// accumBig.Set(coefficientsBig[0])
		// baseBig.SetInt64(x)
		// expBig.Set(baseBig)
		// expModBig.Mod(expBig, primeBig)

		// ShouldEqual(accumBig, accum)
		// ShouldEqual(baseBig, base)
		// ShouldEqual(expBig, exp)
		// ShouldEqual(expModBig, expMod)

		// Evaluate the polyomial at x.
		for j := range coefficients[1:] {
			// ShouldEqual(accumBig, accum)

			// co := (coefficients * expoMod) % prime
			coefficient := coefficients[j].MontClone()
			coM := coefficient.MontMul(&expM)

			// accum = (accum + co) % prime
			accumM = accumM.MontAdd(&coM)

			expM = expM.MontMul(&baseM)
			// expMod = (expMod * base ) % prime
			// exp = exp.Mul(&base)
			// expBig.Mul(expBig, baseBig)
			// expMod = exp.Mod(prime)
			// expModBig.Mod(expBig, primeBig)

			// ShouldEqual(accumBig, accum)
		}
		// ShouldEqual(accumBig, accum)
		shares[x-1] = Share{
			Key:   x,
			Value: accumM.ToInt1024(),
			// ValueBig: big.NewInt(0).Set(accumBig),
		}
	}
	return shares, nil
}

// Join Shares into a secret. Prime is used to define the finite field from
// which the secret was selected. The reconstructed secret, or an error, is
// returned.
func Join(prime *stackint.Int1024 /*, primeBig *big.Int */, shares Shares) *stackint.Int1024 /*big.Int*/ {
	// secretBig := big.NewInt(0)

	// Setup big numbers so that we do not have to keep recreating them in each
	// loop.
	// valueBig := big.NewInt(0)
	// numBig := big.NewInt(1)
	// denBig := big.NewInt(1)
	// startBig := big.NewInt(0)
	// nextBig := big.NewInt(0)
	// nextNegBig := big.NewInt(0)
	// nextDiffBig := big.NewInt(0)

	mont := stackint.PrimeM

	secret := mont.FromUint64(0)

	primeM := mont.M

	// Compute the Lagrange basic polynomial interpolation.
	for i := 0; i < len(shares); i++ {
		num := stackint.OneC
		den := stackint.OneC

		// numBig.SetInt64(1)
		// denBig.SetInt64(1)

		for j := 0; j < len(shares); j++ {
			if i == j {
				continue
			}
			// startposition = shares[formula][0];
			// startBig.SetInt64(int64(shares[i].Key))
			start := mont.FromUint64(uint64(shares[i].Key))

			// nextposition = shares[count][0];
			// nextBig.SetInt64(int64(shares[j].Key))
			next := mont.FromUint64(uint64(shares[j].Key))

			// numerator = (numerator * -nextposition) % prime;
			// nextNegBig.SetInt64(0)
			// nextNegBig.Sub(nextNegBig, nextBig)
			// numBig.Mul(numBig, nextNegBig)
			// numBig.Mod(numBig, primeBig)

			nextGen := num.MontMul(&next)
			num = primeM.MontSub(&nextGen)

			// denominator = (denominator * (startposition - nextposition)) % prime;
			// nextDiffBig.Sub(startBig, nextBig)
			// denBig.Mul(denBig, nextDiffBig)
			// denBig.Mod(denBig, primeBig)

			// oldDen := den
			nextDiff := start.MontSub(&next)
			den = den.MontMul(&nextDiff)

			// if denBig.String() != den.String() {
			// 	fmt.Println("!!!")
			// 	fmt.Printf("start: %s\n", start.String())
			// 	fmt.Printf("next: %s\n", next.String())
			// 	fmt.Printf("prime: %s\n", prime.String())
			// 	fmt.Printf("nextDiff: %s\n", nextDiff.String())
			// 	fmt.Printf("oldDen: %s\n", oldDen.String())
			// 	fmt.Printf("den: %s\n", den.String())
			// 	fmt.Printf("denBig: %s\n", denBig.String())
			// 	fmt.Println("!!!")
			// }

			// ShouldEqual(denBig, den)
		}

		// valueBig = shares[formula][1]
		// accumBig = (primeBig + accumBig + (valueBig * numeratorBig * modInverse(denominatorBig))) % primeBig
		// denBig.ModInverse(denBig, primeBig)
		// valueBig.Mul(shares[i].ValueBig, numBig)
		// valueBig.Mul(valueBig, denBig)
		// secretBig.Add(secretBig, valueBig)
		// secretBig.Mod(secretBig, primeBig)

		den = den.MontInv()
		value := mont.ToMont(&shares[i].Value)
		value = value.MontMul(&num)
		value = value.MontMul(&den)
		// value := shares[i].Value.MulModulo(&num, prime)
		// value = value.MulModulo(&den, prime)
		secret = secret.MontAdd(&value)
		// secret = secret.AddModulo(&value, prime)
	}

	secretInt := secret.ToInt1024()
	return &secretInt //, secretBig
}

// ToBytes encodes the Share into a slice of bytes.
func ToBytes(share Share) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, share.Key)
	binary.Write(buf, binary.LittleEndian, share.Value.ToBytes())
	return buf.Bytes()
}

// FromBytes decodes a slice of bytes into a Share.
func FromBytes(bs []byte) (Share, error) {
	key := int64(0)
	buf := bytes.NewReader(bs)
	if err := binary.Read(buf, binary.LittleEndian, &key); err != nil {
		return Share{}, err
	}
	data := make([]byte, buf.Len())
	if err := binary.Read(buf, binary.LittleEndian, data); err != nil {
		return Share{}, err
	}
	return Share{
		Key:   key,
		Value: stackint.FromBytes(data),
		// ValueBig: big.NewInt(0).SetBytes(data),
	}, nil
}
