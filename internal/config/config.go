package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type config struct {
	S3 struct {
		Vault struct {
			Endpoint        string `json:"endpoint"`
			Region          string `json:"region"`
			AccessKeyId     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
			Bucket          string `json:"bucket"`
		} `json:"vault"`
		Asset struct {
			Endpoint        string `json:"endpoint"`
			Region          string `json:"region"`
			AccessKeyId     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
			Bucket          string `json:"bucket"`
		} `json:"asset"`
	} `json:"s3"`
	VaultPost         string `json:"vaultPost"`
	VaultAsset        string `json:"vaultAsset"`
	AssetCacheControl string `json:"assetCacheControl"`
	AssetUrl          string `json:"assetUrl"`
	Hugo              struct {
		SitePath string   `json:"sitePath"`
		PostPath string   `json:"postPath"`
		Cmd      []string `json:"cmd"`
	} `json:"hugo"`
}

// X saves unmarshaled config data.
var X = new(config)

func init() {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(execPath), "config.json"))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, X)
	if err != nil {
		fmt.Println(err)
	}
}
