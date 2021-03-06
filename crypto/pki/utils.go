package pki

import (
	"crypto/sha256"

	"github.com/dedis/kyber"
	"github.com/ethereum/go-ethereum/log"
	"github.com/harmony-one/harmony/crypto"
)

// GetAddressFromPublicKey returns address given a public key.
func GetAddressFromPublicKey(pubKey kyber.Point) [20]byte {
	bytes, err := pubKey.MarshalBinary()
	if err != nil {
		log.Error("Failed to serialize challenge")
	}
	address := [20]byte{}
	hash := sha256.Sum256(bytes)
	copy(address[:], hash[12:])
	return address
}

// GetAddressFromPrivateKey returns address given a private key.
func GetAddressFromPrivateKey(priKey kyber.Scalar) [20]byte {
	return GetAddressFromPublicKey(GetPublicKeyFromScalar(priKey))
}

// GetAddressFromPrivateKeyBytes returns address from private key in bytes.
func GetAddressFromPrivateKeyBytes(priKey [32]byte) [20]byte {
	return GetAddressFromPublicKey(GetPublicKeyFromScalar(crypto.Ed25519Curve.Scalar().SetBytes(priKey[:])))
}

// GetAddressFromInt is the temporary helper function for benchmark use
func GetAddressFromInt(value int) [20]byte {
	return GetAddressFromPublicKey(GetPublicKeyFromScalar(GetPrivateKeyScalarFromInt(value)))
}

// GetPrivateKeyScalarFromInt return private key scalar.
func GetPrivateKeyScalarFromInt(value int) kyber.Scalar {
	return crypto.Ed25519Curve.Scalar().SetInt64(int64(value))
}

// GetPrivateKeyFromInt returns private key in bytes given an interger.
func GetPrivateKeyFromInt(value int) [32]byte {
	priKey, err := crypto.Ed25519Curve.Scalar().SetInt64(int64(value)).MarshalBinary()
	priKeyBytes := [32]byte{}
	if err == nil {
		copy(priKeyBytes[:], priKey[:])
	}
	return priKeyBytes
}

// GetPublicKeyFromPrivateKey return public key from private key.
func GetPublicKeyFromPrivateKey(priKey [32]byte) kyber.Point {
	suite := crypto.Ed25519Curve
	scalar := suite.Scalar()
	scalar.UnmarshalBinary(priKey[:])
	return suite.Point().Mul(scalar, nil)
}

// GetPublicKeyFromScalar is the same as GetPublicKeyFromPrivateKey, but it directly works on kyber.Scalar object.
func GetPublicKeyFromScalar(priKey kyber.Scalar) kyber.Point {
	return crypto.Ed25519Curve.Point().Mul(priKey, nil)
}

// GetBytesFromPublicKey converts public key point to bytes
func GetBytesFromPublicKey(pubKey kyber.Point) [32]byte {
	bytes, err := pubKey.MarshalBinary()
	result := [32]byte{}
	if err == nil {
		copy(result[:], bytes)
	}
	return result
}
