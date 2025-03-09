package scope

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (service *Service) Validate(inputs []string) []Scope {

	var result []Scope

	for _, input := range inputs {

		scope, err := service.repository.FindById(input)
		if err == nil {
			result = append(result, scope)
		}
	}

	return result
}
