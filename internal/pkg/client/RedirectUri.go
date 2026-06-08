package client

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type RedirectUri string
type RedirectUris []RedirectUri

var _ sql.Scanner = (*RedirectUris)(nil)

var ErrUnsupportedRedirectUrisSource = errors.New("unsupported source type for RedirectUris")

func (r *RedirectUris) Scan(raw any) error {
	switch source := raw.(type) {
	case string:
		return json.Unmarshal([]byte(source), r)
	case []byte:
		return json.Unmarshal(source, r)
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedRedirectUrisSource, raw)
	}
}