package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AzureRestAPI struct {
	Client       HTTPClient
	ClientId     string
	ClientSecret string
	Resource     string
	TenantId     string
}

func (api AzureRestAPI) GetAccessToken() (*string, error) {
	reqURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", api.TenantId)

	reqBody := url.Values{}
	reqBody.Set("grant_type", "client_credentials")
	reqBody.Set("client_id", api.ClientId)
	reqBody.Set("client_secret", api.ClientSecret)
	reqBody.Set("resource", api.Resource)

	request, err := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(reqBody.Encode()))
	if err != nil {
		errorMsg := fmt.Errorf("Could not retrieve acesss token: %v", err)
		return nil, errorMsg
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(reqBody.Encode())))

	resp, err := api.Client.Do(request)
	if err != nil {
		errorMsg := fmt.Errorf("Could perform request: %v", err)
		return nil, errorMsg
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMsg := fmt.Errorf("Could read body: %v", err)
		return nil, errorMsg
	}

	var credentialData map[string]string
	if err := json.Unmarshal(body, &credentialData); err != nil {
		errorMsg := fmt.Errorf("Could not unmarshal credentials json to object: %v", err)
		return nil, errorMsg
	}

	var token string = ""
	if credentialData["access_token"] != "" {
		token = credentialData["access_token"]
	}

	return &token, nil
}

func (api AzureRestAPI) GetSecretValue(accessToken string, vaultURL string, secretName string) (map[string]string, error) {
	reqURL := fmt.Sprintf("%s//secrets/%s?api-version=7.2", vaultURL, secretName)

	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		errorMsg := fmt.Errorf("Could not retrieve acesss token: %v", err)
		return nil, errorMsg
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")

	resp, err := api.Client.Do(request)
	if err != nil {
		errorMsg := fmt.Errorf("Could perform request: %v", err)
		return nil, errorMsg
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMsg := fmt.Errorf("Could read body: %v", err)
		return nil, errorMsg
	}

	var rawData map[string]string
	if err := json.Unmarshal(body, &rawData); err != nil {
		errorMsg := fmt.Errorf("Could not unmarshal body json to object: %v", err)
		return nil, errorMsg
	}

	var credentialData map[string]string
	if err := json.Unmarshal([]byte(rawData["value"]), &credentialData); err != nil {
		errorMsg := fmt.Errorf("Could not unmarshal credentials json to object: %v", err)
		return nil, errorMsg
	}

	return credentialData, nil
}
