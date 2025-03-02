package scope

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (service *Service) Validate(inputs []string) []Scope { // TODO - Review the significance of * in functions of structs

	var result []Scope

	for _, input := range inputs {

		scope := service.repository.FindById(input)
		if scope != nil {
			result = append(result, *scope)
		}
	}

	return result
}
