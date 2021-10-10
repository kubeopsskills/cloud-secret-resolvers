package cloud

import (
	"context"
	"errors"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type GoogleCloudProvider struct {
	Credentials string
	ProjectId   string
	session     *secretmanager.Client
	context     context.Context
}

func (gcProvider GoogleCloudProvider) GetName() string {
	return "gcloud"
}

func (gcProvider GoogleCloudProvider) InitialCloudSession() provider.CloudProvider {
	gcProvider.context = context.Background()
	client, err := secretmanager.NewClient(
		gcProvider.context,
		option.WithCredentialsFile(gcProvider.Credentials),
	)
	if err != nil {
		fmt.Printf("Google Cloud cannot authentication: %v", err)
	}
	gcProvider.session = client
	return gcProvider
}

func (gcProvider GoogleCloudProvider) RetrieveCredentials() (map[string]string, error) {
	secretName := utils.GetEnv("GC_SECRET_NAME", "")

	req := &secretmanagerpb.GetSecretRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s", gcProvider.ProjectId, secretName),
	}
	resp, err := gcProvider.session.GetSecret(
		gcProvider.context,
		req,
	)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	// FIXME: Cannot read secret vaule
	// TODO: Convert string to map, But want see what is result
	fmt.Printf("Result -> %v \n\n", resp)

	credentialData := make(map[string]string)
	credentialData[secretName] = resp.String()

	return credentialData, nil
}
