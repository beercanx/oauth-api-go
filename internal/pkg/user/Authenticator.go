package user

import "fmt"

type Success struct {
	Username AuthenticatedUsername
}

type Failure struct {
	Reason Reason
}

func (f Failure) Error() string {
	return fmt.Sprintf("Failure: %s", f.Reason)
}

// assert Failure implements error
var _ error = (*Failure)(nil)

type Reason string

const (
	Missing    Reason = "missing"
	Mismatched Reason = "mismatched"
	Locked     Reason = "locked"
)

type Authenticator interface {
	Authenticate(username string, password string) (Success, error)
}
