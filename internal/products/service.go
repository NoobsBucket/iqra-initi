package products
import (
	"context"
)
type service interface {
	// Define the methods that the service should implement
	getProducts(ctx context.Context) ([]string, error)
}
type scv struct {
	//repository repository
}

func NewService() service {
	return &scv{
		//repository: repository,
	}
}
func (s *scv) getProducts(ctx context.Context) ([]string, error) {
	// Implement the logic to retrieve products from the repository
	// For now, return a static list of products
	return nil, nil

}