package main

import (
	"fmt"
	"sync"
)

type Item struct {
	Name  string
	Price int
	Stock int
}

type VendingMachine struct {
	UserBalance int
	Inventory   map[string]Item
	sync.Mutex
}

var validCoins = map[int]bool{
	1:  true,
	5:  true,
	10: true,
	25: true,
}

func main() {
	// 1. Initialize the machine with some inventory
	inventory := map[string]Item{
		"Cola":  {Name: "Cola", Price: 25, Stock: 5},
		"Chips": {Name: "Chips", Price: 35, Stock: 10},
		"Candy": {Name: "Candy", Price: 10, Stock: 20},
	}
	vm := NewVendingMachine(inventory)

	fmt.Println("Vending Machine is ready.")
	fmt.Println("---")

	// --- Scenario 1: Successful Purchase ---
	fmt.Println("## Scenario 1: Successful Purchase of Candy ##")
	fmt.Println("Inserting a 10 coin...")
	if err := vm.InsertCoin(10); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("Selecting 'Candy' (Price 10)...")
	change, err := vm.SelectProduct("Candy")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Product dispensed. Change: %d\n", change)
	}
	fmt.Println("---")

	// --- Scenario 2: Insufficient Funds & Cancel ---
	fmt.Println("## Scenario 2: Insufficient Funds & Cancel ##")
	fmt.Println("Inserting a 25 coin...")
	if err := vm.InsertCoin(25); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("Selecting 'Chips' (Price 35)...")
	change, err = vm.SelectProduct("Chips")
	if err != nil {
		// This is the expected outcome
		fmt.Printf("Error (as expected): %v\n", err)
	} else {
		fmt.Printf("Success! Product dispensed. Change: %d\n", change)
	}
	fmt.Println("Transaction incomplete. Cancelling...")
	refund := vm.Cancel()
	fmt.Printf("Refunded amount: %d\n", refund)
	fmt.Println("---")

	// --- Scenario 3: Purchase with Change ---
	fmt.Println("## Scenario 3: Purchase with Change ##")
	fmt.Println("Inserting a 25 coin...")
	vm.InsertCoin(25)
	fmt.Println("Inserting another 25 coin...")
	vm.InsertCoin(25)
	fmt.Println("Total inserted: 50")
	fmt.Println("Selecting 'Chips' (Price 35)...")
	change, err = vm.SelectProduct("Chips")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Product dispensed. Change: %d\n", change)
	}
	fmt.Println("---")

	// --- Scenario 4: Out of Stock ---
	fmt.Println("## Scenario 4: Item is Out of Stock ##")
	// Buy the remaining 4 Colas
	for i := 0; i < 5; i++ {
		vm.InsertCoin(25)
		vm.SelectProduct("Cola")
	}
	fmt.Println("All Colas have been purchased.")

	fmt.Println("Attempting to buy one more Cola...")
	vm.InsertCoin(25)
	_, err = vm.SelectProduct("Cola")
	if err != nil {
		fmt.Printf("Error (as expected): %v\n", err)
	}
	fmt.Println("---")
}

func NewVendingMachine(inventory map[string]Item) *VendingMachine {
	return &VendingMachine{
		UserBalance: 0,
		Inventory:   inventory,
	}
}

func (v *VendingMachine) InsertCoin(coinValue int) error {
	if !validCoins[coinValue] {
		return fmt.Errorf("not a valid coin")
	}
	v.Lock()
	v.UserBalance += coinValue
	v.Unlock()
	return nil
}

func (v *VendingMachine) SelectProduct(productName string) (int, error) {
	v.Lock()
	defer v.Unlock()
	prod, ok := v.Inventory[productName]
	if !ok {
		return 0, fmt.Errorf("no product found with name : %v", productName)
	}
	if prod.Stock == 0 {
		return 0, fmt.Errorf("%s not in stock", productName)
	}
	if v.UserBalance < prod.Price {
		return 0, fmt.Errorf("insuffcient bal : %d", v.UserBalance)
	}
	prod.Stock--
	v.Inventory[productName] = prod
	change := v.UserBalance - prod.Price
	v.UserBalance = 0
	return change, nil
}

func (v *VendingMachine) Cancel() int {
	v.Lock()
	defer v.Unlock()
	change := v.UserBalance
	v.UserBalance = 0
	return change
}
