package main

import (
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/csr"
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/pkg/provider"
)

func main() {
	// load environment configs
	var keyValueEnvMap = csr.LoadCredentialKeyFromEnvironment()
	var cloudType = "aws"
	switch {
	case cloudType == "aws":
		awsProvider := provider.AwsProvider{Region: "ap-southeast-1", SecretName: "prod/test"}
		csr.SyncCredentialKeyFromCloud(awsProvider, keyValueEnvMap)
	}

}
