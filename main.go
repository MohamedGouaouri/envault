package main

import (
	"github.com/MohamedGouaouri/envault/cmd"
)

func main() {
	cmd.Execute()
	// ctx := context.Background()

	// // prepare a client with the given base address
	// client, err := vault.New(
	// 	vault.WithAddress("http://127.0.0.1:8200"),
	// 	vault.WithRequestTimeout(30*time.Second),
	// )

	// if err != nil {
	// 	log.Fatalf("Error connectin %v", err)
	// }

	// // authenticate with a root token (insecure)
	// if err := client.SetToken("hvs.1ibNc9rOM6azW8ST2YyOaEHn"); err != nil {
	// 	log.Fatalf("Error authenticating %v", err)
	// }

	// // read the secret
	// s, err := client.Secrets.KvV2Read(ctx, "dev", vault.WithMountPath("secret"))
	// if err != nil {
	// 	log.Fatalf("Error reading %v", err)
	// }
	// log.Println("secret retrieved:", s.Data.Data["DB_URI"])

}
