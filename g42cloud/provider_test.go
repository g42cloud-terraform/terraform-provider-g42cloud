package g42cloud

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/pathorcontents"
)

var (
	G42_REGION_NAME                = os.Getenv("G42_REGION_NAME")
	G42_ACCOUNT_NAME               = os.Getenv("G42_ACCOUNT_NAME")
	G42_ADMIN                      = os.Getenv("G42_ADMIN")
	G42_ENTERPRISE_PROJECT_ID_TEST = os.Getenv("G42_ENTERPRISE_PROJECT_ID_TEST")
	G42_ACCESS_KEY                 = os.Getenv("G42_ACCESS_KEY")
	G42_SECRET_KEY                 = os.Getenv("G42_SECRET_KEY")
	G42_DNS_ENVIRONMENT            = os.Getenv("G42_DNS_ENVIRONMENT")
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"g42cloud": testAccProvider,
	}
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	if G42_REGION_NAME == "" {
		t.Fatal("G42_REGION_NAME must be set for acceptance tests")
	}
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
}

func testAccPreCheckAdminOnly(t *testing.T) {
	if G42_ADMIN == "" {
		t.Skip("G42_ADMIN must be set for acceptance tests")
	}
}

func testAccPreCheckEpsID(t *testing.T) {
	if G42_ENTERPRISE_PROJECT_ID_TEST == "" {
		t.Skip("This environment does not support EPS_ID tests")
	}
}

func testAccPreCheckOBS(t *testing.T) {
	if G42_ACCESS_KEY == "" || G42_SECRET_KEY == "" {
		t.Skip("G42_ACCESS_KEY and G42_SECRET_KEY must be set for OBS acceptance tests")
	}
}

func testAccPreCheckDNS(t *testing.T) {
	if G42_DNS_ENVIRONMENT == "" {
		t.Skip("This environment does not support DNS tests")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func envVarContents(varName string) (string, error) {
	contents, _, err := pathorcontents.Read(os.Getenv(varName))
	if err != nil {
		return "", fmt.Errorf("Error reading %s: %s", varName, err)
	}
	return contents, nil
}

func envVarFile(varName string) (string, error) {
	contents, err := envVarContents(varName)
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", varName)
	if err != nil {
		return "", fmt.Errorf("Error creating temp file: %s", err)
	}
	if _, err := tmpFile.Write([]byte(contents)); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error writing temp file: %s", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Error closing temp file: %s", err)
	}
	return tmpFile.Name(), nil
}
