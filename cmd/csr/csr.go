package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/csr"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud/service"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
)

func main() {
	// load environment configs
	keyValueEnvMap := csr.LoadCredentialKeyFromEnvironment()
	cloudType := utils.GetEnv("CLOUD_TYPE", "aws")
	log.Infof("Syncing credentials from %s ...", cloudType)
	switch {
	case cloudType == "aws":
		awsRegion := utils.GetEnv("AWS_REGION", "ap-southeast-1")
		awsSecretName := utils.GetEnv("AWS_SECRET_NAME", "")
		if awsSecretName == "" {
			log.Fatal("No AWS_SECRET_NAME is defined.")
		}
		service := service.AWSServiceImpl{}
		awsProvider := cloud.AwsProvider{
			Service:    &service,
			Region:     awsRegion,
			SecretName: awsSecretName,
		}
		environmentVariableString, err := csr.SyncCredentialKeyFromCloud(awsProvider, keyValueEnvMap)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed as it could not sync any credentials from the cloud provider: %v\n", err)
			log.Fatal(errorMessage)
		}
		if *environmentVariableString == "" {
			log.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
		}
	case cloudType == "azure":
		azVaultName := utils.GetEnv("AZ_VAULT_NAME", "")
		if azVaultName == "" {
			log.Fatal("No AZ_VAULT_NAME is defined.")
		}
		service := service.AzureServiceImpl{}
		azureProvider := cloud.AzureProvider{
			Service:   &service,
			VaultName: azVaultName,
		}
		environmentVariableString, err := csr.SyncCredentialKeyFromCloud(azureProvider, keyValueEnvMap)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed as it could not sync any credentials from the cloud provider: %v\n", err)
			log.Fatal(errorMessage)
		}
		if *environmentVariableString == "" {
			log.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
		}
	case cloudType == "gcloud":
		gcCreds := utils.GetEnv("GOOGLE_APPLICATION_CREDENTIALS", "")
		if gcCreds == "" {
			log.Fatal("No GOOGLE_APPLICATION_CREDENTIALS is defined.")
		}
		gcProjectId := utils.GetEnv("GOOGLE_PROJECT_ID", "")
		if gcProjectId == "" {
			log.Fatal("No GOOGLE_PROJECT_ID is defined.")
		}

		service := service.GoogleCloudServiceImpl{}
		gcProvider := cloud.GoogleCloudProvider{
			ProjectId: gcProjectId,
			Service:   &service,
		}
		environmentVariableString, err := csr.SyncCredentialKeyFromCloud(gcProvider, keyValueEnvMap)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed as it could not sync any credentials from the cloud provider: %v\n", err)
			log.Fatal(errorMessage)
		}
		if *environmentVariableString == "" {
			log.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
		}
	}
	log.Info("Synced")

}
