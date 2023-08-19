package config

import (
	"os"

	"github.com/MohamedGouaouri/envault/util"
)

var VAULT_ADDR string
var VAULT_TOKEN, _ = os.LookupEnv("VAULT_TOKEN")

func init() {
	// Set VAULT_ADDR to its default value if it's not found
	vaultAddr, found := os.LookupEnv("VAULT_ADDR")
	if found {
		VAULT_ADDR = vaultAddr
	}
	VAULT_ADDR = util.VAULT_ADDR_DEFAULT
}
