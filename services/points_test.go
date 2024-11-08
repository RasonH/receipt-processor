// services/points_test.go
// The tests for the points calculation service.
// Put tests for points calculation functions (and related helpers) all in this file.

package services

import (
	"testing"

    "receipt-processor/models"

    "github.com/stretchr/testify/assert"
)

////////////////////////////
//     GENERAL CASES 	  //
////////////////////////////

// General test, using the examples from the instruction
// expected: correct points calculation and no error
func TestCalculateTotalPoints_Examples(t *testing.T) {
	receipt := &models.Receipt{
		Retailer: 	  "Target",
		PurchaseDate: "2022-01-01",
	  PurchaseTime: "13:01",
	  Items: []models.Item{
		  {ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
		  {ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		  {ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		  {ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
		  {ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
	},
	Total: "35.35",
  }

  points, err := CalculateTotalPoints(receipt)
  assert.NoError(t, err)
  assert.Equal(t, int64(28), points)

    receipt = &models.Receipt{
        Retailer:     "M&M Corner Market",
        PurchaseDate: "2022-03-20",
        PurchaseTime: "14:33",
        Total:        "9.00",
        Items: []models.Item{
            {ShortDescription: "Gatorade", Price: "2.25"},
            {ShortDescription: "Gatorade", Price: "2.25"},
            {ShortDescription: "Gatorade", Price: "2.25"},
            {ShortDescription: "Gatorade", Price: "2.25"},
        },
    }

    points, err = CalculateTotalPoints(receipt)
    assert.NoError(t, err)
    assert.Equal(t, int64(109), points)

	receipt = &models.Receipt{
        Retailer:     "   M&-M Corner Market s    ", // added spaces to test the trimming, and allowed non-alphanumeric characters
        PurchaseDate: "2022-03-20",
        PurchaseTime: "14:33",
        Total:        "9.00",
        Items: []models.Item{
            {ShortDescription: "Gatorade", Price: "2.25"},
            {ShortDescription: "Gatorade", Price: "2.25"},
            {ShortDescription: "Gatorade", Price: "2.25"},
            {ShortDescription: "Gatorade", Price: "2.25"},
        },
    }
	points, err = CalculateTotalPoints(receipt)
    assert.NoError(t, err)
    assert.Equal(t, int64(110), points)
}

// Test on the helper function `calculateRetailerNamePoints`
// expected: correct points calculation and no error
func TestCalculateRetailerNamePoints(t *testing.T) {
	points, err := calculateRetailerNamePoints("Target abc")
	assert.NoError(t, err)
	assert.Equal(t, int64(9), points)

	points, err = calculateRetailerNamePoints("     Target-abc")
	assert.NoError(t, err)
	assert.Equal(t, int64(9), points)
}

// Test on the helper function `calculatePurchaseDatePoints`
// expected: correct points calculation and no error
func TestCalculatePurchaseDatePoints(t *testing.T) {
	points, err := calculatePurchaseDatePoints("2022-01-01")
	assert.NoError(t, err)
	assert.Equal(t, int64(6), points)

	points, err = calculatePurchaseDatePoints("2022-01-02")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), points)
}

// Test on the helper function `calculatePurchaseTimePoints`
// expected: correct points calculation and no error
func TestCalculatePurchaseTimePoints(t *testing.T) {
	points, err := calculatePurchaseTimePoints("13:01")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), points)
	
	points, err = calculatePurchaseTimePoints("13:59")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculatePurchaseTimePoints("14:00")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), points)

	points, err = calculatePurchaseTimePoints("14:33")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), points)

	points, err = calculatePurchaseTimePoints("15:59")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), points)

	points, err = calculatePurchaseTimePoints("16:00")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), points)
}

