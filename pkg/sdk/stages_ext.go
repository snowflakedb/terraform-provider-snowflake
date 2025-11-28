package sdk

import "fmt"

func (v *Stage) Location() string {
	return NewStageLocation(v.ID(), "").ToSql()
}

func (s *CreateInternalStageRequest) ID() SchemaObjectIdentifier {
	return s.name
}

func (s *StageCopyOnErrorOptionsRequest) WithSkipFile() *StageCopyOnErrorOptionsRequest {
	s.SkipFile = String("SKIP_FILE")
	return s
}

func (s *StageCopyOnErrorOptionsRequest) WithSkipFileX(x int) *StageCopyOnErrorOptionsRequest {
	s.SkipFile = String(fmt.Sprintf("SKIP_FILE_%d", x))
	return s
}

func (s *StageCopyOnErrorOptionsRequest) WithSkipFileXPercent(x int) *StageCopyOnErrorOptionsRequest {
	s.SkipFile = String(fmt.Sprintf("'SKIP_FILE_%d%%'", x))
	return s
}
