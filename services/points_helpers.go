// services/v1/points_helpers.go
// Implementions for points calculation on different recipt info.
// Seperated some functions in case of future need to use them in other services, and rule expansion.
// TODO: Could reorganize into different files if the number of rules increase.

package services

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"receipt-processor/models"
)

/*
	--- Points calculation rules summary ---
	1: One point for every alphanumeric character in the retailer name.
	2: 50 points if the total is a round dollar amount with no cents.
	3: 25 points if the total is a multiple of 0.25.
	4: 5 points for every two items on the receipt.
	5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	6. 6 points if the day in the purchase date is odd.
	7. 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	...
*/


////////////////////////
//    MAIN HELPERS    //
////////////////////////


// calculateRetailerNamePoints
// @Description    Calculate points based on retailer's name.
//                 Included rules:
//				   		1: One point for every alphanumeric character in the retailer name.
//                 Assumption:
// 						retailer name compiles pattern "^[\\w\\s\\-&]+$"
// @Param          retailer: string
// @Return         points from retailer name: int64, error: error
func calculateRetailerNamePoints(retailer string) (int64, error) {
	var points int64 = 0

	// Assumption: retailer name compiles pattern "^[\\w\\s\\-&]+$" 
	// TODO: can refactor validation to models, and simplify the function
	validRetailerName := regexp.MustCompile(`^[\w\s&-]+$`)
	if !validRetailerName.MatchString(retailer) {
		return 0, fmt.Errorf("[calculateRetailerNamePoints] Retailer name is invalid %v", retailer)
	}

	// 1: One point for every alphanumeric character in the retailer name.
	count := countAlphanumericChar(retailer)
	points += count
	
	return points, nil 
}


// calculatePurchaseDatePoints
// @Description    Calculate points based on purchase date.
//                 Included rules:
//				   		6: 6 points if the day in the purchase date is odd.
// @Param          purchaseDate: string
// @Return         points from purchase date: int64, error: error
func calculatePurchaseDatePoints(purchaseDate string) (int64, error) {
	var points int64 = 0

	// parse purchase date
	date, err := time.Parse("2006-01-02", purchaseDate)
	if err != nil {
		return 0, fmt.Errorf("[calculatePurchaseDatePoints] Failed to parse purchase date %v: %w", purchaseDate, err)
	}

	// 6: 6 points if the day in the purchase date is odd.
	if date.Day() & 1 == 1 {
		points += 6
	}

	return points, nil
}


// calculatePurchaseTimePoints
// @Description    Calculate points based on purchase time.
//                 Included rules:
//				   		7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
// @Param          purchaseTime: string
// @Return         points from purchase time: int64, error: error
func calculatePurchaseTimePoints(purchaseTime string) (int64, error) {
	var points int64 = 0

	// parse purchase time
	time, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		return 0, fmt.Errorf("[calculatePurchaseTimePoints] Failed to parse purchase time %v: %w", purchaseTime, err)
	}

	// 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	// TODO: definition check for before and after
	if time.Hour() >= 14 && time.Hour() < 16 {	
		points += 10
	}

	return points, nil
}


// calculateItemsPoints
// @Description    Calculate points based on items.
//                 Included rules:
//				   		4: 5 points for every two items on the receipt.
//				   		5: If the trimmed length of the item description is a multiple of 3,
//						   multiply the price by 0.2 and round up to the nearest integer.
// 				   Assumptions:
//						items should have at least 1 item.
//						trimmed description should not be empty.
//						item price compiles pattern "^\\d+\\.\\d{2}$", which is a non-negative float with 2 decimal places.
// @Param          items: []models.Item
// @Return         points from items: int64, error: error
func calculateItemsPoints(items []models.Item) (int64, error) {
	var points int64 = 0

	// 4: 5 points for every two items on the receipt.
	itemsCount := len(items)
	// Assumption: items should have at least 1 item.
	if itemsCount == 0 { 
		return 0, fmt.Errorf("[calculateItemsPoints] No items found in the receipt %v", items)
	}
	points += int64(itemsCount / 2 * 5)

	// 5: If the trimmed length of the item description is a multiple of 3,
	//    multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range items {
		trimmed, err := trimDescription(item.ShortDescription)
		if err != nil {
			return 0, fmt.Errorf("[calculateItemsPoints] Failed to check item description %v: %w", item.ShortDescription, err)
		}

		// Assumption: trimmed description should not be empty
		if len(trimmed) == 0 { 
			return 0, fmt.Errorf("[calculateItemsPoints] Trimmed description is empty for %v", item.ShortDescription)
		}

		// Assumption: item price compiles pattern "^\\d+\\.\\d{2}$", which is a non-negative float with 2 decimal places.
		validItemPrice := regexp.MustCompile(`^\d+\.\d{2}$`)
		if !validItemPrice.MatchString(item.Price) {
			return 0, fmt.Errorf("[calculateItemsPoints] Invalid item price %v", item.Price)
		}
		// check if the trimmed length of the item description is a multiple of 3
		if len(trimmed) % 3 == 0 {
			// parse item price
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, fmt.Errorf("[calculateItemsPoints] Failed to parse item price %v: %w", item.Price, err)
			}
			// Round up to the nearest integer
			points += int64(math.Ceil(price * 0.2)) 
		}
	}
	return points, nil
}

// calculateTotalAmountPoints
// @Description    Calculate points based on total amount.
//                 Included rules:
//				   		2: 50 points if the total is a round dollar amount with no cents.
//				   		3: 25 points if the total is a multiple of 0.25.
// 				   Assumptions:
//						total amount compiles pattern "^\\d+\\.\\d{2}$", which is a non-negative float with 2 decimal places.
// @Param          total: string
// @Return         points from total amount: int64, error: error
func calculateTotalAmountPoints(total string) (int64, error) {
	var points int64 = 0
	
	// Assumption: total amount compiles pattern "^\\d+\\.\\d{2}$", which is a non-negative float with 2 decimal places.
	validTotal := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !validTotal.MatchString(total) {
		return 0, fmt.Errorf("[calculateTotalAmountPoints] Invalid total amount %v", total)
	}

	// parse total amount
	floatTotal, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, fmt.Errorf("[calculateTotalAmountPoints] Failed to parse total amount %v: %w", total, err)
	}

	totalToCents := int64(floatTotal * 100)

	// 2: 50 points if the total is a round dollar amount with no cents.
	if totalToCents % 100 == 0 {
		points += 50
	}

	// 3: 25 points if the total is a multiple of 0.25.
	if totalToCents % 25 == 0 {
		points += 25
	}

	return points, nil
}


////////////////////////
// HELPERS OF HELPERS //
////////////////////////

// countAlphanumericChar
// @Description    Count the number of alphanumeric characters in a string.
// @Param          str: string
// @Return         count of alphanumeric characters: int64, error: error
func countAlphanumericChar(str string) int64 {
	var count int64 = 0

	// Count the number of alphanumeric characters in the string.
	for _, c := range str {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			count++
		}
	}

	return count
}

// trimDescription
// @Description    Trim the description string.
// 				   Assumption:
// 						trimmed description should not be empty 
// @Param          description: string
// @Return         trimmed description: string, error: error
// TODO: might not need it if receipt data were always normalized (trimmed) before hashing, which is prior to calculations.
func trimDescription(description string) (string, error) {
	trimmed := strings.TrimSpace(description)
	
	if len(trimmed) == 0 {
		return "", fmt.Errorf("[trimDescription] Trimmed description is empty for %v", description)
	}
	return trimmed, nil
}