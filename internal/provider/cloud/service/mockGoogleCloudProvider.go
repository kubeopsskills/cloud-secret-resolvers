package service

import (
	"context"
	"errors"

	gax "github.com/googleapis/gax-go/v2"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type MockGoogleCloudService struct {
	IsFail bool
}

func (service MockGoogleCloudService) NewClient(context context.Context) error {
	if service.IsFail {
		return errors.New("cannot make an authentication to Google Cloud")
	}
	return nil
}

func (service MockGoogleCloudService) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	if service.IsFail {
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

	return &secretmanagerpb.AccessSecretVersionResponse{
		Name: req.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(val),
		},
	}, nil
}
