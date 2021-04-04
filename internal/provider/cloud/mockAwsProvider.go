package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
)

type MockAwsProvider struct {
	Region        string
	secretManager MockSecretManagerClient
	SecretName    string
}

type MockSecretManagerClient struct {
	secretsmanageriface.SecretsManagerAPI
}

func (mockSecretManagerClient *MockSecretManagerClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
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

func (awsProvider MockAwsProvider) InitialCloudSession() provider.CloudProvider {
	return awsProvider
}

func (awsProvider MockAwsProvider) RetrieveCredentials() (map[string]string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(awsProvider.SecretName),
	}

	result, err := awsProvider.secretManager.GetSecretValue(input)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve any credentials: %v", err)
		return nil, errors.New(errorMessage)
	}

	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	}

	var credentialData map[string]string
	if err := json.Unmarshal([]byte(secretString), &credentialData); err != nil {
		errorMessage := fmt.Sprintf("Could not unmarshal credentials json to object: %v", err)
		return nil, errors.New(errorMessage)
	}

	return credentialData, nil
}
