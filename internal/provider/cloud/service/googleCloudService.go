package service

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type GoogleCloudService interface {
	NewClient(context context.Context) error
	AccessSecretVersion(
		ctx context.Context,
		req *secretmanagerpb.AccessSecretVersionRequest,
		opts ...gax.CallOption,
	) (*secretmanagerpb.AccessSecretVersionResponse, error)
}

type GoogleCloudServiceImpl struct {
	session *secretmanager.Client
}

func (service *GoogleCloudServiceImpl) NewClient(context context.Context) error {
	session, err := secretmanager.NewClient(context)
	if err != nil {
		return err
	}
	service.session = session
	return err
}

func (service *GoogleCloudServiceImpl) AccessSecretVersion(
	ctx context.Context,
	req *secretmanagerpb.AccessSecretVersionRequest,
	opts ...gax.CallOption,
) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return service.session.AccessSecretVersion(
		ctx,
		req,
	)
}
