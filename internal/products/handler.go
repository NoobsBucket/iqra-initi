package products

import (
	"net/http"
	"github.com/NoobsBucket/iqra-initi/internal/json"
)

type handler struct {
	service service
}

func NewHandler(service service) *handler {
	return &handler{
		service: service,
	}
}
func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// call the service method to get the products
	// return the json in a http response
	products := []string{"Product 1", "Product 2", "Product 3"}
	json.SendJson(w, products, http.StatusOK)

}
