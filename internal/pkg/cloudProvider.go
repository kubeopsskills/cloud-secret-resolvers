package pkg

type cloudProvider interface {
	initialCloudSession() bool
	retrieveCredentials() map[string]string
}
