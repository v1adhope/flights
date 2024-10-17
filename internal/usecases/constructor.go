package usecases

type Usecases struct {
	repos Reposer
}

func New(r Reposer) *Usecases {
	return &Usecases{
		repos: r,
	}
}