// Test on the helper function `calculateItemsPoints`
// expected: correct points calculation and no error
func TestCalculateItemsPoints(t *testing.T) {
	points, err := calculateItemsPoints([]models.Item{
		{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
		{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
		{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
		{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(16), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "0.00"}, // 0 points - TODO: discuss if this is the expected behavior
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "5.01"}, // 2 points
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "5.00"}, // 1 points
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "0.01"}, // 1 points
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "5.01"}, // 2 points
		{ShortDescription: "def", Price: "5.01"}, // 2 points
		{ShortDescription: "ghi", Price: "5.01"}, // 2 points
		{ShortDescription: "jklm", Price: "5.01"}, // 0 points
		// + 5 points for every two items
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(16), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "5.01"}, // 2 points
		{ShortDescription: "def", Price: "5.01"}, // 2 points
		{ShortDescription: "ghi", Price: "5.01"}, // 2 points
		{ShortDescription: "jklm", Price: "5.01"}, // 0 points
		{ShortDescription: "jklm", Price: "5.01"}, // 0 points
		// + 10 points for every two items
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(16), points)
}

// Test on the helper function `calculateTotalAmountPoints`
// expected: correct points calculation and no error
func TestCalculateTotalAmountPoints(t *testing.T) {
	points, err := calculateTotalAmountPoints("35.35")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), points)
	
	points, err = calculateTotalAmountPoints("35.00") // 50 + 25 = 75 points
	assert.NoError(t, err)
	assert.Equal(t, int64(75), points)

	points, err = calculateTotalAmountPoints("35.25") // 25 points
	assert.NoError(t, err)
	assert.Equal(t, int64(25), points)

	points, err = calculateTotalAmountPoints("0.00") // 75 points - TODO: discuss if this is the expected behavior
	assert.NoError(t, err)
	assert.Equal(t, int64(75), points)
}


////////////////////////////
//      EDGE CASES 		  //
////////////////////////////

// Tests on Retailer Name `^[\w\s&-]+$`

// 1. Empty Retailer Name
//    expected: should return an error since it should contain at least one character

// test on helper function `calculateRetailerNamePoints`
func TestCalculateRetailerNamePoints_EmptyRetailer(t *testing.T) {	
	points, err := calculateRetailerNamePoints("")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}

// test on main function `CalculateTotalPoints`
func TestCalculateTotalPoints_EmptyRetailer(t *testing.T) {
	receipt := &models.Receipt{
		Retailer:     "",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err := CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}

// 2. Retailer Name with only spaces
//    expected: should return 0 points, but no error

// test on helper function `calculateRetailerNamePoints`
func TestCalculateRetailerNamePoints_OnlySpacesRetailer(t *testing.T) {
	points, err := calculateRetailerNamePoints("     ")
	assert.NoError(t, err) // No error should be returned TODO: Discuss if this is the expected behavior
	assert.Equal(t, int64(0), points)
}
// test on main helper function `CalculateTotalPoints`
func TestCalculateTotalPoints_OnlySpacesRetailer(t *testing.T) {
	receipt := &models.Receipt{
		Retailer:     "     ",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "00:00",
		Total:        "9.01",
		Items: []models.Item{
			{ShortDescription: "ab", Price: "2.25"},
		},
	}

	points, err := CalculateTotalPoints(receipt)
	assert.NoError(t, err) // No error should be returned TODO: Discuss if this is the expected behavior
	assert.Equal(t, int64(0), points)
}

// 3. Retailer Name with unallowed special characters
//    expected: should return an error since it should contain only letters, numbers, spaces, and the characters "&" and "-"

// test on helper function `calculateRetailerNamePoints`
func TestCalculateRetailerNamePoints_SpecialCharactersRetailer(t *testing.T) {
	points, err := calculateRetailerNamePoints("Invalid@Retailer!")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}
// test on main helper function `CalculateTotalPoints`
func TestCalculateTotalPoints_SpecialCharactersRetailer(t *testing.T) {
	receipt := &models.Receipt{
		Retailer:     "Invalid@Retailer!",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err := CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}

// Tests on Purchase Date `^\d{4}-\d{2}-\d{2}$`
// 1. Invalid Purchase Date
//    expected: should return an error since it should be in the format "YYYY-MM-DD"
// test on helper function `calculatePurchaseDatePoints`
func TestCalculatePurchaseDatePoints_InvalidDate(t *testing.T) {
	points, err := calculatePurchaseDatePoints("invalid_date")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculatePurchaseDatePoints("22-03-20")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculatePurchaseDatePoints("")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}
// test on main helper function `CalculateTotalPoints`
func TestCalculateTotalPoints_InvalidDate(t *testing.T) {
	receipt := &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "invalid_date",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err := CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	receipt = &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "22-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err = CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}


// Tests on Purchase Time `^\d{2}:\d{2}$`
// 1. Invalid Purchase Time
//    expected: should return an error since it should be in the format "HH:MM"
// test on helper function `calculatePurchaseTimePoints`
func TestCalculatePurchaseTimePoints_InvalidTime(t *testing.T) {
	points, err := calculatePurchaseTimePoints("invalid_time")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculatePurchaseTimePoints("14:33:00")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculatePurchaseTimePoints("")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}
// test on main helper function `CalculateTotalPoints`
func TestCalculateTotalPoints_InvalidTime(t *testing.T) {
	receipt := &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "invalid_time",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err := CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	receipt = &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33:00",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err = CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}

// Tests on Items
// 1. No Items
//    expected: should return an error since it should contain at least one item
// test on helper function `calculateItemsPoints`
func TestCalculateItemsPoints_NoItems(t *testing.T) {
	points, err := calculateItemsPoints([]models.Item{})
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}
// test on main helper function `CalculateTotalPoints`
func TestCalculateTotalPoints_NoItems(t *testing.T) {
	receipt := &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{},
	}

	points, err := CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}

// 2. Empty Item Description
//    expected: should return an error since it should contain at least one character after trimming
//    TODO: Discuss if these are the expected behavior
// test on helper function `calculateItemsPoints`
func TestCalculateItemsPoints_EmptyDescription(t *testing.T) {
	points, err := calculateItemsPoints([]models.Item{
		{ShortDescription: "", Price: "5.01"},
	})
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "   ", Price: "5.01"},
	})
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}

