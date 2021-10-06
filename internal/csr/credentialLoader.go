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
		var result map[string]string
		for localKey := range credentialKey {
			os.Setenv("AZ_SECRET_NAME", localKey)
			result, err = cloudSession.RetrieveCredentials()
			if err != nil {
				errorMessage := fmt.Sprintf("%v", err)
				return nil, errors.New(errorMessage)
			}
			if result[localKey] != "" {
				credentialData[localKey] = result[localKey]
			}
		}
	}

	if err != nil {
		errorMessage := fmt.Sprintf("%v", err)
		return nil, errors.New(errorMessage)
	}
	re := regexp.MustCompile(`\W+`)
	environmentVariableString := ""
	for key, value := range credentialKey {
		environmentVariableString = environmentVariableString + fmt.Sprintf("export %s=%s\n", key, credentialData[re.ReplaceAllString(value, "")])
		fmt.Print(environmentVariableString)
	}
	return &environmentVariableString, nil
}
