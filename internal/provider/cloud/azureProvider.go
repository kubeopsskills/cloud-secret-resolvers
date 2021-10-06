package cloud

import (
	"errors"
	"fmt"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/restapi"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
)

type AzureProvider struct {
	Region   string
	VaultURL string

	API restapi.IAzureRestAPI

	accessToken string
}

func (azProvider AzureProvider) GetName() string {
	return "azure"
}

func (azProvider AzureProvider) InitialCloudSession() provider.CloudProvider {
	result, err := azProvider.API.GetAccessToken()
	if err != nil {
		fmt.Printf("Could not retrieve access token: %v", err)
	}
	azProvider.accessToken = result.Token
	return azProvider
}

func (azProvider AzureProvider) RetrieveCredentials() (map[string]string, error) {
	secretName := utils.GetEnv("AZ_SECRET_NAME", "")
	result, err := azProvider.API.GetSecretValue(
		azProvider.accessToken,
		azProvider.VaultURL,
		secretName,
	)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}
	return result, nil
}
