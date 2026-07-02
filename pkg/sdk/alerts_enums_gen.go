package sdk

type AlertAction string

var (
	// AlertActionResume makes a suspended alert active.
	AlertActionResume AlertAction = "RESUME"
	// AlertActionSuspend puts the alert into a "Suspended" state.
	AlertActionSuspend AlertAction = "SUSPEND"
)

type AlertState string

var (
	AlertStateStarted   AlertState = "started"
	AlertStateSuspended AlertState = "suspended"
)
