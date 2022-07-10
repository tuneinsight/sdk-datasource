package credentials

// Provider is the interface that all credentials providers must implement.
type Provider interface {
	// GetCredentials returns the Credentials identified by @credID.
	GetCredentials(credID string) (Credentials, error)
}
