package client

import (
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

type SecretRepository interface {
	FindById(id uuid.UUID) (Secret, bool)
	FindByClient(client Id) ([]Secret, bool)
	FindByClientId(clientId string) ([]Secret, bool)
}

type inMemorySecretRepository struct {
	byId       map[uuid.UUID]Secret
	byClientId map[Id][]Secret
}

// assert inMemorySecretRepository implements SecretRepository
var _ SecretRepository = (*inMemorySecretRepository)(nil)

func NewInMemorySecretRepository() SecretRepository {
	repository := &inMemorySecretRepository{make(map[uuid.UUID]Secret), make(map[Id][]Secret)}

	// TODO - Remove once we've got a means of creating new clients
	aardvarkHash, _ := argon2id.CreateHash("badger", argon2id.DefaultParams)
	repository.insert(Secret{uuid.New(), "aardvark", aardvarkHash})

	dodoHash, _ := argon2id.CreateHash("echidna", argon2id.DefaultParams)
	repository.insert(Secret{uuid.New(), "dodo", dodoHash})

	return repository
}

func (i inMemorySecretRepository) insert(secret Secret) {
	i.byId[secret.id] = secret
	i.byClientId[secret.clientId] = append(i.byClientId[secret.clientId], secret)
}

func (i inMemorySecretRepository) FindById(id uuid.UUID) (Secret, bool) {
	secret, ok := i.byId[id]
	return secret, ok
}

func (i inMemorySecretRepository) FindByClient(client Id) ([]Secret, bool) {
	secrets, ok := i.byClientId[client]
	return secrets, ok
}

func (i inMemorySecretRepository) FindByClientId(clientId string) ([]Secret, bool) {
	secrets, ok := i.byClientId[Id(clientId)]
	return secrets, ok
}
