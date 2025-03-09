package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {

	const FindByUsername = "FindByUsername"

	credentialRepoError := errors.New("credential repository error")
	statusRepoError := errors.New("status repository error")

	validUsername := "aardvark"
	validPassword := "P@55w0rd"
	validHash := "$argon2id$v=19$m=8,t=1,p=1$ZHl1SVFDUVBlT3JkYkpJRQ$smHA3mizJ+fSojqdxJC+Pg" // P@55w09rd
	validCredential := Credential{validUsername, validHash, time.Now(), time.Now()}

	//
	// Error
	//

	t.Run("when credential repository errors", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.On(FindByUsername, "cred-repo-error").Return(Credential{}, credentialRepoError).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		assert.Panics(t, func() { // TODO - Check for credentialRepoError
			underTest.Authenticate("cred-repo-error", validPassword)
		})

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertNotCalled(t, FindByUsername, mock.Anything)
	})

	t.Run("when argon2 hash checking errors", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.
			On(FindByUsername, "argon2-error").
			Return(&Credential{"argon2-error", "aardvark", time.Now(), time.Now()}, nil).
			Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		assert.Panics(t, func() { // TODO - Check for argon2id.ErrInvalidHash
			underTest.Authenticate("argon2-error", validPassword)
		})

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertNotCalled(t, FindByUsername, mock.Anything)
	})

	t.Run("when status repository errors", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.On(FindByUsername, validUsername).Return(validCredential, nil).Once()
		statusRepository.On(FindByUsername, validUsername).Panic(statusRepoError.Error()).Maybe()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		assert.Panics(t, func() { // TODO - Check for statusRepoError
			underTest.Authenticate(validUsername, validPassword)
		})

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertExpectations(t)
	})

	//
	// Failure
	//

	t.Run("when there is no such credential", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.On(FindByUsername, validUsername).Return(Credential{}, ErrNoSuchCredential).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, "badger")
		assert.Nil(t, success, "success should be nil")
		assert.NotNil(t, failure, "failure should not be nil")
		assert.Equal(t, Missing, failure.Reason)

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertExpectations(t)
	})

	t.Run("when there is a credential mismatch", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.On(FindByUsername, validUsername).Return(validCredential, nil).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, "badger")
		assert.Nil(t, success, "success should be nil")
		assert.NotNil(t, failure, "failure should not be nil")
		assert.Equal(t, Mismatched, failure.Reason)

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertExpectations(t)
	})

	t.Run("when there is no such status", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.On(FindByUsername, validUsername).Return(validCredential, nil).Once()
		statusRepository.On(FindByUsername, validUsername).Return(Status{}, ErrNoSuchStatus).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, validPassword)
		assert.Nil(t, success, "success should be nil")
		assert.NotNil(t, failure, "failure should not be nil")
		assert.Equal(t, Missing, failure.Reason)

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertExpectations(t)
	})

	t.Run("when there is a locked status set", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.On(FindByUsername, validUsername).Return(validCredential, nil).Once()
		statusRepository.On(FindByUsername, validUsername).Return(Status{validUsername, true}, nil).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, validPassword)
		assert.Nil(t, success, "success should be nil")
		assert.NotNil(t, failure, "failure should not be nil")
		assert.Equal(t, Locked, failure.Reason)

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertExpectations(t)
	})

	//
	// Success
	//

	t.Run("when it is all successful", func(t *testing.T) {

		credentialRepository := new(MockedCredentialRepository)
		statusRepository := new(MockedStatusRepository)

		credentialRepository.On(FindByUsername, validUsername).Return(validCredential, nil).Once()
		statusRepository.On(FindByUsername, validUsername).Return(Status{validUsername, false}, nil).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, validPassword)
		assert.Nil(t, failure, "failure should be nil")
		assert.NotNil(t, success, "success should not be nil")
		assert.Equal(t, AuthenticatedUsername{validUsername}, success.Username)

		credentialRepository.AssertExpectations(t)
		statusRepository.AssertExpectations(t)
	})
}

type MockedCredentialRepository struct {
	mock.Mock
}

// assert MockedCredentialRepository implements CredentialRepository
var _ CredentialRepository = &MockedCredentialRepository{}

func (repo *MockedCredentialRepository) Insert(cred Credential) error {
	args := repo.Called(cred)
	return args.Error(0)
}

func (repo *MockedCredentialRepository) FindByUsername(username string) (Credential, error) {
	args := repo.Called(username)
	return args.Get(0).(Credential), args.Error(1)
}

type MockedStatusRepository struct {
	mock.Mock
}

// assert MockedStatusRepository implements StatusRepository
var _ StatusRepository = &MockedStatusRepository{}

func (repo *MockedStatusRepository) Insert(status Status) error {
	args := repo.Called(status)
	return args.Error(0)
}

func (repo *MockedStatusRepository) FindByUsername(username string) (Status, error) {
	args := repo.Called(username)
	return args.Get(0).(Status), args.Error(1)
}
