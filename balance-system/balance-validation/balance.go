package main

import "fmt"

type Transaction struct {
	from   string
	to     string
	amount float64
}
type Blockchain struct {
	Transactions []Transaction
}

func (bc *Blockchain) getbalance(address string) float64 {
	balance := 0.0
	for _, tx := range bc.Transactions {
		if tx.from == address {
			balance -= tx.amount
		}
		if tx.to == address {
			balance += tx.amount
		}
	}
	return balance
}
func (bc *Blockchain) addtransac(from, to string, amount float64) bool {
	if from!="system"{
		senderbalance := bc.getbalance(from)
	if senderbalance < amount {
		fmt.Println("transaction failde : ",from, "has insufficient balance")
		return false
	}
	}
	tx:=Transaction{from: from,to: to,amount: amount}
	bc.Transactions=append(bc.Transactions, tx)
	fmt.Println("transaction successfull:", from, "->", to,amount)
	return true
}
func main(){
	blockchain:=Blockchain{}
	//mining reward
	blockchain.addtransac("system","akash",10)
	blockchain.addtransac("akash","boss",4)  //valid
	blockchain.addtransac("boss","dragon",10) //invalid
	blockchain.addtransac("boss","dragon",2) //valid
	fmt.Println("\nfinal balance:")
	fmt.Println("akash : ",blockchain.getbalance("akash"))
	fmt.Println("boss : ",blockchain.getbalance("boss"))
	fmt.Println("dragon : ",blockchain.getbalance("dragon"))
}