package authentication

import (
	"github.com/alexedwards/argon2id"
	"strings"
	"time"
)

type UserCredential struct {
	username     string
	hashedSecret string
	createdAt    time.Time
	modifiedAt   time.Time
}

type UserCredentialRepository interface {
	Insert(new UserCredential) error
	FindByUsername(username string) (*UserCredential, error)
}

type InMemoryUserCredentialRepository struct {
	store map[string]UserCredential
}

func NewInMemoryUserCredentialRepository() *InMemoryUserCredentialRepository {
	repository := &InMemoryUserCredentialRepository{make(map[string]UserCredential)}

	// TODO - Remove once we've got a means of creating new users
	hash, _ := argon2id.CreateHash("P@55w0rd", argon2id.DefaultParams)
	_ = repository.Insert(UserCredential{"aardvark", hash, time.Now(), time.Now()})

	return repository
}

func (repository InMemoryUserCredentialRepository) Insert(new UserCredential) error {
	repository.store[strings.ToLower(new.username)] = new
	return nil
}

func (repository InMemoryUserCredentialRepository) FindByUsername(username string) (*UserCredential, error) {
	credential, ok := repository.store[strings.ToLower(username)]
	if ok {
		return &credential, nil
	} else {
		return nil, nil
	}
}
