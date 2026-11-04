package v1

type api struct {
	//orderv1.UnimplementedHandler
	service Service
}

func NewApi(service Service) *api {
	return &api{
		service: service,
	}
}
