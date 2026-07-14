package client

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type Action string

const (
	Authorise               Action = "authorise"
	Introspect              Action = "introspect"
	ProofKeyForCodeExchange Action = "pkce"
)

type Actions []Action

var _ sql.Scanner = (*Actions)(nil)

var ErrUnsupportedActionsSource = errors.New("unsupported source type for Actions")

func (r *Actions) Scan(raw any) error {
	switch source := raw.(type) {
	case string:
		return json.Unmarshal([]byte(source), r)
	case []byte:
		return json.Unmarshal(source, r)
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedActionsSource, raw)
	}
}