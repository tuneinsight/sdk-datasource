package credentials

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAzureKeyVault(t *testing.T) {

	// credentials of the TI test AKV
	os.Setenv("AZURE_TENANT_ID", "e6021d6c-8bdc-4c91-b88f-e3333caae8b8")
	os.Setenv("AZURE_CLIENT_ID", "00d75d9b-9524-4428-bd6a-a6dc2a08ac19")
	os.Setenv("AZURE_CLIENT_SECRET", "P7I8Q~KpPiCphM2LOe25Wil6c80vHo3KJqSr3b~r")
	os.Setenv("AZURE_KEY_VAULT_URI", "https://ti-test-vault.vault.azure.net/")

	// these are test credentials previously stored on the TI test AKV
	testCredentialsID := "test-credentials"
	testCredentials := &Credentials{
		username: "test-user",
		password: "test-password",
	}

	akv, err := NewAzureKeyVault(nil)
	require.NoError(t, err)

	// cannot retrieve not existing credentials
	_, err = akv.GetCredentials("not-existing-credentials")
	require.Error(t, err)

	// correctly retrieving existing credentials
	creds, err := akv.GetCredentials(testCredentialsID)
	require.NoError(t, err)
	require.EqualValues(t, testCredentials, creds)

}
