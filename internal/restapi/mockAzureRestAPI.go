package restapi

import "fmt"

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

	credentialData = map[string]string{
		"db_name": "admin",
	}

	return credentialData, nil
}
