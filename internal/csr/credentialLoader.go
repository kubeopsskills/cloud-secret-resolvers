package csr

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/provider"
)

func LoadCredentialKeyFromEnvironment() map[string]string {
	var keyValueEnvMap map[string]string = make(map[string]string)
	for _, env := range os.Environ() {
		keyValueEnv := strings.SplitN(env, "=", 2)
		credentialPattern := regexp.MustCompile(`\${.*}`)
		if credentialPattern.MatchString(keyValueEnv[1]) {
			keyValueEnvMap[keyValueEnv[0]] = keyValueEnv[1]
		}
	}
	return keyValueEnvMap
}

func SyncCredentialKeyFromCloud(cloudProvider provider.CloudProvider, credentialKey map[string]string) (*string, error) {
	cloudSession := cloudProvider.InitialCloudSession()
	var credentialData = make(map[string]string)

	var err error
	switch cloudProvider.GetName() {
	case "aws":
		credentialData, err = cloudSession.RetrieveCredentials()
	case "azure":
		credentialData, err = getSecretWithLocalKey("AZ_SECRET_NAME", credentialKey, cloudSession)
	case "gcloud":
		credentialData, err = getSecretWithLocalKey("GC_SECRET_NAME", credentialKey, cloudSession)
	}

	if err != nil {
		errorMessage := fmt.Sprintf("%v", err)
		return nil, errors.New(errorMessage)
	}
	re := regexp.MustCompile(`\W+`)
	environmentVariableString := ""
	for key, value := range credentialKey {
		environmentVariableString = environmentVariableString + fmt.Sprintf("export %s=%s\n", key, credentialData[re.ReplaceAllString(value, "")])
	}
	fmt.Print(environmentVariableString)
	return &environmentVariableString, nil
}

func getSecretWithLocalKey(
	secretKey string,
	credentialKey map[string]string,
	cloudSession provider.CloudProvider,
) (map[string]string, error) {
	var credentialData = make(map[string]string)
	var result map[string]string
	var err error

	re := regexp.MustCompile(`\W+`)

	for _, localValue := range credentialKey {
		pureLocalValue := re.ReplaceAllString(localValue, "")

		os.Setenv(secretKey, pureLocalValue)
		result, err = cloudSession.RetrieveCredentials()
		if err != nil {
			errorMessage := fmt.Sprintf("%v", err)
			return nil, errors.New(errorMessage)
		}
		if result[pureLocalValue] != "" {
			credentialData[pureLocalValue] = result[pureLocalValue]
		}
	}
	return credentialData, nil
}
