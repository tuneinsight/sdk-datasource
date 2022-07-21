package credentials

// ProviderType is the Provider type.
type ProviderType string

// ProviderFactory is a function that instantiates a Provider.
type ProviderFactory func(...interface{}) (Provider, error)

// Provider is the interface that all credentials providers must implement.
type Provider interface {
	// Type returns the Provider type.
	Type() ProviderType
	// GetCredentials returns the Credentials identified by @credID.
	GetCredentials(credID string) (*Credentials, error)
}

// ProviderFactories contains the ProviderFactory for the implemented Provider s.
var ProviderFactories = map[ProviderType]ProviderFactory{
	LocalProviderType:         LocalFactory,
	AzureKeyVaultProviderType: AzureKeyVaultFactory,
}
