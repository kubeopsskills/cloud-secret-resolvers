package main

import (
	"github.com/kubeopsskills/cloud-secret-resolvers/internal/csr"
)

func main() {
	// load environment configs
	var keyValueEnvMap = csr.LoadCredentialKeyFromEnvironment()
}
