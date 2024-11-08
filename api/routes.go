// api/routes.go
// Implement the API routes for the application.

package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRouter
// @Description    Set up the router for the API.
// @Param          router: *mux.Router (pointer to the router)
// @Return         none
func SetupRouter(router *mux.Router) {
	// Skip cleaning the URL path (enabling empty {id} requests and return 404 instead of 301 redirect)
	router.SkipClean(true)
	
	router.HandleFunc("/receipts/process", ProcessReceiptHandler).Methods(http.MethodPost)
    router.HandleFunc("/receipts/{id}/points", GetPointsHandler).Methods(http.MethodGet)
}