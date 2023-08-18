package config

import "os"

var VAULT_ADDR, _ = os.LookupEnv("VAULT_ADDR")
var VAULT_TOKEN, _ = os.LookupEnv("VAULT_TOKEN")
