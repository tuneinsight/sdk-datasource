package credentials

import "fmt"

// LocalCredentials is a collection of locally stored credentials.
type LocalCredentials struct {
	credentials map[string]Credentials
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
