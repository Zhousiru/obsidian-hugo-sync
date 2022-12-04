package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	S3 struct {
		Vault struct {
			Endpoint        string `json:"endpoint"`
			Region          string `json:"region"`
			AccessKeyID     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
			Bucket          string `json:"bucket"`
		} `json:"vault"`
		Asset struct {
			Endpoint        string `json:"endpoint"`
			Region          string `json:"region"`
			AccessKeyID     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
			Bucket          string `json:"bucket"`
		} `json:"asset"`
	} `json:"s3"`
	VaultPost  string `json:"vaultPost"`
	VaultAsset string `json:"vaultAsset"`
	AssetUrl   string `json:"assetUrl"`
	SitePath   string `json:"sitePath"`
	HugoCmd    string `json:"hugoCmd"`
}

// X saves unmarshaled config data.
var X = new(config)

func init() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, X)
	if err != nil {
		fmt.Println(err)
	}
}