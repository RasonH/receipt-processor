// storage/storage.go
// In memory storage for data

// Storage package provides in-memory storage for the application.
package storage

import (
	"sync"

	"receipt-processor/models"
)

// ReceiptData is a struct that holds the receipt info and the calculated points associated with it.
type ReceiptData struct {
	Receipt models.Receipt
	Points int64
}

// Storage is where we map receipt IDs to their data.
// We do not store invalid receipts in the storage
type Storage struct {
	mu sync.RWMutex
	data map[string]ReceiptData
}

// ensuring the singleton pattern
var (
	storageInstance *Storage
	once sync.Once
)

// GetStorageInstance
// @Description    Get the singleton instance of the storage
// @Param          none
// @Return         pointer to the storage instance: *Storage
func GetStorageInstance() *Storage {
	once.Do(func() {
		storageInstance = &Storage{
			data: make(map[string]ReceiptData),
		}
	})
	return storageInstance
}

// SaveReceipt
// @Description    Save a receipt and its calculated points to the storage
// @Param          id: string, receipt: models.Receipt, points: int64
// @Return         none
func (s *Storage) SaveReceipt(id string, receipt models.Receipt, points int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = ReceiptData{
		Receipt: receipt,
		Points: points,
	}
}

// GetReceiptData
// @Description    Retrieve the receipt data from the storage based on the receipt ID
// @Param          id: string
// @Return         receipt data: ReceiptData, found: bool
func (s *Storage) GetReceiptData(id string) (ReceiptData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, found := s.data[id]
	return data, found
}