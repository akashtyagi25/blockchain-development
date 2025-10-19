// transaction & block integration
package main

import (
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
}
//block structure
type Block struct{
	index int
	timestamp string
	Transactions []Transaction
	prehash string
	hash string
	nonce int
}
func calchash(block Block)string{
	record:=fmt.Sprintf("%d%s%s%d",block.index,block.timestamp,block.prehash,block.nonce)
	for _,tx:= range block.Transactions{
		record+=fmt.Sprintf("%s%s%f%s",tx.from,tx.to,tx.amount,tx.timestamp)
	}
	h:=sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}
//create block generates a new block from previous & new transaction
func createblock(preblock Block,txs []Transaction)Block{
	newblock:=Block{
		index: preblock.index+1,
		timestamp: time.Now().String(),
		Transactions: txs,
		prehash: preblock.hash,
		nonce: 0,
	}
	newblock.hash=calchash(newblock)
	return newblock
}
//create genesis block create first block
func creategenesis()Block{
	genesistx:=Transaction{"network","akash",100.0,time.Now().String()}
	block:=Block{
		index: 0,
		timestamp: time.Now().String(),
		Transactions: []Transaction{genesistx},
		prehash: "",
	}
	block.hash=calchash(block)
	return block
}
//bllokchain validations
func isblockvalid(newblock,oldblock Block)bool{
	if oldblock.index+1 !=newblock.index{
		return false
	}
	if oldblock.hash !=newblock.prehash{
		return false
	}
	if calchash(newblock) != newblock.hash{
		return false
	}
	return true
}
func main(){
	genesisblock:=creategenesis()
	blockchain:=[]Block{genesisblock}
	//new transactions
	txs:=[]Transaction{
		{"akash","tyagi",30.5,time.Now().String()},
		{"tyagi","rohan",10.5,time.Now().String()},
	}
	//create new block with transactions
	newblock:=createblock(genesisblock,txs)
	//validate&append
	if isblockvalid(newblock,genesisblock){
		blockchain=append(blockchain, newblock)
	}
	//print blockchain
	for _,block:=range blockchain{
		fmt.Printf("\n======block %d ======\n",block.index)
		fmt.Printf("timestamp : %s\n",block.timestamp)
		fmt.Printf("prehash : %s\n",block.prehash)
		fmt.Printf("hash : %s\n",block.hash)
		fmt.Println("transactions:")
		for _,tx:=range block.Transactions{
			fmt.Printf(" %s -> %s : %.2f coins\n",tx.from,tx.to,tx.amount)
		}
	}
}