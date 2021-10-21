package service

import (
	"errors"

	"github.com/hashicorp/vault/api"
)

type MockVaultService struct {
	Role string
}

func (service MockVaultService) New() (*api.Client, error) {
	return nil, nil
}

func (service MockVaultService) Read(path string) (*api.Secret, error) {
	if path == "kv/data/backend/dev/app" {
		return nil, errors.New("the vault path is not available")
	}
	return &api.Secret{
		RequestID:     "",
		LeaseID:       "",
		LeaseDuration: 0,
		Renewable:     false,
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"username": "username",
				"password": "password",
			},
		},
		Warnings: []string{},
		Auth:     &api.SecretAuth{},
		WrapInfo: &api.SecretWrapInfo{},
	}, nil
}
