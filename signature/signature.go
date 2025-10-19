package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Transaction struct {
	from      string
	to        string
	amount    float64
	signature string
}
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

// creates a digital signature using private key
func signintransaction(tx *Transaction, privatekey *ecdsa.PrivateKey){
	//create a hash of transaction data
	data:= fmt.Sprintf("%s:%s:%f",tx.from,tx.to,tx.amount)
	hash:=sha256.Sum256([]byte(data))
	//generate the digital signature
	r,s,_:=ecdsa.Sign(rand.Reader,privatekey,hash[:])
	signature:=append(r.Bytes(),s.Bytes()...)
	tx.signature=hex.EncodeToString(signature)
}
//checks if transaction signature is valid
func verifytransaction(tx *Transaction,publickey ecdsa.PublicKey)bool{
	data:= fmt.Sprintf("%s:%s:%f",tx.from,tx.to,tx.amount)
	hash:=sha256.Sum256([]byte(data))
	signinbytes,_:=hex.DecodeString(tx.signature)
	r:=big.Int{}
	s:=big.Int{}
	signlen:=len(signinbytes)
	r.SetBytes(signinbytes[:(signlen/2)])
	s.SetBytes(signinbytes[(signlen/2):])
	return ecdsa.Verify(&publickey,hash[:],&r,&s)
}
func main(){
	//create wallets
	sender:=generatewallet()
	receiver:=generatewallet()
	//create a transaction
	tx:=Transaction{
		from: sender.address,
		to: receiver.address,
		amount: 25.5,
	}
	//signin & verify
	signintransaction(&tx,sender.PrivateKey)
	fmt.Println("transaction signed:")
	fmt.Println("from : ",tx.from)
	fmt.Println("to : ",tx.to)
	fmt.Println("amount : ",tx.amount)
	fmt.Println("signature : ",tx.signature)
	//verify signature
	isvalid:=verifytransaction(&tx,sender.PrivateKey.PublicKey)
	fmt.Println("\n signature verified : ",isvalid)
}