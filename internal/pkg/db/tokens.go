package db

import "time"

func (token *AccessToken) HasExpired() bool {
	return time.Now().After(token.ExpiresAt)
}

func (token *AccessToken) IsBefore() bool {
	return time.Now().Before(token.NotBefore)
}
