package provider

type CloudProvider interface {
	InitialCloudSession() CloudProvider
	RetrieveCredentials() (map[string]string, error)
}
