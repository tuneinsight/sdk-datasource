package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// AzureKeyVaultProviderType is the type of AzureKeyVault.
const AzureKeyVaultProviderType ProviderType = "azure-key-vault"

// AzureKeyVault is a collection of credentials stored in an Azure Key Vault (https://azure.microsoft.com/en-us/services/key-vault/#product-overview) instance.
type AzureKeyVault struct {
	client *azsecrets.Client
	// credsMapping maps the data source creds IDs with the secret IDs stored in the Azure Key Vault
	credsMapping map[string]string
}

// NewAzureKeyVault returns a new instance of AzureKeyVault.
// It requires the following env variables to be set: AZURE_TENANT_ID, AZURE_CLIENT_ID,
// AZURE_CLIENT_CERTIFICATE_PATH (or, in alternative, AZURE_CLIENT_SECRET), AZURE_KEY_VAULT_URI.
func NewAzureKeyVault(credsMapping map[string]string) (*AzureKeyVault, error) {

	azureKeyVault := new(AzureKeyVault)

	// check that all the env variables required by the Azure SDK are set
	tenantID := os.Getenv("AZURE_TENANT_ID")
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("AZURE_TENANT_ID is not set")
	}

	clientID := os.Getenv("AZURE_CLIENT_ID")
	if len(clientID) == 0 {
		return nil, fmt.Errorf("AZURE_CLIENT_ID is not set")
	}

	clientCertificatePath := os.Getenv("AZURE_CLIENT_CERTIFICATE_PATH")
	if len(clientCertificatePath) == 0 {
		clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
		if len(clientSecret) == 0 {
			return nil, fmt.Errorf("AZURE_CLIENT_CERTIFICATE_PATH and AZURE_CLIENT_SECRET are not set")
		}
	}

	keyVaultURI := os.Getenv("AZURE_KEY_VAULT_URI")
	if len(keyVaultURI) == 0 {
		return nil, fmt.Errorf("AZURE_KEY_VAULT_URI is not set")
	}

	// create a credential using the NewDefaultAzureCredential type
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain a credential: %v", err)
	}

	// exponential backoff options
	clientOptions := &azcore.ClientOptions{
		Retry: policy.RetryOptions{
			RetryDelay:    2 * time.Second,
			MaxRetryDelay: 16 * time.Second,
			MaxRetries:    5,
		},
	}
	// create a Key Vault client
	azureKeyVault.client = azsecrets.NewClient(keyVaultURI, cred, clientOptions)

	azureKeyVault.credsMapping = credsMapping

	return azureKeyVault, nil
}

// AzureKeyVaultFactory is the ProviderFactory for AzureKeyVault. It aceepts the same arguments as AzureKeyVault.
func AzureKeyVaultFactory(args ...interface{}) (Provider, error) {
	if args == nil {
		return NewAzureKeyVault(nil)
	}
	credsMapping, ok := args[0].(map[string]string)
	if !ok {
		return nil, fmt.Errorf("wrong input type for azure key vault factory: expected: %T, actual: %T", map[string]string{}, args[0])
	}
	return NewAzureKeyVault(credsMapping)
}

// Type returns the ProviderType of AzureKeyVault.
func (akv AzureKeyVault) Type() ProviderType {
	return AzureKeyVaultProviderType
}

// GetCredentials returns the credentials, stored in @akv, identified by @credID.
// If @credID is mapped to a different Azure Key Vault secret ID in akv.credsMapping,
// the function will use this last one to retrieve the credentials.
func (akv AzureKeyVault) GetCredentials(credID string) (*Credentials, error) {

	if id, ok := akv.credsMapping[credID]; ok {
		logrus.Infof("using mapped id -> %s: %s", credID, id)
		credID = id
	}

	resp, err := akv.client.GetSecret(context.TODO(), credID, "", nil)
	if err != nil {
		return nil, fmt.Errorf("while retrieving credentials with ID: %s from azure key vault: %s, %w", credID, os.Getenv("AZURE_KEY_VAULT_URI"), err)
	}

	credsJSON := *resp.Value
	creds := struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		ConnectionString string `json:"connectionString"`
	}{}
	if err := json.Unmarshal([]byte(credsJSON), &creds); err != nil {
		return nil, fmt.Errorf("while unmarhsalling retrieved credentials with ID: %s: %s", credID, credsJSON)
	}

	return NewCredentials(creds.Username, creds.Password, creds.ConnectionString), nil
}
