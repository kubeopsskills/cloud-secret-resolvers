package cloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
)

type AzureProvider struct {
	SubscribeId string
	Region      string
	VaultName   string
	session     *keyvault.BaseClient
}

func (azProvider AzureProvider) GetName() string {
	return "azure"
}

func (azProvider AzureProvider) InitialCloudSession() provider.CloudProvider {
	client := keyvault.New()
	authorizer, err := auth.NewAuthorizerFromCLI()
	authorizer.WithAuthorization()

	if err != nil {
		fmt.Println(err)
	} else {
		client.Authorizer = authorizer
		azProvider.session = &client
	}

	return azProvider
}

func (azProvider AzureProvider) RetrieveCredentials() (map[string]string, error) {
	secretName := utils.GetEnv("AZ_SECRET_NAME", "")

	result, err := azProvider.session.GetSecret(
		context.Background(),
		fmt.Sprintf("https://%s.vault.azure.net", azProvider.VaultName),
		secretName,
		"",
	)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	credentialData := make(map[string]string)
	credentialData[secretName] = *result.Value

	return credentialData, nil
}
