package sdk

import "fmt"

func (s *ServiceFromSpecificationRequest) WithStageWrapped(stage string) *ServiceFromSpecificationRequest {
	stage = fmt.Sprintf(`@%s`, stage)
	s.Stage = &stage
	return s
}

func (s *ServiceFromSpecificationTemplateRequest) WithStageWrapped(stage string) *ServiceFromSpecificationTemplateRequest {
	stage = fmt.Sprintf(`@%s`, stage)
	s.Stage = &stage
	return s
}
