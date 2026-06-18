package user

import (
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
)

type inMemoryCredentialRepository struct {
	store map[string]Credential
}

func (repository *inMemoryCredentialRepository) Insert(new Credential) error {
	repository.store[strings.ToLower(new.username)] = new
	return nil
}

func (repository *inMemoryCredentialRepository) FindByUsername(username string) (Credential, error) {
	if credential, ok := repository.store[strings.ToLower(username)]; ok {
		return credential, nil
	}
	return Credential{}, ErrNoSuchCredential
}

// assert inMemoryCredentialRepository implements CredentialRepository
var _ CredentialRepository = (*inMemoryCredentialRepository)(nil)

func NewInMemoryCredentialRepository() CredentialRepository {
	repository := &inMemoryCredentialRepository{make(map[string]Credential)}

	// TODO - Remove once we've got a means of creating new users
	hash, _ := argon2id.CreateHash("P@55w0rd", argon2id.DefaultParams)
	_ = repository.Insert(Credential{"aardvark", hash, time.Now(), time.Now()})

	return repository
}
