package restapi

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

type AzureRestAPI struct {
	Client       *resty.Client
	ClientId     string
	ClientSecret string
	Resource     string
	TenantId     string
}

type AzureAccessToken struct {
	Token string `json:"access_token"`
}

type AzureSecretValue struct {
	Value string `json:"value"`
}

func (api AzureRestAPI) GetAccessToken() (*AzureAccessToken, error) {
	reqURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", api.TenantId)

	reqBody := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     api.ClientId,
		"client_secret": api.ClientSecret,
		"resource":      api.Resource,
	}

	resp, err := api.Client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(reqBody).
		SetResult(AzureAccessToken{}).
		Post(reqURL)
	if err != nil {
		errorMsg := fmt.Errorf("could not perfrom request: %v", err)
		return nil, errorMsg
	}

	var result AzureAccessToken
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		errorMsg := fmt.Errorf("could not unmarshal credentials json to object: %v", err)
		return nil, errorMsg
	}

	return &result, nil
}

func (api AzureRestAPI) GetSecretValue(accessToken string, vaultURL string, secretName string) (map[string]string, error) {
	reqURL := fmt.Sprintf("%s/secrets/%s?api-version=7.2", vaultURL, secretName)
	resp, err := api.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(AzureSecretValue{}).
		Get(reqURL)

	if err != nil {
		errorMsg := fmt.Errorf("could not perfrom request: %v", err)
		return nil, errorMsg
	}

	var secretValue AzureSecretValue
	if err := json.Unmarshal(resp.Body(), &secretValue); err != nil {
		errorMsg := fmt.Errorf("could not unmarshal credentials json to object: %v", err)
		return nil, errorMsg
	}

	credentialData := map[string]string{
		secretName: secretValue.Value,
	}

	return credentialData, nil
}
