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
	re := regexp.MustCompile(`\W+`)
	var err error
	switch cloudProvider.GetName() {
	case "aws":
		credentialData, err = cloudSession.RetrieveCredentials()
	case "azure":
		var result map[string]string
		for _, localValue := range credentialKey {
			pureLocalValue := re.ReplaceAllString(localValue, "")
			os.Setenv("AZ_SECRET_NAME", pureLocalValue)
			result, err = cloudSession.RetrieveCredentials()
			if err != nil {
				errorMessage := fmt.Sprintf("%v", err)
				return nil, errors.New(errorMessage)
			}
			if result[pureLocalValue] != "" {
				credentialData[pureLocalValue] = result[pureLocalValue]
			}
		}
	}

	if err != nil {
		errorMessage := fmt.Sprintf("%v", err)
		return nil, errors.New(errorMessage)
	}
	environmentVariableString := ""
	for key, value := range credentialKey {
		pureValue := re.ReplaceAllString(value, "")
		environmentVariableString = environmentVariableString + fmt.Sprintf("export %s=%s\n", key, credentialData[pureValue])
	}
	fmt.Print(environmentVariableString)
	return &environmentVariableString, nil
}
