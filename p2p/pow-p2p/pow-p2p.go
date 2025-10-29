// pow integration with p2p network
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	index     int
	timestamp string
	data      string
	prevhash  string
	hash      string
	nonce     int
}
type Blockchain struct {
	Blocks []Block
}
type Node struct {
	address    string
	Blockchain *Blockchain
	peers      []string
	difficulty int
}

func calchash(block Block) string {
	rec := strconv.Itoa(block.index) + block.timestamp + block.data + block.prevhash + strconv.Itoa(block.nonce)
	hash := sha256.Sum256([]byte(rec))
	return hex.EncodeToString(hash[:])
}
func mineblock(data string, prevhash string, difficulty int, index int) Block {
	var newblock Block
	newblock.index = index
	newblock.timestamp = time.Now().String()
	newblock.data = data
	newblock.prevhash = prevhash
	prefix := strings.Repeat("0", difficulty)
	for {
		hash := calchash(newblock)
		if strings.HasPrefix(hash, prefix) {
			newblock.hash = hash
			break
		}
		newblock.nonce++
	}
	fmt.Println("block mined ", newblock.hash)
	return newblock
}
func (bc *Blockchain) addblock(block Block) {
	bc.Blocks = append(bc.Blocks, block)
}
func (bc *Blockchain) getlasthash() string {
	return bc.Blocks[len(bc.Blocks)-1].hash
}
func createnode(address string, difficulty int) *Node {
	genesis := Block{index: 0, timestamp: time.Now().String(), data: "genesis block", prevhash: "", hash: "GENESIS"}
	return &Node{
		address:    address,
		Blockchain: &Blockchain{Blocks: []Block{genesis}},
		peers:      []string{},
		difficulty: difficulty,
	}
}
func (n *Node) connectpeer(peer string) {
	for _, p := range n.peers {
		if p == peer {
			return
		}
	}
	n.peers = append(n.peers, peer)
	fmt.Println("connected to : ", peer)
}
func (n *Node) broadcastblock(block Block) {
	for _, peer := range n.peers {
		url := "http://" + peer + "/receiveblock"
		body, _ := json.Marshal(block)
		http.Post(url, "application/json", bytes.NewBuffer(body))
		fmt.Println("broadcasted block to", peer)
	}
}
func (n *Node) handlemineblock(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	data := body["data"]
	prevhash := n.Blockchain.getlasthash()
	newblock := mineblock(data, prevhash, n.difficulty, len(n.Blockchain.Blocks))
	n.Blockchain.addblock(newblock)
	n.broadcastblock(newblock)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "block mined and broadcasted successfully",
		"hash":    newblock.hash,
	})
}
func (n *Node) handlereceiveblock(w http.ResponseWriter, r *http.Request) {
	var block Block
	json.NewDecoder(r.Body).Decode(&block)
	lasthash := n.Blockchain.getlasthash()
	if block.prevhash == lasthash {
		n.Blockchain.addblock(block)
		fmt.Println("block received and added ", block.hash)
	} else {
		fmt.Println("block rejected: invalid chain sequence")
	}
	json.NewEncoder(w).Encode(map[string]string{"status":"received"})
}
func (n *Node)handleaddpeer(w http.ResponseWriter, r *http.Request){
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	peer:=body["address"]
	n.connectpeer(peer)
	json.NewEncoder(w).Encode(map[string]string{"message":"peer added"})
}
func (n *Node)handlegetchain(w http.ResponseWriter,r *http.Request){
	json.NewEncoder(w).Encode(n.Blockchain.Blocks)
}
func (n *Node)startserver(){
	http.HandleFunc("/mineblock",n.handlemineblock)
	http.HandleFunc("/receiveblock",n.handlereceiveblock)
	http.HandleFunc("/addpeer",n.handleaddpeer)
	http.HandleFunc("/chain",n.handlegetchain)
	fmt.Println("node running at : ",n.address)
	http.ListenAndServe(n.address,nil)
}
func main(){
	node:=createnode("localhost:5000",3)
	node.connectpeer("localhost:5001")
	node.startserver()
}
