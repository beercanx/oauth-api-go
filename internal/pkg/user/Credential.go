package user

import (
	"time"
)

type Credential struct {
	username     string
	hashedSecret string
	createdAt    time.Time
	modifiedAt   time.Time
}
