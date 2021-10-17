package csr

import (
	"os"
	"testing"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider/cloud/service"
)

func TestMain(t *testing.T) {

	var expectedKeys []string
	expectedKeys = append(expectedKeys, "db_username")
	expectedKeys = append(expectedKeys, "db_password")
	var expectedValues []string
	expectedValues = append(expectedValues, "${db_username}")
	expectedValues = append(expectedValues, "${db_password}")

	os.Setenv("db_username", "${db_username}")
	os.Setenv("db_password", "${db_password}")
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	counter := 0
	for key, value := range keyValueEnvMap {
		if key != expectedKeys[counter] {
			t.Fatalf("Failed as key [%s] is not equal to expected key [%s].", key, expectedKeys[counter])
		}
		if value != expectedValues[counter] {
			t.Fatalf("Failed as value [%s] is not equal to expected value [%s].", value, expectedValues[counter])
		}
		counter = counter + 1
	}

}

func TestSyncAWSCredentialKeyFromCloud_SecretNameAvailable(t *testing.T) {
	service := service.MockAwsService{}
	awsProvider := cloud.AwsProvider{
		Service:    &service,
		Region:     "ap-southeast-1",
		SecretName: "prod/profile",
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	environmentVariableString, err := SyncCredentialKeyFromCloud(awsProvider, keyValueEnvMap)
	if err != nil {
		t.Fatal("Failed as it could not sync any credentials from the cloud provider")
	}
	if *environmentVariableString == "" {
		t.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
	}
}

func TestSyncAWSCredentialKeyFromCloud_SecretNameNotAvailable(t *testing.T) {
	service := service.MockAwsService{}
	awsProvider := cloud.AwsProvider{
		Service:    &service,
		Region:     "ap-southeast-1",
		SecretName: "prod/customer",
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	_, err := SyncCredentialKeyFromCloud(awsProvider, keyValueEnvMap)
	if err == nil {
		t.Fatal("Failed as it could not handle in case of the secret name is not available")
	}
}

func TestSyncAzureCredentialKeyFromCloud_SecretNameAvailable(t *testing.T) {
	service := service.MockAzureService{
		IsFail: false,
	}
	azureProvider := cloud.AzureProvider{
		Service:   &service,
		VaultName: "mock_vault",
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	environmentVariableString, err := SyncCredentialKeyFromCloud(azureProvider, keyValueEnvMap)
	if err != nil {
		t.Fatal("Failed as it could not sync any credentials from the cloud provider")
	}
	if *environmentVariableString == "" {
		t.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
	}
}

func TestSyncAzureCredentialKeyFromCloud_SecretNameNotAvailable(t *testing.T) {
	service := service.MockAzureService{
		IsFail: true,
	}
	azureProvider := cloud.AzureProvider{
		Service:   &service,
		VaultName: "mock_vault",
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	_, err := SyncCredentialKeyFromCloud(azureProvider, keyValueEnvMap)
	if err == nil {
		t.Fatal("Failed as it could not handle in case of the secret name is not available")
	}
}

func TestSyncGoogleCloudCredentialKeyFromCloud_SecretNameAvailable(t *testing.T) {
	service := service.MockGoogleCloudService{
		IsFail: false,
	}
	cloudProvider := cloud.GoogleCloudProvider{
		ProjectId: "mock-id",
		Service:   &service,
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	environmentVariableString, err := SyncCredentialKeyFromCloud(cloudProvider, keyValueEnvMap)
	if err != nil {
		t.Fatal("Failed as it could not sync any credentials from the cloud provider")
	}
	if *environmentVariableString == "" {
		t.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
	}
}

func TestSyncGoogleCredentialKeyFromCloud_SecretNameNotAvailable(t *testing.T) {
	service := service.MockGoogleCloudService{
		IsFail: true,
	}
	cloudProvider := cloud.GoogleCloudProvider{
		ProjectId: "mock-id",
		Service:   &service,
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	_, err := SyncCredentialKeyFromCloud(cloudProvider, keyValueEnvMap)
	if err == nil {
		t.Fatal("Failed as it could not handle in case of the secret name is not available")
	}
}

func TestSyncVaultCredentialKeyFromCloud_SecretNameAvailable(t *testing.T) {
	service := service.MockVaultService{}
	vaultProvider := cloud.VaultProvider{
		Service: &service,
		Role:    "users-dev",
		Path:    "kv/data/backend/dev/users",
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	environmentVariableString, err := SyncCredentialKeyFromCloud(vaultProvider, keyValueEnvMap)
	if err != nil {
		t.Fatal("Failed as it could not sync any credentials from the cloud provider")
	}
	if *environmentVariableString == "" {
		t.Fatal("Failed as it could not map local environment variables with the credentials from the cloud provider")
	}
}

func TestSyncVaultCredentialKeyFromCloud_SecretNameNotAvailable(t *testing.T) {
	service := service.MockVaultService{}
	vaultProvider := cloud.VaultProvider{
		Service: &service,
		Role:    "users-dev",
		Path:    "kv/data/backend/dev/app",
	}
	keyValueEnvMap := LoadCredentialKeyFromEnvironment()
	_, err := SyncCredentialKeyFromCloud(vaultProvider, keyValueEnvMap)
	if err == nil {
		t.Fatal("Failed as it could not handle in case of the secret name is not available")
	}
}
