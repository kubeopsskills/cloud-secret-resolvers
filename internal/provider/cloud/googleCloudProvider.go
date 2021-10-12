package cloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud/service"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type GoogleCloudProvider struct {
	ProjectId string
	Service   service.GoogleCloudService
	context   context.Context
}

func (gcProvider GoogleCloudProvider) GetName() string {
	return "gcloud"
}

func (gcProvider GoogleCloudProvider) InitialCloudSession() provider.CloudProvider {
	gcProvider.context = context.Background()
	err := gcProvider.Service.NewClient(gcProvider.context)
	if err != nil {
		fmt.Printf("Cannot make an authentication to Google Cloud: %v", err)
	}
	return gcProvider
}

func (gcProvider GoogleCloudProvider) RetrieveCredentials() (map[string]string, error) {
	secretName := utils.GetEnv("GC_SECRET_NAME", "")

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", gcProvider.ProjectId, secretName),
	}
	// FIXME: Pointer invalid
	resp, err := gcProvider.Service.AccessSecretVersion(
		gcProvider.context,
		req,
	)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	credentialData := make(map[string]string)
	credentialData[secretName] = string(resp.Payload.Data)

	return credentialData, nil
}
