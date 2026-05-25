package scope

type Service struct {
	repository Repository
}

func (service *Service) Validate(inputs []string) Scopes {

	var result Scopes

	for _, input := range inputs {

		scope, err := service.repository.FindById(input)
		if err == nil {
			result = append(result, scope)
		}
	}

	return result
}

func NewService(repository Repository) *Service {
	return &Service{repository}
}
