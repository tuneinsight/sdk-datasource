package credentials

// Credentials contains the credentials to access a resource.
type Credentials struct {
	username string
	password string

	connectionString string
}

// NewCredentials returns a new instance of Credentials containing the passed parameters.
func NewCredentials(username, password, connectionString string) *Credentials {
	return &Credentials{
		username:         username,
		password:         password,
		connectionString: connectionString,
	}
}

// Username returns the username linked to a Credentials instance.
// If no username is linked to the credentials, an empty string is returned.
func (cred Credentials) Username() string {
	return cred.username
}

// Password returns the password linked to a Credentials instance.
// If no password is linked to the credentials, an empty string is returned.
func (cred Credentials) Password() string {
	return cred.password
}

// ConnectionString returns the connection string linked to a Credentials instance.
// If no connection string is linked to the credentials, an empty string is returned.
func (cred Credentials) ConnectionString() string {
	return cred.connectionString
}
