// peer communication(http routes + json data exchange)
package main

import(
	"fmt"
	"net/http"
	"encoding/json"
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
func (n *Node)connectpeer(peer string){
	for _,p:=range n.peers{
		if p==peer{
			fmt.Println("already connected to ",peer)
			return
		}
	}
	n.peers=append(n.peers, peer)
	fmt.Println("connected to ",peer)
}
func (n *Node)handlegetPeers(w http.ResponseWriter,r *http.Request){
	json.NewEncoder(w).Encode(n.peers)
}
func (n *Node)handleaddPeer(w http.ResponseWriter,r *http.Request){
	var body map[string]string
	json.NewDecoder(r.Body).Decode(&body)
	peer:=body["address"]
	n.connectpeer(peer)
	json.NewEncoder(w).Encode(map[string]string{"message": "peer added","peer": peer})
}
func (n *Node) handlegetChain(w http.ResponseWriter,r *http.Request){
	json.NewEncoder(w).Encode(n.Blockchain.Blocks)
}
func (n *Node)startserver(){
	http.HandleFunc("/peers",n.handlegetPeers)
	http.HandleFunc("/addpeer",n.handleaddPeer)
	http.HandleFunc("/chain",n.handlegetChain)
	fmt.Println("node running at",n.address)
	http.ListenAndServe(n.address,nil)
}
func main(){
	node:=createNode(":5000")
	node.startserver()
}