package cloud

import (
	"fmt"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud/service"
)

type VaultProvider struct {
	Service service.VaultService
	Role    string
	Path    string
}

func (vaultProvider VaultProvider) GetName() string {
	return "vault"
}

func (vaultProvider VaultProvider) InitialCloudSession() provider.CloudProvider {
	vaultProvider.Service.New()
	return vaultProvider
}

func (vaultProvider VaultProvider) RetrieveCredentials() (map[string]string, error) {
	// get secret from Vault
	secret, err := vaultProvider.Service.Read(vaultProvider.Path)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data type assertion failed: %T %#v", secret.Data["data"], secret.Data["data"])
	}

	finalData := make(map[string]string)
	for key, value := range data {
		finalData[key] = value.(string)
	}
	// data map can contain more than one key-value pair
	return finalData, nil
}
