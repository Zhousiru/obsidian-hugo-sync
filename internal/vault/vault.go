package post

import (
	"github.com/Zhousiru/obsidian-hugo-sync/internal/config"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/s3"
)

var client *s3.Client = nil

// ListVaultPost returns a slice of S3 Obsidian vault post objects.
func ListVaultPost() ([]*s3.Object, error) {
	return listVault(config.X.VaultPost)
}

// ListVaultAsset returns a slice of S3 Obsidian vault asset objects.
func ListVaultAsset() ([]*s3.Object, error) {
	return listVault(config.X.VaultAsset)
}

func listVault(prefix string) ([]*s3.Object, error) {
	if client == nil {
		// setup vault client
		client = new(s3.Client)
		err := client.Setup(config.X.S3.Vault.Endpoint,
			config.X.S3.Vault.Region,
			config.X.S3.Vault.AccessKeyID,
			config.X.S3.Vault.SecretAccessKey,
			config.X.S3.Vault.Bucket)
		if err != nil {
			return nil, err
		}
	}

	return client.List(prefix)
}
