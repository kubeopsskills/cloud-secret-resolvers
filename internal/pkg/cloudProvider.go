package pkg

type CloudProvider interface {
	InitialCloudSession() CloudProvider
	RetrieveCredentials() map[string]string
}
