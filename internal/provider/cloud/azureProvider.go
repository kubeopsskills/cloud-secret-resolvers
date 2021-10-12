package cloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud/service"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
)

type AzureProvider struct {
	VaultName string
	Service   service.AzureService
}

func (azProvider AzureProvider) GetName() string {
	return "azure"
}

func (azProvider AzureProvider) InitialCloudSession() provider.CloudProvider {
	azProvider.Service.New()
	return azProvider
}

func (azProvider AzureProvider) RetrieveCredentials() (map[string]string, error) {
	secretName := utils.GetEnv("AZ_SECRET_NAME", "")

	result, err := azProvider.Service.GetSecret(
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
