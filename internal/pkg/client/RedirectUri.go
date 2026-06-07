package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type RedirectUri string
type RedirectUris []RedirectUri

var _ sql.Scanner = (*RedirectUris)(nil)

func (r *RedirectUris) Scan(raw any) error {
	switch source := raw.(type) {
	case string:
		return json.Unmarshal([]byte(source), r)
	case []byte:
		return json.Unmarshal(source, r)
	default:
		return fmt.Errorf("unsupported source type for RedirectUris: %T", raw)
	}
}