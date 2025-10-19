// mining proof of work code
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)
type Block struct{
	index int
	data string
	prehash string
	hash string
	nonce int
}
func calchash(block Block) string{
	record:=strconv.Itoa(block.index)+block.data+block.prehash+strconv.Itoa(block.nonce)
	hash:=sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
//mining function
func miningblock(block *Block,difficulty int){
	prefix:=strings.Repeat("0",difficulty)
	for{
		block.hash=calchash(*block)
		if strings.HasPrefix(block.hash,prefix){
			fmt.Println("block mined successfully: ",block.hash)
			break
		}
		block.nonce++
	}
}
func main(){
	difficulty:=5
	genesis:=Block{index: 0,data: "genesis block",prehash: ""}
	miningblock(&genesis,difficulty)
	block2:=Block{index: 1,data: "a sent 5 coins to b",prehash: genesis.hash}
	miningblock(&block2,difficulty)
	fmt.Println("\nblockchain:")
	fmt.Println(genesis)
	fmt.Println(block2)
}