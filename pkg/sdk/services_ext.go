package sdk

import "fmt"

func (s *CreateServiceRequest) WithSpecificationFileWrapped(spec string) *CreateServiceRequest {
	spec = fmt.Sprintf(`$$%s$$`, spec)
	s.WithFromSpecification(*NewServiceFromSpecificationRequest().WithSpecification(spec))
	return s
}

func (s *CreateServiceRequest) WithSpecificationTemplateFileWrapped(spec string, using []ListItem) *CreateServiceRequest {
	spec = fmt.Sprintf(`$$%s$$`, spec)
	s.WithFromSpecificationTemplate(*NewServiceFromSpecificationTemplateRequest(using).WithSpecificationTemplate(spec))
	return s
}

func (s *AlterServiceRequest) WithSpecificationFileWrapped(spec string) *AlterServiceRequest {
	spec = fmt.Sprintf(`$$%s$$`, spec)
	s.WithFromSpecification(*NewServiceFromSpecificationRequest().WithSpecification(spec))
	return s
}

func (s *AlterServiceRequest) WithSpecificationTemplateFileWrapped(spec string, using []ListItem) *AlterServiceRequest {
	spec = fmt.Sprintf(`$$%s$$`, spec)
	s.WithFromSpecificationTemplate(*NewServiceFromSpecificationTemplateRequest(using).WithSpecificationTemplate(spec))
	return s
}

func (s *ServiceFromSpecificationRequest) WithStageWrapped(Stage string) *ServiceFromSpecificationRequest {
	Stage = fmt.Sprintf(`@%s`, Stage)
	s.Stage = &Stage
	return s
}

func (s *ServiceFromSpecificationTemplateRequest) WithStageWrapped(Stage string) *ServiceFromSpecificationTemplateRequest {
	Stage = fmt.Sprintf(`@%s`, Stage)
	s.Stage = &Stage
	return s
}
