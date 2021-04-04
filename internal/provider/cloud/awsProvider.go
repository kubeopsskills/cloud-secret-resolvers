package cloud

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
)

type AwsProvider struct {
	session       *session.Session
	Region        string
	secretManager *secretsmanager.SecretsManager
	SecretName    string
}

func (awsProvider AwsProvider) InitialCloudSession() provider.CloudProvider {
	awsProvider.session = session.Must(session.NewSession())
	awsProvider.secretManager = secretsmanager.New(
		awsProvider.session,
		aws.NewConfig().WithRegion(awsProvider.Region),
	)

	return awsProvider
}

func (awsProvider AwsProvider) RetrieveCredentials() (map[string]string, error) {
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
