package provider

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type awsProvider struct {
	session       *session.Session
	region        string
	secretManager *secretsmanager.SecretsManager
	secretName    string
}

func (awsProvider awsProvider) initialCloudSession() bool {
	awsProvider.session = session.Must(session.NewSession())
	awsProvider.secretManager = secretsmanager.New(
		awsProvider.session,
		aws.NewConfig().WithRegion(awsProvider.region),
	)
	return awsProvider.secretManager != nil
}

func (awsProvider awsProvider) retrieveCredentials() map[string]interface{} {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(awsProvider.secretName),
	}
	result, err := awsProvider.secretManager.GetSecretValue(input)
	if err != nil {
		panic(err.Error())
	}

	var secretBinary []byte
	if result.SecretBinary != nil {
		secretBinary = result.SecretBinary
	}

	var credentialData map[string]interface{}
	if err := json.Unmarshal(secretBinary, &credentialData); err != nil {
		panic(err)
	}

	return credentialData
}
