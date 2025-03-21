package user

import (
	"github.com/alexedwards/argon2id"
	"strings"
	"time"
)

type InMemoryCredentialRepository struct {
	store map[string]Credential
}

// assert InMemoryCredentialRepository implements CredentialRepository
var _ CredentialRepository = &InMemoryCredentialRepository{}

func NewInMemoryCredentialRepository() *InMemoryCredentialRepository {
	repository := &InMemoryCredentialRepository{make(map[string]Credential)}

	// TODO - Remove once we've got a means of creating new users
	hash, _ := argon2id.CreateHash("P@55w0rd", argon2id.DefaultParams)
	_ = repository.Insert(Credential{"aardvark", hash, time.Now(), time.Now()})

	return repository
}

func (repository *InMemoryCredentialRepository) Insert(new Credential) error {
	repository.store[strings.ToLower(new.username)] = new
	return nil
}

func (repository *InMemoryCredentialRepository) FindByUsername(username string) (Credential, error) {
	if credential, ok := repository.store[strings.ToLower(username)]; ok {
		return credential, nil
	}
	return Credential{}, ErrNoSuchCredential
}
