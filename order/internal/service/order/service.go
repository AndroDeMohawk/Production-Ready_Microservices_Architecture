package order

import "github.com/go-faster/errors"

var ErrNilDependency = errors.New("зависимость не передана")

type orderService struct {
	repository      Repository
	inventoryClient InventoryClient
	paymentClient   PaymentClient
}

func NewService(repository Repository, inventoryClient InventoryClient, paymentClient PaymentClient) *orderService {
	return &orderService{
		repository:      repository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
