package signv2

import (
	"crypto"
)

// KeyAlgorithm is used to map strings used in e.g. config files to implementations.
type KeyAlgorithm string

const (
	RSA KeyAlgorithm = "RSA"
	EC               = "EC"
)

// HashAlgorithm is used to map strings used in e.g. config files to implementations. This is
// partially redundant with crypto.Hash, but its purpose is to be able to basically map a string
// from a config file into a crypto.Hash elsewhere in code
type HashAlgorithm string

const (
	SHA256 HashAlgorithm = "SHA256"
	SHA512               = "SHA512"
)

// AsHash turns our string-based enum type into a Go crypto.Hash value.
func (h HashAlgorithm) AsHash() crypto.Hash {
	switch h {
	case SHA256:
		return crypto.SHA256
	case SHA512:
		return crypto.SHA512
	default:
		// panic is a smidge aggressive here, but we can't return nil and caller shouldn't have called
		// us on a string not listed above. in normal operation this is pretty bad.
		panic("unknown hash algorithm requested")
	}
}

// AlgorithmID labels the Android APK signing scheme v2 magic constants. Note that these constants
// serve the same function as the usual ASN.1 object ID registered constants, but in an integer
// format.
type AlgorithmID uint32
