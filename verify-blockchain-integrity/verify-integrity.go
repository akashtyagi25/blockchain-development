//verify blockchain integrity
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	index   int
	data    string
	hash    string
	prehash string
}
type Blockchain struct {
	Blocks []Block
}

func calchash(block Block) string {
	record := string(block.index) + block.data + block.prehash
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
func (bc *Blockchain) addblock(data string){
	preblock:=bc.Blocks[len(bc.Blocks)-1]
	newblock:=Block{
		index: len(bc.Blocks),
		data: data,
		prehash: preblock.hash,
	}
	newblock.hash=calchash(newblock)
	bc.Blocks=append(bc.Blocks, newblock)
}
func (bc *Blockchain) IsValid() bool{
	for i:=1;i<len(bc.Blocks);i++{
		current:=bc.Blocks[i]
		prev:=bc.Blocks[i-1]
		if current.hash!=calchash(current){
			return false
		}
		if current.prehash!=prev.hash{
			return false
		}
	}
	return true
}
func main(){
	genesis:=Block{index: 0,data: "genesis broo",prehash: ""}
	genesis.hash=calchash(genesis)
	blockchain:=Blockchain{Blocks: []Block{genesis}}
	blockchain.addblock("a sent 10 coins to b")
	blockchain.addblock("b sent 2coins to c")
	for _,block:=range blockchain.Blocks{
		fmt.Println("\nblock #",block.index)
		fmt.Println("data : ",block.data)
		fmt.Println("prehash : ",block.prehash)
		fmt.Println("hash : ",block.hash)
	}
	fmt.Println("\nis blockchain valid ? ",blockchain.IsValid());
}