package provider

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/pkg"
)

type AwsProvider struct {
	session       *session.Session
	Region        string
	secretManager *secretsmanager.SecretsManager
	SecretName    string
}

func (awsProvider AwsProvider) InitialCloudSession() pkg.CloudProvider {
	awsProvider.session = session.Must(session.NewSession())
	awsProvider.secretManager = secretsmanager.New(
		awsProvider.session,
		aws.NewConfig().WithRegion(awsProvider.Region),
	)

	return awsProvider
}

func (awsProvider AwsProvider) RetrieveCredentials() map[string]string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(awsProvider.SecretName),
	}

	result, err := awsProvider.secretManager.GetSecretValue(input)
	if err != nil {
		panic(err.Error())
	}

	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	}

	var credentialData map[string]string
	if err := json.Unmarshal([]byte(secretString), &credentialData); err != nil {
		panic(err)
	}

	return credentialData
}
