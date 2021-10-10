package cloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type MockGoogleCloudProvider struct {
	Session MockGoogleCloudSession
}

type MockGoogleCloudSession struct {
	IsFail bool
	secretmanager.Client
}

func (session MockGoogleCloudSession) GetSecret(ctx context.Context, req *secretmanagerpb.GetSecretRequest, opts ...gax.CallOption) (*secretmanagerpb.Secret, error) {
	if session.IsFail {
		return nil, errors.New("the secret name is not available")
	}

	val := ""
	if req.Name != "" {
		switch req.Name {
		case "db_username":
			val = "admin"
		case "db_password":
			val = "p@ssw0rd"
		}
	}

	return &secretmanagerpb.Secret{
		Name: val,
	}, nil
}

func (gcProvider MockGoogleCloudProvider) GetName() string {
	return "gcloud"
}

func (gcProvider MockGoogleCloudProvider) InitialCloudSession() provider.CloudProvider {
	return gcProvider
}

func (gcProvider MockGoogleCloudProvider) RetrieveCredentials() (map[string]string, error) {
	secretName := utils.GetEnv("GC_SECRET_NAME", "")

	req := &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	}
	resp, err := gcProvider.Session.GetSecret(
		context.Background(),
		req,
	)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	data := resp.String()
	data = strings.Replace(data, "name:", "", -1)
	data = strings.Replace(data, "\"", "", -1)

	credentialData := make(map[string]string)
	credentialData[secretName] = data

	return credentialData, nil
}
