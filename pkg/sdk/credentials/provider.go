package credentials

// ProviderType is the Provider type.
type ProviderType string

// Provider is the interface that all credentials providers must implement.
type Provider interface {
	// Type returns the Provider type.
	Type() ProviderType
	// GetCredentials returns the Credentials identified by @credID.
	GetCredentials(credID string) (*Credentials, error)
}
