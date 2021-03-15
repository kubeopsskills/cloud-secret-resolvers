package csr

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kubeopsskills/cloud-secret-resolvers/internal/pkg"
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

func SyncCredentialKeyFromCloud(cloudProvider pkg.CloudProvider, credentialKey map[string]string) {
	cloudSession := cloudProvider.InitialCloudSession()
	var credentialData = cloudSession.RetrieveCredentials()
	re := regexp.MustCompile(`\W+`)
	for key, value := range credentialKey {
		fmt.Printf("export %s=%s\n", key, credentialData[re.ReplaceAllString(value, "")])
	}
}
