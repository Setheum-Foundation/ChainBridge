// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package keystore

import (
	"fmt"

	"github.com/ChainSafe/ChainBridgeV2/crypto"
	"github.com/ChainSafe/ChainBridgeV2/crypto/secp256k1"
	"github.com/ChainSafe/ChainBridgeV2/crypto/sr25519"
)

// The Constant "keys". These are the name that the keys are based on. This can be expanded, but
// any additions must be added to TestKeyRing and to insecureKeyFromAddress
const AliceKey = "alice"
const BobKey = "bob"
const CharlieKey = "charlie"
const DaveKey = "dave"
const EveKey = "eve"

// The Chain type Constants
const EthChain = "ethereum"
const SubChain = "substrate"

var TestKeyRing *TestKeyRingHolder

//var TestKeyStoreMap map[string]*Keystore

// TestKeyStore is a struct that holds a Keystore of all the test keys
type TestKeyRingHolder struct {
	EthereumKeys   map[string]*secp256k1.Keypair
	CentrifugeKeys map[string]*sr25519.Keypair
}

// KeyRing holds the keypair related to a specfic keypair type
type KeyRing map[string]crypto.Keypair

// Init function to create a keyRing that can be accessed anywhere without having to recreate the data
func init() {
	TestKeyRing = &TestKeyRingHolder{
		EthereumKeys:   makeETHRing(createKeyRing(EthChain)),
		CentrifugeKeys: makeSUBRing(createKeyRing(SubChain)),
	}

}

func makeETHRing(k KeyRing) map[string]*secp256k1.Keypair {
	ring := map[string]*secp256k1.Keypair{}
	for key, pair := range k {
		ring[key] = pair.(*secp256k1.Keypair)
	}

	return ring
}

func makeSUBRing(k KeyRing) map[string]*sr25519.Keypair {
	ring := map[string]*sr25519.Keypair{}
	for key, pair := range k {
		ring[key] = pair.(*sr25519.Keypair)
	}

	return ring
}

// padWithZeros adds on extra 0 bytes to make a byte array of a specified length
func padWithZeros(key []byte, targetLength int) []byte {
	res := make([]byte, targetLength-len(key))
	return append(res, key...)
}

// errorWrap is a helper function that panics on errors, to make the code cleaner
func errorWrap(in interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return in
}

// createKeyRing creates a KeyRing for the specfied chain/key type
func createKeyRing(chain string) KeyRing {
	ring := map[string]crypto.Keypair{
		AliceKey:   createKeypair(AliceKey, chain),
		BobKey:     createKeypair(BobKey, chain),
		CharlieKey: createKeypair(CharlieKey, chain),
		DaveKey:    createKeypair(DaveKey, chain),
		EveKey:     createKeypair(EveKey, chain),
	}

	return ring

}

// createKeypair creates keypairs based on the private key seed inputted for the specfied chain
func createKeypair(key, chain string) crypto.Keypair {
	switch chain {
	case EthChain:
		bz := padWithZeros([]byte(key), secp256k1.PrivateKeyLength)
		return errorWrap(secp256k1.NewKeypairFromPrivateKey(bz)).(*secp256k1.Keypair)
	case SubChain:
		return errorWrap(sr25519.NewKeypairFromSeed("//" + key)).(*sr25519.Keypair)
	}
	return nil

}

// insecureKeypairFromAddress is used for resolving addresses to test keypairs.
func insecureKeypairFromAddress(key string, chainType string) (crypto.Keypair, error) {
	var kp crypto.Keypair
	var ok bool

	if chainType == EthChain {
		kp, ok = TestKeyRing.EthereumKeys[key]
	} else if chainType == SubChain {
		kp, ok = TestKeyRing.CentrifugeKeys[key]
	} else {
		return nil, fmt.Errorf("unrecognized chain type: %s", chainType)
	}

	if !ok {
		return nil, fmt.Errorf("invalid test key selection: %s", key)
	}

	return kp, nil
}