// 3. Invalid Item Price
//    expected: should return an error since it should be a non-negative float with 2 decimal places
// test on helper function `calculateItemsPoints`
func TestCalculateItemsPoints_InvalidPrice(t *testing.T) {
	points, err := calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "invalid_price"},
	})
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "5.001"},
	})
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculateItemsPoints([]models.Item{
		{ShortDescription: "abc", Price: "-5.01"},
	})
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}
// test on main helper function `CalculateTotalPoints`
func TestCalculateTotalPoints_InvalidPrice(t *testing.T) {
	// length is not a multiple of 3 but invalid price cases
	receipt := &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "invalid_price"}, 
		},
	}

	points, err := CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	receipt = &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "5.001"},
		},
	}

	points, err = CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	// length is a multiple of 3 and invalid price cases
	receipt = &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "abc", Price: "-5.01"},
		},
	}

	points, err = CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}

// Tests on Total Amount `^\d+\.\d{2}$`
// 1. Invalid Total Amount
//    expected: should return an error since it should be a non-negative float with 2 decimal places
// test on helper function `calculateTotalAmountPoints`
func TestCalculateTotalAmountPoints_InvalidTotal(t *testing.T) {
	points, err := calculateTotalAmountPoints("invalid_total")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculateTotalAmountPoints("5.001")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	points, err = calculateTotalAmountPoints("-5.01")
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)
}
// test on main helper function `CalculateTotalPoints`
func TestCalculateTotalPoints_InvalidTotal(t *testing.T) {
    receipt := &models.Receipt{
        Retailer:     "ValidRetailer",
        PurchaseDate: "2022-03-20",
        PurchaseTime: "14:33",
        Total:        "invalid_total",
        Items: []models.Item{
            {ShortDescription: "Gatorade", Price: "2.25"},
        },
    }

    points, err := CalculateTotalPoints(receipt)
    assert.Error(t, err)
    assert.Equal(t, int64(0), points)

	receipt = &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "5.001",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err = CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

	receipt = &models.Receipt{
		Retailer:     "ValidRetailer",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "-0.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	points, err = CalculateTotalPoints(receipt)
	assert.Error(t, err)
	assert.Equal(t, int64(0), points)

}