package user

import "fmt"

type Authenticated struct {
	Username AuthenticatedUsername
}

type AuthenticationFailure struct {
	Reason Reason
}

func (f AuthenticationFailure) Error() string {
	return fmt.Sprintf("AuthenticationFailure: %s", f.Reason)
}

// assert AuthenticationFailure implements error
var _ error = (*AuthenticationFailure)(nil)

type Reason string

const (
	Missing    Reason = "missing"
	Mismatched Reason = "mismatched"
	Locked     Reason = "locked"
)

type Authenticator interface {
	Authenticate(username string, password string) (Authenticated, error)
}
