package gen

import "fmt"

// SuspendIfStartedConfig holds the data needed to generate SuspendIfStarted.
type SuspendIfStartedConfig struct {
	// StateField is the Go field of the SHOW object that holds the state, such as State.
	StateField string
	// StartedConstant and SuspendedConstant are the constants of the state enum.
	StartedConstant, SuspendedConstant string
	// ResumeIfSuspended adds IF SUSPENDED to RESUME, so a resume of an object that someone else already resumed succeeds.
	ResumeIfSuspended bool
}

// SuspendIfStarted declares SuspendIfStarted: it suspends the object when SHOW says it is started, waits until SHOW says it
// is suspended (five checks, one second apart), and returns the function that resumes it. The returned function is never
// nil, and the caller has to call it on every path once the method returned, even when it returned an error.
// The body calls ShowByID and Alter with WithSuspend and WithResume (and WithIfSuspended when ResumeIfSuspended is set).
func (i *Interface) SuspendIfStarted(config SuspendIfStartedConfig) *Interface {
	if config.StateField == "" || config.StartedConstant == "" || config.SuspendedConstant == "" {
		panic(fmt.Sprintf("interface %q: SuspendIfStarted needs StateField, StartedConstant, and SuspendedConstant", i.Name))
	}
	i.SuspendIfStartedConfig = &config
	return i.WithCustomInterfaceMethod(
		"SuspendIfStarted",
		"SuspendIfStarted suspends the object when it is started and waits until it is suspended. The returned function resumes it; it is never nil and has to be called even when an error is returned.",
		[]*MethodParameter{NewMethodParameter("id", i.IdentifierKind)},
		"func(context.Context) error", "error",
	)
}
