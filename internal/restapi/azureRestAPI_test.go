package restapi

import (
	"fmt"
	"reflect"
	"testing"

	resty "github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

func TestAzureGetAccessToken(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	defer httpmock.DeactivateAndReset()

	api := AzureRestAPI{
		Client:       client,
		ClientSecret: "mock_secret",
		ClientId:     "mock_id",
		Resource:     "resource",
		TenantId:     "tenant_id",
	}

	mockResp := `{"access_token": "token_mock"}`
	responder := httpmock.NewBytesResponder(200, []byte(mockResp))
	fakeUrl := "https://login.microsoftonline.com/tenant_id/oauth2/token"
	httpmock.RegisterResponder("POST", fakeUrl, responder)

	result, err := api.GetAccessToken()
	if err != nil {
		t.Errorf("Azure GetAccessToken has error : %v", err)
	}

	got := result.Token
	want := "token_mock"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Azure GetAccessToken = %v, want %v", got, want)
	}
}

func TestAzureGetSecret(t *testing.T) {
	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())

	defer httpmock.DeactivateAndReset()

	api := AzureRestAPI{
		Client:       client,
		ClientSecret: "mock_secret",
		ClientId:     "mock_id",
		Resource:     "resource",
		TenantId:     "tenant_id",
	}

	mockResp := `{"value": "mock_secret_value"}`
	responder := httpmock.NewBytesResponder(200, []byte(mockResp))
	fakeUrl := fmt.Sprintf("%s/secrets/%s?api-version=7.2", "http://mock.vault.com", "mock_secret_name")
	httpmock.RegisterResponder("GET", fakeUrl, responder)

	result, err := api.GetSecretValue("token_mock", "http://mock.vault.com", "mock_secret_name")
	if err != nil {
		t.Errorf("Azure GetSecretValue has error : %v", err)
	}

	got := result
	want := "mock_secret_value"
	if !reflect.DeepEqual(got["mock_secret_name"], want) {
		t.Errorf("Azure GetSecretValue = %v, want %v", got, want)
	}
}
