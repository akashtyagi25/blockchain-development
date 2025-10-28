// block broadcast + chain synchronization
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)
type Blockchain struct{
	Blocks []string
}
type Node struct{
	address string
	Blockchain *Blockchain
	peers []string
}
func createNode(address string) *Node{
	return &Node{
		address: address,
		Blockchain: &Blockchain{Blocks: []string{"genesis block"}},
		peers: []string{},
	}
}
func (n *Node) connectpeer(peer string){
	for _,p:=range n.peers{
		if p==peer{
			return
		}
	}
	n.peers=append(n.peers, peer)
	fmt.Println("connected to: ",peer)
}
func (n *Node) addblock(data string){
	newblock:=fmt.Sprintf("block %d: %s",len(n.Blockchain.Blocks),data)
	n.Blockchain.Blocks=append(n.Blockchain.Blocks, newblock)
	fmt.Println("new block added locally : ",newblock)
	n.broadcastblock(newblock)
}
func (n *Node) broadcastblock(block string){
	for _,peer:=range n.peers{
		url:="http://"+peer+"/receiveblock"
		body,_:=json.Marshal(map[string]string{"block": block})
		http.Post(url,"application/json",bytes.NewBuffer(body))
		fmt.Println("broadcast block to",peer)
	}
}
func (n *Node) handlereceiveblock(w http.ResponseWriter,r *http.Request){
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	block:=data["block"]
	//add block
	n.Blockchain.Blocks=append(n.Blockchain.Blocks, block)
	fmt.Println("received new block: ",block)
	json.NewEncoder(w).Encode(map[string]string{"status":"block received"})
}
func (n *Node) handlesyncChain(w http.ResponseWriter,r *http.Request){
	for _,peer:=range n.peers{
		resp,err:=http.Get("http://"+peer+"/chain")
		if err!=nil{
			continue
		}
		body,_:=ioutil.ReadAll(resp.Body)
		var peerchain []string
		json.Unmarshal(body,&peerchain)
		if len(peerchain) > len(n.Blockchain.Blocks){
			n.Blockchain.Blocks=peerchain
			fmt.Println("chain updated from peer:",peer)
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"status":"chain synchronized"})
}
func (n *Node) handlegetchain(w http.ResponseWriter,r *http.Request){
	json.NewEncoder(w).Encode(n.Blockchain.Blocks)
}
func (n *Node) handleaddpeer(w http.ResponseWriter,r *http.Request){
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	peer:=body["address"]
	n.connectpeer(peer)
	json.NewEncoder(w).Encode(map[string]string{"message":"peer added"})
}
func (n *Node) startserver(){
	http.HandleFunc("/chain",n.handlegetchain)
	http.HandleFunc("/receiveblock",n.handlereceiveblock)
	http.HandleFunc("/addpeer",n.handleaddpeer)
	http.HandleFunc("/sync",n.handlesyncChain)
	fmt.Println("node running at",n.address)
	http.ListenAndServe(n.address,nil)
}
func main(){
	node:=createNode("localhost:5000")
	node.connectpeer("localhost:5001")
	node.addblock("Transaction from Node 1")
	node.startserver()
}