package restapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockHttpClient struct{}

func (c *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	mockResponse := `{"access_token": "token_mock"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockResponse)
	}))
	defer ts.Close()
	return http.Get(ts.URL)
}

func TestAzureGetAccessToken(t *testing.T) {
	mockClient := &MockHttpClient{}
	api := AzureRestAPI{
		Client:       mockClient,
		ClientSecret: "mock_secret",
		ClientId:     "mock_id",
		Resource:     "resource",
		TenantId:     "tenant_id",
	}

	want := `{"access_token": "token_mock"}`
	got, _ := api.GetAccessToken()

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("\n Received %s \n but want %s", got, want)
	}

}
