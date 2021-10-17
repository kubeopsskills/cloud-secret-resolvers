package service

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultService interface {
	New() (*api.Client, error)
	Read(path string) (*api.Secret, error)
}

type VaultServiceImpl struct {
	Role   string
	Client *api.Client
}

func (vaultService *VaultServiceImpl) New() (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultClient, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	vaultJwt, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, fmt.Errorf("unable to read file containing service account token: %w", err)
	}

	params := map[string]interface{}{
		"jwt":  string(vaultJwt),
		"role": vaultService.Role, // the name of the role in Vault that was created with this app's Kubernetes service account bound to it
	}
	// log in to Vault's Kubernetes auth method
	resp, err := vaultClient.Logical().Write("auth/kubernetes/login", params)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.ClientToken == "" {
		return nil, fmt.Errorf("login response did not return client token")
	}

	// now you will use the resulting Vault token for making all future calls to Vault
	vaultClient.SetToken(resp.Auth.ClientToken)
	vaultService.Client = vaultClient
	return vaultClient, nil
}

func (vaultService *VaultServiceImpl) Read(path string) (*api.Secret, error) {
	// get secret from Vault
	return vaultService.Client.Logical().Read(path)
}
