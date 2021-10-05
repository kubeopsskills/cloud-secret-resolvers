package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
)

type AzureProvider struct {
	AceessTokenRequest AzureAceessTokenRequest
	Region             string
	SecretName         string
	TenantId           string
	VaultURL           string

	accessToken string
}

type AzureAceessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Resource     string `json:"resource"`
}

func (azProvider AzureProvider) InitialCloudSession() provider.CloudProvider {
	timeout := 5000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	reqURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", azProvider.TenantId)

	reqBody := url.Values{}
	reqBody.Set("grant_type", "client_credentials")
	reqBody.Set("client_id", azProvider.AceessTokenRequest.ClientId)
	reqBody.Set("client_secret", azProvider.AceessTokenRequest.ClientSecret)
	reqBody.Set("resource", azProvider.AceessTokenRequest.Resource)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Add("Content-Length", strconv.Itoa(len(reqBody.Encode())))

	res, err := client.Post(reqURL, strings.NewReader(reqBody.Encode()), nil)
	if err != nil {
		fmt.Printf("Could not retrieve acesss token: %v", err)
	}

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Could not read body from api: %v", err)
	}

	var credentialData map[string]string
	if err := json.Unmarshal(resp, &credentialData); err != nil {
		fmt.Printf("Could not unmarshal credentials json to object: %v", err)
	}

	token := credentialData["access_token"]
	if token != "" {
		azProvider.accessToken = credentialData["access_token"]
	}

	return azProvider
}

func (azProvider AzureProvider) RetrieveCredentials() (map[string]string, error) {
	// TODO: Use rest api to get secret value
	// https://docs.microsoft.com/en-us/rest/api/keyvault/get-secret/get-secret
	// GET {vaultBaseUrl}/secrets/{secret-name}/?api-version=7.2

	timeout := 5000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", azProvider.accessToken))
	headers.Set("Content-Type", "application/json")

	url := fmt.Sprintf("%s/secrets/%s?api-version=7.2", azProvider.VaultURL, azProvider.SecretName)
	res, err := client.Get(url, nil)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	secretString, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not read body from api: %v", err)
		return nil, errors.New(errorMessage)
	}

	var credentialData map[string]string
	if err := json.Unmarshal(secretString, &credentialData); err != nil {
		errorMessage := fmt.Sprintf("Could not unmarshal credentials json to object: %v", err)
		return nil, errors.New(errorMessage)
	}

	return credentialData, nil
}
