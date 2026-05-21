package user

import (
	"errors"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errorCredentialRepo = errors.New("credential repository error")
var errorStatusRepo = errors.New("status repository error")

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	validUsername := "aardvark"
	validPassword := "P@55w0rd"
	validHash := "$argon2id$v=19$m=8,t=1,p=1$ZHl1SVFDUVBlT3JkYkpJRQ$smHA3mizJ+fSojqdxJC+Pg" // P@55w09rd
	validCredential := Credential{validUsername, validHash, time.Now(), time.Now()}

	//
	// Error
	//

	t.Run("when credential repository errors", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.EXPECT().FindByUsername("cred-repo-error").Return(Credential{}, errorCredentialRepo).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		_, err := underTest.Authenticate("cred-repo-error", validPassword)
		require.ErrorIs(t, err, errorCredentialRepo)
	})

	t.Run("when argon2 hash checking errors", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.
			EXPECT().
			FindByUsername("argon2-error").
			Return(Credential{"argon2-error", "aardvark", time.Now(), time.Now()}, nil).
			Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		_, err := underTest.Authenticate("argon2-error", validPassword)
		require.ErrorIs(t, err, argon2id.ErrInvalidHash)
	})

	t.Run("when status repository errors", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.EXPECT().FindByUsername(validUsername).Return(validCredential, nil).Once()
		statusRepository.EXPECT().FindByUsername(validUsername).Return(Status{}, errorStatusRepo).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		_, err := underTest.Authenticate(validUsername, validPassword)
		require.ErrorIs(t, err, errorStatusRepo)
	})

	//
	// AuthenticationFailure
	//

	t.Run("when there is no such credential", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.EXPECT().FindByUsername(validUsername).Return(Credential{}, ErrNoSuchCredential).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, "badger")
		assert.Zero(t, success, "success should be nil")
		require.Error(t, failure, "failure should not be nil")
		require.ErrorIs(t, failure, AuthenticationFailure{Missing})
	})

	t.Run("when there is a credential mismatch", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.EXPECT().FindByUsername(validUsername).Return(validCredential, nil).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, "badger")
		assert.Zero(t, success, "success should be nil")
		require.Error(t, failure, "failure should not be nil")
		require.ErrorIs(t, failure, AuthenticationFailure{Mismatched})
	})

	t.Run("when there is no such status", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.EXPECT().FindByUsername(validUsername).Return(validCredential, nil).Once()
		statusRepository.EXPECT().FindByUsername(validUsername).Return(Status{}, ErrNoSuchStatus).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, validPassword)
		assert.Zero(t, success, "success should be nil")
		require.Error(t, failure, "failure should not be nil")
		require.ErrorIs(t, failure, AuthenticationFailure{Missing})
	})

	t.Run("when there is a locked status set", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.EXPECT().FindByUsername(validUsername).Return(validCredential, nil).Once()
		statusRepository.EXPECT().FindByUsername(validUsername).Return(Status{validUsername, true}, nil).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, validPassword)
		assert.Zero(t, success, "success should be nil")
		require.Error(t, failure, "failure should not be nil")
		require.ErrorIs(t, failure, AuthenticationFailure{Locked})
	})

	//
	// Authenticated
	//

	t.Run("when it is all successful", func(t *testing.T) {
		t.Parallel()

		credentialRepository := NewMockCredentialRepository(t)
		statusRepository := NewMockStatusRepository(t)

		credentialRepository.EXPECT().FindByUsername(validUsername).Return(validCredential, nil).Once()
		statusRepository.EXPECT().FindByUsername(validUsername).Return(Status{validUsername, false}, nil).Once()

		underTest := NewAuthenticationService(credentialRepository, statusRepository)

		success, failure := underTest.Authenticate(validUsername, validPassword)
		require.NoError(t, failure, "failure should be nil")
		assert.NotNil(t, success, "success should not be nil")
		assert.Equal(t, AuthenticatedUsername{validUsername}, success.Username)
	})
}
