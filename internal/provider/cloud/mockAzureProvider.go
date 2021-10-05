package cloud

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
)

type MockAzureProvider struct {
	Region     string
	SecretName string
}

func (azProvider MockAzureProvider) InitialCloudSession() provider.CloudProvider {
	return azProvider
}

func (azProvider MockAzureProvider) RetrieveCredentials() (map[string]string, error) {
	secretString := "{\"db_username\": \"admin\", \"db_password\": \"p@ssw0rd\"}"

	var credentialData map[string]string
	if err := json.Unmarshal([]byte(secretString), &credentialData); err != nil {
		errorMessage := fmt.Sprintf("Could not unmarshal credentials json to object: %v", err)
		return nil, errors.New(errorMessage)
	}

	return credentialData, nil
}
