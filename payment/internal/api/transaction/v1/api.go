package v1

import paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"

type api struct {
	paymentv1.UnimplementedPaymentServiceServer
	service Service
}

func NewApi(service Service) *api {
	return &api{
		service: service,
	}
}
