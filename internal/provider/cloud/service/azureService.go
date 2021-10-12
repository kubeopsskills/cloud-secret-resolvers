package service

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzureService interface {
	New()
	GetSecret(ctx context.Context, vaultBaseURL string, secretName string, secretVersion string) (result keyvault.SecretBundle, err error)
}

type AzureServiceImpl struct {
	session *keyvault.BaseClient
}

func (azService *AzureServiceImpl) New() {
	client := keyvault.New()
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("cannot make an authentication to Azure: %v", err)
	} else {
		client.Authorizer = authorizer
		azService.session = &client
	}
}

func (azService *AzureServiceImpl) GetSecret(ctx context.Context, vaultBaseURL string, secretName string, secretVersion string) (result keyvault.SecretBundle, err error) {
	return azService.session.GetSecret(ctx, vaultBaseURL, secretName, secretVersion)
}
