package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)
type Transaction struct{
	from string
	to string
	amount float64
	timestamp string
	signature string
}

// Block structure
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}

// Blockchain
type Blockchain struct {
	Blocks         []Block
	TransactionPool []Transaction
	Difficulty     int
}


// Hash calculation
func calchash(block Block)string{
	record:=fmt.Sprintf("%d%s%s%d",block.Index,block.Timestamp,block.PrevHash,block.Nonce)
	for _,tx:= range block.Transactions{
		record+=fmt.Sprintf("%s%s%f%s",tx.from,tx.to,tx.amount,tx.timestamp)
	}
	h:=sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}

// Add Transaction to pool
func (bc *Blockchain) AddTransaction(tx Transaction) {
	bc.TransactionPool = append(bc.TransactionPool, tx)
}

// Mining function
func (bc *Blockchain) MineBlock() {
	if len(bc.TransactionPool) == 0 {
		fmt.Println("No transactions to mine!")
		return
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: bc.TransactionPool,
		PrevHash:     prevBlock.Hash,
		Nonce:        0,
	}

	// Proof of Work
	for {
		hash := calchash(newBlock)
		if hash[:bc.Difficulty] == "0000" {
			newBlock.Hash = hash
			break
		}
		newBlock.Nonce++
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	bc.TransactionPool = []Transaction{} // clear pool
	fmt.Println("Block mined successfully!")
}

// Genesis block
func CreateGenesisBlock() Block {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{},
		PrevHash:     "",
		Nonce:        0,
	}
	genesis.Hash = calchash(genesis)
	return genesis
}

// Initialize blockchain
func InitBlockchain() *Blockchain {
	genesis := CreateGenesisBlock()
	return &Blockchain{
		Blocks:         []Block{genesis},
		Difficulty:     4,
		TransactionPool: []Transaction{},
	}
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
func main() {
	// Initialize blockchain
	blockchain := InitBlockchain()

	// Create wallets
	alice := generatewallet()
	bob := generatewallet()

	// Create transaction
	tx := Transaction{
		from:   alice.address,
		to:     bob.address,
		amount: 10.0,
	}
	signintransaction(&tx, alice.PrivateKey)

	// Add to transaction pool
	blockchain.AddTransaction(tx)

	// Mine block
	blockchain.MineBlock()

	fmt.Println("\nBlockchain State:")
	for _, block := range blockchain.Blocks {
		fmt.Printf("\nBlock #%d\nHash: %s\nPrevHash: %s\nTransactions: %d\n",
			block.Index, block.Hash, block.PrevHash, len(block.Transactions))
	}
}