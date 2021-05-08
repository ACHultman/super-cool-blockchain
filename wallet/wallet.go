package wallet

import (
	"blockchain/blockchain"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey // elliptical curve digital-signing algorithm
	PublicKey  []byte
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}

// NEwKeyPair generate new private key - public key pair
func NEwKeyPair() (ecdsa.PrivateKey, []byte) {
	// P256 returns a Curve which implements NIST P-256
	curve := elliptic.P256()

	// generates a public and private key pair.
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	blockchain.HandleError(err)

	// combine X, Y components of public key
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

// MakeWallet make wallet
func MakeWallet() *Wallet {
	private, public := NEwKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// PublicKeyHash hash public key with sha256 and ripemd160
func PublicKeyHash(pubKey []byte) []byte {
	// get sha256 hashed public key
	pubHash := sha256.Sum256(pubKey)

	// get ripemd160 hash function
	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	blockchain.HandleError(err)

	// get ripemd160 hashed sha256 pubkley hash
	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func Checksum(payload []byte) []byte {
	// hash twice
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}
