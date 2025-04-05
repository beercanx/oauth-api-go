package token_introspection

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
)

type request struct {
	principal client.Principal
	token     uuid.UUID
}
