package restapi

import (
	"fmt"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
)

type MockAzureRestAPI struct {
	IsGetAccessFail bool
	IsGetSecretFail bool
}

func (api MockAzureRestAPI) GetAccessToken() (*AzureAccessToken, error) {
	if api.IsGetAccessFail {
		errorMsg := fmt.Errorf("could not retrieve any credentials")
		return nil, errorMsg
	}

	return &AzureAccessToken{Token: "mock_token"}, nil
}

func (api MockAzureRestAPI) GetSecretValue(accessToken string, vaultURL string, secretName string) (map[string]string, error) {
	credentialData := make(map[string]string)

	if api.IsGetSecretFail {
		errorMsg := fmt.Errorf("could not retrieve any credentials")
		return credentialData, errorMsg
	}

	secretName = utils.GetEnv("AZ_SECRET_NAME", "")

	switch secretName {
	case "db_username":
		credentialData["db_username"] = "admin"
	case "db_password":
		credentialData["db_password"] = "p@ssw0rd"
	}

	return credentialData, nil
}
