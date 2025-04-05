package scope

type Service struct {
	repository Repository
}

func (service *Service) Validate(inputs []string) Scopes {

	var result []Scope

	for _, input := range inputs {

		scope, err := service.repository.FindById(input)
		if err == nil {
			result = append(result, scope)
		}
	}

	return Scopes{result}
}

func NewService(repository Repository) *Service {
	return &Service{repository}
}
