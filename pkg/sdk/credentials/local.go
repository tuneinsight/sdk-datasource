package credentials

import "fmt"

// LocalProviderType is the type of Local.
const LocalProviderType ProviderType = "local"

// Local is a collection of locally stored credentials.
type Local struct {
	credentials map[string]Credentials
}

// NewLocal returns a new instance of Local initialized with @credentials.
// @credentials are deep copied into the returned Local instance.
func NewLocal(credentials map[string]Credentials) *Local {
	localCred := new(Local)
	localCred.credentials = make(map[string]Credentials)
	for k, v := range credentials {
		localCred.credentials[k] = v
	}
	return localCred
}

// LocalFactory is the ProviderFactory for Local. It accepts the same arguments as NewLocal.
func LocalFactory(args ...interface{}) (Provider, error) {
	if args == nil {
		return nil, nil
	}
	credentials, ok := args[0].(map[string]Credentials)
	if !ok {
		return nil, fmt.Errorf("wrong input type for local factory: expected: %T, actual: %T", map[string]Credentials{}, args[0])
	}
	return NewLocal(credentials), nil
}

// Type returns the type of Local.
func (local Local) Type() ProviderType {
	return LocalProviderType
}

// GetCredentials returns the credentials, stored in @local, identified by @credID.
// If no credentials are stored under the provided @credID, it returns an error.
func (local Local) GetCredentials(credID string) (cred Credentials, err error) {
	cred, ok := local.credentials[credID]
	if !ok {
		err = fmt.Errorf("no credentials found for ID: %s", credID)
	}
	return
}

// SetCredentials stores @cred in @local under @credID.
// It silently overwrites any credentials previously stored under @credID.
func (local *Local) SetCredentials(credID string, cred Credentials) {
	local.credentials[credID] = cred
}
