package v1

import inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer
	service Service
}

func NewApi(service Service) *api {
	return &api{
		service: service,
	}
}
