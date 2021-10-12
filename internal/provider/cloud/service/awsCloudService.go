package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AWSCloudService interface {
	New()
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

type AWSCloudServiceImpl struct {
	Region        string
	secretManager *secretsmanager.SecretsManager
}

func (awsService *AWSCloudServiceImpl) New() {
	awsSession := session.Must(session.NewSession())
	awsService.secretManager = secretsmanager.New(
		awsSession,
		aws.NewConfig().WithRegion(awsService.Region),
	)
}

func (awsService *AWSCloudServiceImpl) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return awsService.secretManager.GetSecretValue(input)
}
