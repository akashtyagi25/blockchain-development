// consensus conflict resolution+longest chain rule
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Block struct {
	index     int
	timestamp string
	data      string
	prevhash  string
	hash      string
}
type Blockchain struct {
	Blocks []Block
}
type Node struct {
	Blockchain Blockchain
	peers      []string
}

func calchash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s", block.index, block.timestamp, block.data, block.prevhash)
	// hash := sha256.Sum256([]byte(record))
	// return hex.EncodeToString(hash[:])
	return fmt.Sprintf("%x", len(record)) //replace with sha256
}
func isvalidchain(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		prev := chain[i-1]
		curr := chain[i]
		//prev hash must match
		if curr.prevhash != prev.hash {
			return false
		}
		//verify currenrt hash
		if curr.hash !=calchash(curr){
			return false
		}
	}
	return true
}
func(n *Node)resolveConflict()bool{
	var newchain []Block
	maxlength:=len(n.Blockchain.Blocks)
	for _,peer:=range n.peers{
		resp,err:= http.Get("https://"+peer+"/chain")
		if err!=nil{
			continue
		}
		defer resp.Body.Close()
		var peerchain Blockchain
		body,_:=ioutil.ReadAll(resp.Body)
		json.Unmarshal(body,&peerchain)
		length:=len(peerchain.Blocks)
		if length>maxlength &&  isvalidchain(peerchain.Blocks){
			maxlength =length
			newchain = peerchain.Blocks
		}
	}
	if len(newchain)>0{
		n.Blockchain.Blocks=newchain
		fmt.Println("chain replaced with longer valid chain from peers")
		return true
	}
	fmt.Println("current chain is already the longest")
	return false
}
func (n *Node)chainHandler(w http.ResponseWriter,r *http.Request){
	json.NewEncoder(w).Encode(n.Blockchain)
}
func(n *Node)resolvehandler(w http.ResponseWriter,r *http.Request){
	changed:= n.resolveConflict()
	if changed{
		w.Write([]byte("chain replaced with the longest one.\n"))
		
	}else{
		w.Write([]byte("current chain is already the longest \n"))
	}
}
func creategenesisblock()Block{
	return Block{
		index: 0,
		timestamp: time.Now().String(),
		data: "genesis block",
		prevhash: "",
		hash: "genesis-hash",
	}
}
func main(){
	node:=Node{
		Blockchain: Blockchain{Blocks: []Block{creategenesisblock()}},
		peers: []string{"localhost:5001","localhost:5002"}, //add peers here
	}
	http.HandleFunc("/chain",node.chainHandler)
	http.HandleFunc("/resolve",node.resolvehandler)
	fmt.Println("nod running at port 5000..")
	http.ListenAndServe(":5000",nil)
}