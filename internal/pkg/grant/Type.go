package grant

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Type string

const (
	AuthorisationCode Type = "authorization_code"
	Password          Type = "password"
	RefreshToken      Type = "refresh_token"
	Assertion         Type = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

type Types []Type

var _ sql.Scanner = (*Types)(nil)

func (r *Types) Scan(raw any) error {
	switch source := raw.(type) {
	case string:
		return json.Unmarshal([]byte(source), r)
	case []byte:
		return json.Unmarshal(source, r)
	default:
		return fmt.Errorf("unsupported source type for Types: %T", raw)
	}
}