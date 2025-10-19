//simple block structure code
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	index   int
	data    string
	prehash string
	hash    string
}
//generate hash
func calchash(block Block) string {
	record := string(block.index) + block.data + block.prehash
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
func main(){
	//first genesis block
	genesis:=Block{index: 0,data: "genesis block",prehash: ""}
	genesis.hash=calchash(genesis)
	fmt.Println("block#",genesis.index)
	fmt.Println("data:",genesis.data)
	fmt.Println("hash:",genesis.hash)

}