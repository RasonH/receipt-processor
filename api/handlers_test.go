// api/handlers_test.go
// Tests for the API handlers.

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"receipt-processor/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
    router := mux.NewRouter()
    SetupRouter(router)
    return router
}

// Test on ProcessReceiptHandler function
// 1. general case (using example receipt) - 200 OK
func TestProcessReceiptHandler(t *testing.T) {
	router := setupRouter()

    receipt := models.Receipt{
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

    requestBody, _ := json.Marshal(receipt)
    req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var response map[string]string
    err := json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)
    id, exists := response["id"]
    assert.True(t, exists)
    assert.NotEmpty(t, id)
}

// 2. invalid JSON format - 400 Bad Request
func TestProcessReceiptHandlerInvalidJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// 3. hash collision - 409 Conflict
// TODO: hard to test seperately if not refactoring the code

// 4. invalid receipt - 400 Bad Request
func TestProcessReceiptHandlerInvalidReceipt(t *testing.T) {
	router := setupRouter()

	// invalid time format
	receipt := models.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33:00", 
		Total:        "9.00",
		Items: []models.Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}

	requestBody, _ := json.Marshal(receipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// empty receipt
	receipt = models.Receipt{}

	requestBody, _ = json.Marshal(receipt)
	req, _ = http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Tests on GetReceiptHandler function
// simplifying by using the same processed receipt
func TestGetReceiptHandler(t *testing.T) {
	router := setupRouter()

    // process the example receipt
    receipt := models.Receipt{
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

    requestBody, _ := json.Marshal(receipt)
    req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    var response map[string]string
    json.Unmarshal(rr.Body.Bytes(), &response)
    id := response["id"]

    // get the points with the receipt ID
	// 1. general case - 200 OK
    req, _ = http.NewRequest("GET", "/receipts/"+id+"/points", nil)
    rr = httptest.NewRecorder()
    router.ServeHTTP(rr, req)
    assert.Equal(t, http.StatusOK, rr.Code)

    var pointsResponse map[string]int64
    err := json.Unmarshal(rr.Body.Bytes(), &pointsResponse)
    assert.NoError(t, err)
    points, exists := pointsResponse["points"]
    assert.True(t, exists)
    assert.Equal(t, int64(109), points)

	// 2. unexisting ID - 404 Not Found
	req, _ = http.NewRequest("GET", "/receipts/0/points", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// 3. empty ID - 404 Not Found
	req, _ = http.NewRequest("GET", "/receipts//points", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}