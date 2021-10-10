package cloud

import (
	"context"
	"errors"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type GoogleCloudProvider struct {
	SecretName string
	session    *secretmanager.Client
	context    context.Context
}

func (gcProvider GoogleCloudProvider) GetName() string {
	return "gcloud"
}

func (gcProvider GoogleCloudProvider) InitialCloudSession() provider.CloudProvider {
	// TODO: Connect to google cloud

	gcProvider.context = context.Background()
	client, err := secretmanager.NewClient(gcProvider.context)
	if err != nil {
		fmt.Printf("Google Cloud cannot authentication: %v", err)
	}
	defer client.Close()
	gcProvider.session = client
	return gcProvider
}

func (gcProvider GoogleCloudProvider) RetrieveCredentials() (map[string]string, error) {
	req := &secretmanagerpb.GetSecretRequest{
		Name: gcProvider.SecretName,
	}
	resp, err := gcProvider.session.GetSecret(
		gcProvider.context,
		req,
	)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	// TODO: Convert string to map, But want see what is result
	fmt.Printf("Result -> %s", resp.String())

	credentialData := make(map[string]string)

	return credentialData, nil
}
