package service

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
)

type MockAzureService struct {
	IsFail bool
}

func (azService MockAzureService) New() {
}

func (azService MockAzureService) GetSecret(
	ctx context.Context,
	vaultBaseURL string,
	secretName string,
	secretVersion string,
) (result keyvault.SecretBundle, err error) {
	val := ""
	if secretName != "" {
		switch secretName {
		case "db_username":
			val = "admin"
		case "db_password":
			val = "p@ssw0rd"
		}
	}

	if secretName == "" || azService.IsFail {
		err = errors.New("the secret name is not available")
	}

	return keyvault.SecretBundle{
		Value: &val,
	}, err
}
