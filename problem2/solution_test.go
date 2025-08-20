package main

import (
	"testing"
)

func TestNewVendingMachine(t *testing.T) {
	inventory := map[string]Item{
		"Cola": {Name: "Cola", Price: 25, Stock: 5},
	}
	vm := NewVendingMachine(inventory)

	if vm.UserBalance != 0 {
		t.Errorf("NewVendingMachine() initial balance = %v, want 0", vm.UserBalance)
	}
	if len(vm.Inventory) != 1 {
		t.Errorf("NewVendingMachine() inventory size = %v, want 1", len(vm.Inventory))
	}
	if vm.Inventory["Cola"].Stock != 5 {
		t.Errorf("NewVendingMachine() Cola stock = %v, want 5", vm.Inventory["Cola"].Stock)
	}
}

func TestInsertCoin(t *testing.T) {
	tests := []struct {
		name      string
		coinValue int
		wantErr   bool
		errMsg    string
	}{
		{"valid coin 1", 1, false, ""},
		{"valid coin 5", 5, false, ""},
		{"valid coin 10", 10, false, ""},
		{"valid coin 25", 25, false, ""},
		{"invalid coin 2", 2, true, "not a valid coin"},
		{"invalid coin 50", 50, true, "not a valid coin"},
		{"invalid coin 0", 0, true, "not a valid coin"},
		{"invalid coin negative", -5, true, "not a valid coin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewVendingMachine(map[string]Item{})
			err := vm.InsertCoin(tt.coinValue)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InsertCoin() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("InsertCoin() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("InsertCoin() unexpected error = %v", err)
				return
			}

			if vm.UserBalance != tt.coinValue {
				t.Errorf("InsertCoin() balance = %v, want %v", vm.UserBalance, tt.coinValue)
			}
		})
	}
}

func TestInsertMultipleCoins(t *testing.T) {
	vm := NewVendingMachine(map[string]Item{})

	coins := []int{25, 10, 5, 1}
	expectedTotal := 41

	for _, coin := range coins {
		err := vm.InsertCoin(coin)
		if err != nil {
			t.Errorf("InsertCoin(%d) unexpected error = %v", coin, err)
		}
	}

	if vm.UserBalance != expectedTotal {
		t.Errorf("InsertCoin() total balance = %v, want %v", vm.UserBalance, expectedTotal)
	}
}

func TestSelectProduct(t *testing.T) {
	tests := []struct {
		name         string
		inventory    map[string]Item
		userBalance  int
		productName  string
		wantChange   int
		wantErr      bool
		errContains  string
		finalBalance int
		finalStock   int
	}{
		{
			name:         "successful purchase exact amount",
			inventory:    map[string]Item{"Candy": {Name: "Candy", Price: 10, Stock: 5}},
			userBalance:  10,
			productName:  "Candy",
			wantChange:   0,
			wantErr:      false,
			finalBalance: 0,
			finalStock:   4,
		},
		{
			name:         "successful purchase with change",
			inventory:    map[string]Item{"Candy": {Name: "Candy", Price: 10, Stock: 5}},
			userBalance:  15,
			productName:  "Candy",
			wantChange:   5,
			wantErr:      false,
			finalBalance: 0,
			finalStock:   4,
		},
		{
			name:         "insufficient funds",
			inventory:    map[string]Item{"Chips": {Name: "Chips", Price: 35, Stock: 10}},
			userBalance:  25,
			productName:  "Chips",
			wantChange:   0,
			wantErr:      true,
			errContains:  "insuffcient bal",
			finalBalance: 25,
			finalStock:   10,
		},
		{
			name:         "product not found",
			inventory:    map[string]Item{"Candy": {Name: "Candy", Price: 10, Stock: 5}},
			userBalance:  10,
			productName:  "NonExistent",
			wantChange:   0,
			wantErr:      true,
			errContains:  "no product found",
			finalBalance: 10,
			finalStock:   5,
		},
		{
			name:         "out of stock",
			inventory:    map[string]Item{"Cola": {Name: "Cola", Price: 25, Stock: 0}},
			userBalance:  25,
			productName:  "Cola",
			wantChange:   0,
			wantErr:      true,
			errContains:  "not in stock",
			finalBalance: 25,
			finalStock:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewVendingMachine(tt.inventory)
			vm.UserBalance = tt.userBalance

			change, err := vm.SelectProduct(tt.productName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SelectProduct() expected error but got none")
					return
				}
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("SelectProduct() error = %v, want error containing %v", err.Error(), tt.errContains)
				}
				if vm.UserBalance != tt.finalBalance {
					t.Errorf("SelectProduct() balance after error = %v, want %v", vm.UserBalance, tt.finalBalance)
				}
				return
			}

			if err != nil {
				t.Errorf("SelectProduct() unexpected error = %v", err)
				return
			}

			if change != tt.wantChange {
				t.Errorf("SelectProduct() change = %v, want %v", change, tt.wantChange)
			}

			if vm.UserBalance != tt.finalBalance {
				t.Errorf("SelectProduct() final balance = %v, want %v", vm.UserBalance, tt.finalBalance)
			}

			if item, exists := vm.Inventory[tt.productName]; exists {
				if item.Stock != tt.finalStock {
					t.Errorf("SelectProduct() final stock = %v, want %v", item.Stock, tt.finalStock)
				}
			}
		})
	}
}

func TestCancel(t *testing.T) {
	tests := []struct {
		name        string
		userBalance int
		wantRefund  int
	}{
		{"cancel with balance", 25, 25},
		{"cancel with zero balance", 0, 0},
		{"cancel with large balance", 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewVendingMachine(map[string]Item{})
			vm.UserBalance = tt.userBalance

			refund := vm.Cancel()

			if refund != tt.wantRefund {
				t.Errorf("Cancel() refund = %v, want %v", refund, tt.wantRefund)
			}

			if vm.UserBalance != 0 {
				t.Errorf("Cancel() balance after cancel = %v, want 0", vm.UserBalance)
			}
		})
	}
}

func TestCompleteTransactionFlow(t *testing.T) {
	inventory := map[string]Item{
		"Cola":  {Name: "Cola", Price: 25, Stock: 5},
		"Chips": {Name: "Chips", Price: 35, Stock: 10},
		"Candy": {Name: "Candy", Price: 10, Stock: 20},
	}
	vm := NewVendingMachine(inventory)

	vm.InsertCoin(25)
	vm.InsertCoin(25)
	if vm.UserBalance != 50 {
		t.Errorf("After inserting coins, balance = %v, want 50", vm.UserBalance)
	}

	change, err := vm.SelectProduct("Chips")
	if err != nil {
		t.Errorf("SelectProduct() unexpected error = %v", err)
	}
	if change != 15 {
		t.Errorf("SelectProduct() change = %v, want 15", change)
	}
	if vm.UserBalance != 0 {
		t.Errorf("After purchase, balance = %v, want 0", vm.UserBalance)
	}
	if vm.Inventory["Chips"].Stock != 9 {
		t.Errorf("After purchase, Chips stock = %v, want 9", vm.Inventory["Chips"].Stock)
	}

	vm.InsertCoin(25)
	_, err = vm.SelectProduct("Chips")
	if err == nil {
		t.Errorf("SelectProduct() expected error for insufficient funds but got none")
	}

	refund := vm.Cancel()
	if refund != 25 {
		t.Errorf("Cancel() refund = %v, want 25", refund)
	}
	if vm.UserBalance != 0 {
		t.Errorf("After cancel, balance = %v, want 0", vm.UserBalance)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
