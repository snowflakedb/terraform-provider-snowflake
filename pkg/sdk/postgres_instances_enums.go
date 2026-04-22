package sdk

import (
	"fmt"
	"strings"
)

type PostgresInstanceState string

const (
	PostgresInstanceStateCreating   PostgresInstanceState = "CREATING"
	PostgresInstanceStateRestoring  PostgresInstanceState = "RESTORING"
	PostgresInstanceStateStarting   PostgresInstanceState = "STARTING"
	PostgresInstanceStateReplaying  PostgresInstanceState = "REPLAYING"
	PostgresInstanceStateFinalizing PostgresInstanceState = "FINALIZING"
	PostgresInstanceStateReady      PostgresInstanceState = "READY"
	PostgresInstanceStateRestarting PostgresInstanceState = "RESTARTING"
	PostgresInstanceStateResuming   PostgresInstanceState = "RESUMING"
	PostgresInstanceStateSuspending PostgresInstanceState = "SUSPENDING"
	PostgresInstanceStateSuspended  PostgresInstanceState = "SUSPENDED"
)

var allPostgresInstanceStates = []PostgresInstanceState{
	PostgresInstanceStateCreating,
	PostgresInstanceStateRestoring,
	PostgresInstanceStateStarting,
	PostgresInstanceStateReplaying,
	PostgresInstanceStateFinalizing,
	PostgresInstanceStateReady,
	PostgresInstanceStateRestarting,
	PostgresInstanceStateResuming,
	PostgresInstanceStateSuspending,
	PostgresInstanceStateSuspended,
}

func ToPostgresInstanceState(s string) (PostgresInstanceState, error) {
	s = strings.ToUpper(s)
	for _, state := range allPostgresInstanceStates {
		if string(state) == s {
			return state, nil
		}
	}
	return "", fmt.Errorf("invalid PostgresInstanceState: %s", s)
}

type PostgresInstanceAuthenticationAuthority string

const (
	PostgresInstanceAuthenticationAuthorityPostgres            PostgresInstanceAuthenticationAuthority = "POSTGRES"
	PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake PostgresInstanceAuthenticationAuthority = "POSTGRES_OR_SNOWFLAKE"
)

var allPostgresInstanceAuthenticationAuthorities = []PostgresInstanceAuthenticationAuthority{
	PostgresInstanceAuthenticationAuthorityPostgres,
	PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake,
}

func ToPostgresInstanceAuthenticationAuthority(s string) (PostgresInstanceAuthenticationAuthority, error) {
	s = strings.ToUpper(s)
	for _, auth := range allPostgresInstanceAuthenticationAuthorities {
		if string(auth) == s {
			return auth, nil
		}
	}
	return "", fmt.Errorf("invalid PostgresInstanceAuthenticationAuthority: %s", s)
}
