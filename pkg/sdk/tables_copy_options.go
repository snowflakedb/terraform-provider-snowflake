package sdk

import "fmt"

// Builder functions for StageCopyOptionsRequest

func NewStageCopyOptionsRequest() *StageCopyOptionsRequest {
	return &StageCopyOptionsRequest{}
}

func (s *StageCopyOptionsRequest) WithOnError(onError StageCopyOnErrorOptionsRequest) *StageCopyOptionsRequest {
	s.OnError = &onError
	return s
}

func (s *StageCopyOptionsRequest) WithSizeLimit(sizeLimit int) *StageCopyOptionsRequest {
	s.SizeLimit = &sizeLimit
	return s
}

func (s *StageCopyOptionsRequest) WithPurge(purge bool) *StageCopyOptionsRequest {
	s.Purge = &purge
	return s
}

func (s *StageCopyOptionsRequest) WithReturnFailedOnly(returnFailedOnly bool) *StageCopyOptionsRequest {
	s.ReturnFailedOnly = &returnFailedOnly
	return s
}

func (s *StageCopyOptionsRequest) WithMatchByColumnName(matchByColumnName StageCopyColumnMapOption) *StageCopyOptionsRequest {
	s.MatchByColumnName = &matchByColumnName
	return s
}

func (s *StageCopyOptionsRequest) WithEnforceLength(enforceLength bool) *StageCopyOptionsRequest {
	s.EnforceLength = &enforceLength
	return s
}

func (s *StageCopyOptionsRequest) WithTruncatecolumns(truncatecolumns bool) *StageCopyOptionsRequest {
	s.Truncatecolumns = &truncatecolumns
	return s
}

func (s *StageCopyOptionsRequest) WithForce(force bool) *StageCopyOptionsRequest {
	s.Force = &force
	return s
}

// Builder functions for StageCopyOnErrorOptionsRequest

func NewStageCopyOnErrorOptionsRequest() *StageCopyOnErrorOptionsRequest {
	return &StageCopyOnErrorOptionsRequest{}
}

func (s *StageCopyOnErrorOptionsRequest) WithContinue_(continue_ bool) *StageCopyOnErrorOptionsRequest {
	s.Continue_ = &continue_
	return s
}

func (s *StageCopyOnErrorOptionsRequest) WithAbortStatement(abortStatement bool) *StageCopyOnErrorOptionsRequest {
	s.AbortStatement = &abortStatement
	return s
}

// WithSkipFile sets SkipFile to "SKIP_FILE"
func (s *StageCopyOnErrorOptionsRequest) WithSkipFile() *StageCopyOnErrorOptionsRequest {
	s.SkipFile = String("SKIP_FILE")
	return s
}

// WithSkipFileX sets SkipFile to "SKIP_FILE_n" where n is the provided integer
func (s *StageCopyOnErrorOptionsRequest) WithSkipFileX(x int) *StageCopyOnErrorOptionsRequest {
	s.SkipFile = String(fmt.Sprintf("SKIP_FILE_%d", x))
	return s
}

// WithSkipFileXPercent sets SkipFile to "'SKIP_FILE_n%'" where n is the provided integer
func (s *StageCopyOnErrorOptionsRequest) WithSkipFileXPercent(x int) *StageCopyOnErrorOptionsRequest {
	s.SkipFile = String(fmt.Sprintf("'SKIP_FILE_%d%%'", x))
	return s
}
