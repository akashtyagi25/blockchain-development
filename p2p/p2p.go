package main

import "fmt"

type Blockchain struct {
	Blocks []string
}
type Node struct {
	address    string
	Blockchain *Blockchain
	peers      []string
}

func createNode(address string) *Node {
	node := &Node{
		address:    address,
		Blockchain: &Blockchain{Blocks: []string{"Genesis block"}},
		peers:      []string{},
	}
	return node
}
func (n *Node) connectpeer(peeraddress string) {
	for _, peer := range n.peers {
		if peer == peeraddress{
			fmt.Println("already connected to: ",peeraddress)
			return
		}
	}
	n.peers=append(n.peers, peeraddress)
	fmt.Println("connected to: ",peeraddress)
}
func main(){
	n1:=createNode("localhost:5000")
	n2:=createNode("localhost:5001")
	n1.connectpeer(n2.address)
	n2.connectpeer(n1.address)
	fmt.Println("node 1 peers: ",n1.peers)
	fmt.Println("node 2 peers : ",n2.peers)
}
