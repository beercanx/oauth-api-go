package token

import (
	"github.com/google/uuid"
)

type RefreshTokenRepository struct {
	store map[uuid.UUID]RefreshToken
}

func (i *RefreshTokenRepository) Insert(new RefreshToken) error {
	i.store[new.ID] = new
	return nil
}

func (i *RefreshTokenRepository) FindById(id uuid.UUID) (RefreshToken, error) {
	value, ok := i.store[id]
	if ok {
		return value, nil
	}
	return RefreshToken{}, ErrNoSuchToken
}

func (i *RefreshTokenRepository) DeleteById(id uuid.UUID) error {
	delete(i.store, id)
	return nil
}

func (i *RefreshTokenRepository) DeleteByRecord(record RefreshToken) error {
	return i.DeleteById(record.ID)
}

func (i *RefreshTokenRepository) DeletedExpired() error {
	for _, value := range i.store {
		if value.HasExpired() {
			if err := i.DeleteByRecord(value); err != nil {
				return err
			}
		}
	}
	return nil
}

// assert RefreshTokenRepository implements Repository
var _ Repository[RefreshToken] = (*RefreshTokenRepository)(nil)

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{store: make(map[uuid.UUID]RefreshToken)}
}
