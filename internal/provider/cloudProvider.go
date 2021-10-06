package provider

type CloudProvider interface {
	GetName() string
	InitialCloudSession() CloudProvider
	RetrieveCredentials() (map[string]string, error)
}
