package client

import (
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

type InMemorySecretRepository struct {
	byId       map[uuid.UUID]Secret
	byClientId map[Id][]Secret
}

func (i InMemorySecretRepository) insert(secret Secret) {
	i.byId[secret.id] = secret
	i.byClientId[secret.clientId] = append(i.byClientId[secret.clientId], secret)
}

func (i InMemorySecretRepository) FindById(id uuid.UUID) (Secret, bool) {
	secret, ok := i.byId[id]
	return secret, ok
}

func (i InMemorySecretRepository) FindByClient(client Id) ([]Secret, bool) {
	secrets, ok := i.byClientId[client]
	return secrets, ok
}

func (i InMemorySecretRepository) FindByClientId(clientId string) ([]Secret, bool) {
	secrets, ok := i.byClientId[Id(clientId)]
	return secrets, ok
}

// assert InMemorySecretRepository implements SecretRepository
var _ SecretRepository = (*InMemorySecretRepository)(nil)

func NewInMemorySecretRepository() *InMemorySecretRepository {
	repository := &InMemorySecretRepository{make(map[uuid.UUID]Secret), make(map[Id][]Secret)}

	// TODO - Remove once we've got a means of creating new clients
	aardvarkHash, _ := argon2id.CreateHash("badger", argon2id.DefaultParams)
	repository.insert(Secret{uuid.New(), "aardvark", aardvarkHash})

	dodoHash, _ := argon2id.CreateHash("echidna", argon2id.DefaultParams)
	repository.insert(Secret{uuid.New(), "dodo", dodoHash})

	return repository
}
