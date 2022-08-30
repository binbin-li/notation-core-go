package signature

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
)

// Algorithm defines the signature algorithm.
type Algorithm int

// Signature algorithms supported by this library.
//
// Reference: https://github.com/notaryproject/notaryproject/blob/main/signature-specification.md#algorithm-selection
const (
	AlgorithmPS256 Algorithm = 1 + iota // RSASSA-PSS with SHA-256
	AlgorithmPS384                      // RSASSA-PSS with SHA-384
	AlgorithmPS512                      // RSASSA-PSS with SHA-512
	AlgorithmES256                      // ECDSA on secp256r1 with SHA-256
	AlgorithmES384                      // ECDSA on secp384r1 with SHA-384
	AlgorithmES512                      // ECDSA on secp521r1 with SHA-512
)

// KeyType defines the key type.
type KeyType int

const (
	KeyTypeRSA KeyType = 1 + iota // KeyType RSA
	KeyTypeEC                     // KeyType EC
)

// KeySpec defines a key type and size.
type KeySpec struct {
	Type KeyType
	Size int
}

// ExtractKeySpec extracts KeySpec from the signing certificate.
func ExtractKeySpec(signingCert *x509.Certificate) (KeySpec, error) {
	switch key := signingCert.PublicKey.(type) {
	case *rsa.PublicKey:
		switch bitSize := key.Size() << 3; bitSize {
		case 2048, 3072, 4096:
			return KeySpec{
				Type: KeyTypeRSA,
				Size: bitSize,
			}, nil
		default:
			return KeySpec{}, &UnsupportedSigningKeyError{
				Msg: fmt.Sprintf("rsa key size %d is not supported", bitSize),
			}
		}
	case *ecdsa.PublicKey:
		switch bitSize := key.Curve.Params().BitSize; bitSize {
		case 256, 384, 521:
			return KeySpec{
				Type: KeyTypeEC,
				Size: bitSize,
			}, nil
		default:
			return KeySpec{}, &UnsupportedSigningKeyError{
				Msg: fmt.Sprintf("ecdsa key size %d is not supported", bitSize),
			}
		}
	}
	return KeySpec{}, &UnsupportedSigningKeyError{
		Msg: "invalid public key type",
	}
}
