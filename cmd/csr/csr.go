package main

import (
	"fmt"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/csr"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/utils"
	log "github.com/sirupsen/logrus"
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
		awsProvider := cloud.AwsProvider{Region: awsRegion, SecretName: awsSecretName}
		environmentVariableString, err := csr.SyncCredentialKeyFromCloud(awsProvider, keyValueEnvMap)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed as it could not sync any credentials from the cloud provider: %v\n", err)
			log.Fatal(errorMessage)
		}
		if *environmentVariableString == "" {
			log.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
		}
	case cloudType == "azure":
		// TODO: Set default var value
		azRegion := utils.GetEnv("AZ_REGION", "ap-southeast-1")
		azSecretName := utils.GetEnv("AZ_SECRET_NAME", "")
		if azSecretName == "" {
			log.Fatal("No AZ_SECRET_NAME is defined.")
		}
		azTenantId := utils.GetEnv("AZ_TENANT_ID", "")
		if azTenantId == "" {
			log.Fatal("No AZ_TENANT_ID is defined.")
		}
		azVaultURL := utils.GetEnv("AZ_VAULT_URL", "")
		if azVaultURL == "" {
			log.Fatal("No AZ_VAULT_URL is defined.")
		}
		azClientId := utils.GetEnv("AZ_CLIENT_ID", "")
		if azVaultURL == "" {
			log.Fatal("No AZ_CLIENT_ID is defined.")
		}
		azClientSecret := utils.GetEnv("AZ_CLIENT_SECRET", "")
		if azVaultURL == "" {
			log.Fatal("No AZ_CLIENT_SECRET is defined.")
		}
		azResource := utils.GetEnv("AZ_RESOURCE", "resource")

		azureProvider := cloud.AzureProvider{
			AceessTokenRequest: cloud.AzureAceessTokenRequest{
				ClientId:     azClientId,
				ClientSecret: azClientSecret,
				Resource:     azResource,
			},
			Region:     azRegion,
			SecretName: azSecretName,
			TenantId:   azTenantId,
			VaultURL:   azVaultURL,
		}
		environmentVariableString, err := csr.SyncCredentialKeyFromCloud(azureProvider, keyValueEnvMap)
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
