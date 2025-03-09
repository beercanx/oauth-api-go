package user

import (
	"errors"
	"github.com/alexedwards/argon2id"
	"strings"
	"time"
)

type Credential struct {
	username     string
	hashedSecret string
	createdAt    time.Time
	modifiedAt   time.Time
}

type CredentialRepository interface {
	Insert(new Credential) error                        // TODO - Verify if Go DB libraries panic or return errors
	FindByUsername(username string) (Credential, error) // TODO - Verify if Go DB libraries panic or return errors
}

var (
	ErrNoSuchCredential = errors.New("credential does not exist")
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
