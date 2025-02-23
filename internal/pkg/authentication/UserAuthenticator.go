package authentication

type Success struct {
	Username AuthenticatedUsername
}

type Failure struct {
	Reason Reason
}

type Reason string

const (
	Missing    Reason = "missing"
	Mismatched Reason = "mismatched"
	Locked     Reason = "locked"
)

type UserAuthenticator interface {
	Authenticate(username string, password string) (*Success, *Failure, error)
}
