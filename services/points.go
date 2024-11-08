// services/v1/points.go
// Main logic for points calculation.

// Services includes the business logic of the application.
package services

import (
	"fmt"

	"receipt-processor/models"
)


////////////////////////////
//  POINTS SERVICE LOGIC  //
////////////////////////////


// CalculateTotalPoints
// @Description    calculates the points earned from a given receipt.
//				   Assumptions:
// 						any error during the process will stop the calculation and return 0 points with an error.
// @Param          pointer to the receipt object: *models.Receipt
// @Return         total points earned: int64, error: error
func CalculateTotalPoints(receipt *models.Receipt) (int64, error) {
	var totalPoints int64 = 0 // assuming int64 is large enough to avoid overflow, and aligns with the API definition

	// Points calculation rules are based on the following information on a receipt:

	// points received from retailer's name
	retailerNamePoints, err := calculateRetailerNamePoints(receipt.Retailer)
	if err != nil {
		return 0, fmt.Errorf("[CalculateTotalPoints] Failed to calculate retailer name points for receipt ID %v: %w", receipt.ID, err)
	}
	totalPoints += retailerNamePoints

	// points received from puchase date
	puchaseDatePoints, err := calculatePurchaseDatePoints(receipt.PurchaseDate)
	if err != nil {
		return 0, fmt.Errorf("[CalculateTotalPoints] Failed to calculate purchase date points for receipt ID %v: %w", receipt.ID, err)
	}
	totalPoints += puchaseDatePoints

	// points received from purchase time
	purchaseTimePoints, err := calculatePurchaseTimePoints(receipt.PurchaseTime)
	if err != nil {
		return 0, fmt.Errorf("[CalculateTotalPoints] Failed to calculate purchase time points for receipt ID %v: %w", receipt.ID, err)
	}
	totalPoints += purchaseTimePoints

	// points received from items
	itemsPoints, err := calculateItemsPoints(receipt.Items)
	if err != nil {
		return 0, fmt.Errorf("[CalculateTotalPoints] Failed to calculate items points for receipt ID %v: %w", receipt.ID, err)
	}
	totalPoints += itemsPoints
	
	// points received from total amount
	totalAmountPoints, err := calculateTotalAmountPoints(receipt.Total)
	if err != nil {
		return 0, fmt.Errorf("[CalculateTotalPoints] Failed to calculate total amount points for receipt ID %v: %w", receipt.ID, err)
	}
	totalPoints += totalAmountPoints

	return totalPoints, nil
}