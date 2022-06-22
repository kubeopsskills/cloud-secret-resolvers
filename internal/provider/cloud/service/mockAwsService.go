package service

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type MockAwsService struct {
	Region string
}

func (awsService *MockAwsService) New() {
}

func (awsService *MockAwsService) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	secretId := *input.SecretId
	if secretId == "prod/profile" {
		secretString := "{\"db_username\": \"admin\", \"db_password\": \"p@ssw0rd\"}"
		return &secretsmanager.GetSecretValueOutput{
			ARN:           new(string),
			CreatedDate:   &time.Time{},
			Name:          &secretId,
			SecretBinary:  []byte{},
			SecretString:  &secretString,
			VersionId:     new(string),
			VersionStages: []*string{},
		}, nil
	} else {
		return &secretsmanager.GetSecretValueOutput{
			ARN:           new(string),
			CreatedDate:   &time.Time{},
			Name:          &secretId,
			SecretBinary:  []byte{},
			SecretString:  new(string),
			VersionId:     new(string),
			VersionStages: []*string{},
		}, errors.New("the secret name is not available")
	}
}
