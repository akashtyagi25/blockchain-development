// simple balance system utxo based model
package main

import "fmt"
type Transaction struct{
	from string
	to string
	amount float64
}
//blockchain structure
type blockchain struct{
	Transactions []Transaction
}
//balancr calculation func
func(bc *blockchain) getbalance(address string) float64{
	balance:=0.0
	for _,tx:=range bc.Transactions{
		if tx.from==address{
			balance-=tx.amount
		}
		if tx.to==address{
			balance+=tx.amount
		}
	}
	return balance
}
func main(){
	blockchain:=blockchain{
		Transactions: []Transaction{
			{"system","akash",10},
			{"akash","boss",4},
			{"boss","dragon",2},
		},
	}
	fmt.Println("akash balance : ",blockchain.getbalance("akash"))
	fmt.Println("boss balance : ",blockchain.getbalance("boss"))
	fmt.Println("dragon balance : ",blockchain.getbalance("dragon"))
}