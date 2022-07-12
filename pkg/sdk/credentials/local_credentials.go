package credentials

import "fmt"

// LocalCredentials is a collection of locally stored credentials.
type LocalCredentials struct {
	credentials map[string]Credentials
}

// NewLocalCredentials returns a new instance of LocalCredentials initialized with @credentials.
// @credentials are deep copied into the returned LocalCredentials instance.
func NewLocalCredentials(credentials map[string]Credentials) *LocalCredentials {
	localCred := new(LocalCredentials)
	localCred.credentials = make(map[string]Credentials)
	for k, v := range credentials {
		localCred.credentials[k] = v
	}
	return localCred
}

// GetCredentials returns the credentials, stored in @localCred, identified by @credID.
// If no credentials are stored under the provided @credID, it returns an error.
func (localCred LocalCredentials) GetCredentials(credID string) (cred Credentials, err error) {
	cred, ok := localCred.credentials[credID]
	if !ok {
		err = fmt.Errorf("no credentials found for ID: %s", credID)
	}
	return
}

// SetCredentials stores @cred in @localCred under @credID.
// It silently overwrites any credentials previously stored under @credID.
func (localCred *LocalCredentials) SetCredentials(credID string, cred Credentials) {
	localCred.credentials[credID] = cred
}
