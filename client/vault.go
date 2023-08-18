package client

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/MohamedGouaouri/envault/config"
	"github.com/hashicorp/vault-client-go"
)

func GetAllEnvironmentVariables(token string, env string, path string) (map[string]string, error) {
	var secrets = make(map[string]string)
	ctx := context.Background()
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress(config.VAULT_ADDR),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Connection %v\n", err)
	}

	// authenticate with a root token (insecure)
	if err := client.SetToken(config.VAULT_TOKEN); err != nil {
		log.Fatalf("Auth %v\n", err)
	}

	log.Default().Printf("Reading from path %v", filepath.Join("secret", path))
	// read the secrets
	s, err := client.Secrets.KvV2Read(ctx, path, vault.WithMountPath(env))
	if err != nil {
		log.Default().Printf("Reading vault %v\n", err)
		return secrets, nil
	}
	for k, v := range s.Data.Data {

		secrets[k] = v.(string)
	}
	return secrets, nil
}
