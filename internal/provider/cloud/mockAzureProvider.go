package cloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
)

type MockAzureProvider struct {
	SubscribeId string
	Region      string
	VaultName   string
	IsFail      bool
	session     MockAzureSession
}

type MockAzureSession struct {
	keyvault.BaseClient
}

func (session MockAzureSession) GetSecret(
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

	if secretName == "" {
		err = errors.New("the secret name is not available")
	}

	return keyvault.SecretBundle{
		Value: &val,
	}, err
}

func (azProvider MockAzureProvider) GetName() string {
	return "azure"
}

func (azProvider MockAzureProvider) InitialCloudSession() provider.CloudProvider {
	return azProvider
}

func (azProvider MockAzureProvider) RetrieveCredentials() (map[string]string, error) {
	secretName := utils.GetEnv("AZ_SECRET_NAME", "")

	result, err := azProvider.session.GetSecret(
		context.Background(),
		fmt.Sprintf("https://%s.vault.azure.net", azProvider.VaultName),
		secretName,
		"",
	)

	if azProvider.IsFail {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	credentialData := make(map[string]string)
	credentialData[secretName] = *result.Value

	return credentialData, nil
}
