package client

import (
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

type InMemorySecretRepository struct {
	byId       map[uuid.UUID]Secret
	byClientId map[string][]Secret
}

func (i InMemorySecretRepository) insert(secret Secret) {
	i.byId[secret.id] = secret
	i.byClientId[secret.clientId.Value] = append(i.byClientId[secret.clientId.Value], secret)
}

func (i InMemorySecretRepository) FindById(id uuid.UUID) (Secret, bool) {
	secret, ok := i.byId[id]
	return secret, ok
}

func (i InMemorySecretRepository) FindByClient(client Id) ([]Secret, bool) {
	secrets, ok := i.byClientId[client.Value]
	return secrets, ok
}

func (i InMemorySecretRepository) FindByClientId(clientId string) ([]Secret, bool) {
	secrets, ok := i.byClientId[clientId]
	return secrets, ok
}

// assert InMemorySecretRepository implements SecretRepository
var _ SecretRepository = (*InMemorySecretRepository)(nil)

func NewInMemorySecretRepository() *InMemorySecretRepository {
	repository := &InMemorySecretRepository{make(map[uuid.UUID]Secret), make(map[string][]Secret)}

	// TODO - Remove once we've got a means of creating new clients
	hash, _ := argon2id.CreateHash("badger", argon2id.DefaultParams)
	repository.insert(Secret{uuid.New(), Id{"aardvark"}, hash})

	return repository
}
