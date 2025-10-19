package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// wallet structure
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	publickey  []byte
	address    string
}

// create a wallet with key and address
func generatewallet() *Wallet {
	privatekey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubkey := append(privatekey.PublicKey.X.Bytes(), privatekey.PublicKey.Y.Bytes()...)
	//hash public key to generate address
	pubhash := sha256.Sum256(pubkey)
	address := hex.EncodeToString(pubhash[:])[:20]
	return &Wallet{
		PrivateKey: privatekey,
		publickey:  pubkey,
		address:    address,
	}
}
func (w *Wallet) showWallet() {
	fmt.Println("wallet generated : ")
	fmt.Println("private key : ", hex.EncodeToString(w.PrivateKey.D.Bytes()))
	fmt.Println("public key : ", hex.EncodeToString(w.publickey))
	fmt.Println("address : ", w.address)
}
func main(){
	//create two wallet
	wallet1:=generatewallet()
	wallet2:=generatewallet()
	wallet1.showWallet()
	wallet2.showWallet()
}