// models/models_test.go
// Tests for models package

package models

import (
	"testing"

    "github.com/stretchr/testify/assert"
)

// Test on equals function (comparing two receipts)
func TestEquals(t *testing.T){
	receipt1 := Receipt{
		ID: "1",
		Retailer: "Walmart",
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}
	assert.True(t, receipt1.Equals(&receipt1))

	receipt2 := Receipt{
		ID: "1",
		Retailer: "Walmart",
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}
	assert.True(t, receipt1.Equals(&receipt2))

	receipt3 := Receipt{
		ID: "2", // different ID
		Retailer: "Walmart",
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}
	assert.False(t, receipt1.Equals(&receipt3))

	receipt4 := Receipt{
		ID: "1",
		Retailer: "Walmarts", // different retailer
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}
	assert.False(t, receipt1.Equals(&receipt4))

	receipt5 := Receipt{
		ID: "1",
		Retailer: "Walmart",
		PurchaseDate: "2021-01-02", // different purchase date
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}
	assert.False(t, receipt1.Equals(&receipt5))

	receipt6 := Receipt{
		ID: "1",
		Retailer: "Walmart",
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:01", // different purchase time
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.50",
	}
	assert.False(t, receipt1.Equals(&receipt6))

	receipt7 := Receipt{
		ID: "1",
		Retailer: "Walmart",
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
			{ShortDescription: "Banana", Price: "0.50"}, // different items
		},
		Total: "1.50",
	}
	assert.False(t, receipt1.Equals(&receipt7))

	receipt8 := Receipt{
		ID: "1",
		Retailer: "Walmart",
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Banana", Price: "0.50"}, // different items order
			{ShortDescription: "Apple", Price: "1.00"},
		},
		Total: "1.50",
	}
	assert.False(t, receipt1.Equals(&receipt8))

	receipt9 := Receipt{
		ID: "1",
		Retailer: "Walmart",
		PurchaseDate: "2021-01-01",
		PurchaseTime: "12:00",
		Items: []Item{
			{ShortDescription: "Apple", Price: "1.00"},
			{ShortDescription: "Banana", Price: "0.50"},
		},
		Total: "1.51", // different total amount
	}
	assert.False(t, receipt1.Equals(&receipt9))
}