package transaction

type service struct {
	repo Repository
}

func New(repository Repository) *service {
	return &service{
		repo: repository,
	}
}
