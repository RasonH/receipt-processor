// models/v1/models.go

// Package models defines the data models used in the application.
package models

// @Title        models/models.go
// @Description  Data models definitions.

// Receipt defines a single receipt.
type Receipt struct {
    ID              string `json:"id"`
    Retailer        string `json:"retailer"`
    PurchaseDate    string `json:"purchaseDate"`
    PurchaseTime    string `json:"purchaseTime"`
    Items           []Item `json:"items"`
    Total           string `json:"total"`
}

// Item defines a single item purchased in a receipt.
type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}


// Equals
// @Description    Check if two Receipt structs are equal. Currently only used for hash colision checks.
// @Param          other: *Receipt
// @Return         true if the two Receipt structs are equal, false otherwise: bool
func (r *Receipt) Equals(other *Receipt) bool {
    if  r.ID != other.ID ||
        r.Retailer != other.Retailer ||
        r.PurchaseDate != other.PurchaseDate ||
        r.PurchaseTime != other.PurchaseTime ||
        r.Total != other.Total ||
        len(r.Items) != len(other.Items) {
        return false
    }

    for i := range r.Items {
        // compare the item fields
        if r.Items[i].ShortDescription != other.Items[i].ShortDescription ||
            r.Items[i].Price != other.Items[i].Price {
            return false
        }
    }

    return true
}
