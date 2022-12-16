package main

import (
	"os"
	"sync"

	"github.com/Zhousiru/obsidian-hugo-sync/internal/assetconv"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/config"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/logger"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/mapping"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/postprocess"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/s3"
	"github.com/Zhousiru/obsidian-hugo-sync/internal/util"
)

func main() {
	vaultClient := new(s3.Client)
	err := vaultClient.Setup(
		config.X.S3.Vault.Endpoint,
		config.X.S3.Vault.Region,
		config.X.S3.Vault.AccessKeyID,
		config.X.S3.Vault.SecretAccessKey,
		config.X.S3.Vault.Bucket,
	)
	if err != nil {
		logger.Fatal("Failed to set up vault bucket client: %s", err)
	}

	assetClient := new(s3.Client)
	err = assetClient.Setup(
		config.X.S3.Asset.Endpoint,
		config.X.S3.Asset.Region,
		config.X.S3.Asset.AccessKeyID,
		config.X.S3.Asset.SecretAccessKey,
		config.X.S3.Asset.Bucket,
	)
	if err != nil {
		logger.Fatal("Failed to set up asset bucket client: %s", err)
	}

	logger.Info("Load present asset mapping from local")

	presentAsset := new(mapping.Mapping)
	presentAsset.SetType(mapping.AssetMapping)
	err = presentAsset.Load()
	if err != nil {
		logger.Fatal("Failed to load present asset mapping: %s", err)
	}

	logger.Info("Dump latest asset mapping from S3")

	assetObjSlice, err := vaultClient.List(config.X.VaultAsset)
	if err != nil {
		logger.Fatal("Failed to dump latest asset mapping: %s", err)
	}

	latestAsset := mapping.VaultToMapping(assetObjSlice, config.X.VaultAsset, mapping.AssetMapping)

	logger.Info("Update assets bucket")

	assetAdd, assetDel := presentAsset.Diff(latestAsset)
	var wg sync.WaitGroup

	for _, ent := range assetDel {
		logger.DelFile(ent.RawFilename, ent.Hash)
		wg.Add(1)

		go func(ent *mapping.Entry) {

			defer wg.Done()

			err := assetClient.Remove(ent.ProcessedFilename)
			if err != nil {
				logger.Err("Failed to remove %s: %s", ent.RawFilename, err)
				return
			}

			presentAsset.Remove(ent.Hash)

		}(ent)
	}

	wg.Wait()

	for _, ent := range assetAdd {
		logger.NewFile(ent.RawFilename, ent.Hash)
		wg.Add(1)

		go func(ent *mapping.Entry) {

			defer wg.Done()

			asset, err := vaultClient.Get(config.X.VaultAsset + ent.RawFilename)
			if err != nil {
				logger.Err("Failed to get asset %s: %s", ent.RawFilename, err)
				return
			}

			if assetconv.CanToWebP(util.GetExt(ent.RawFilename)) {
				asset, err = assetconv.ToWebP(asset)
				if err != nil {
					logger.Err("Failed to convert asset %s: %s", ent.RawFilename, err)
					return
				}
			}

			err = assetClient.Put(ent.ProcessedFilename, asset, util.GetContentType(util.GetExt(ent.ProcessedFilename)))
			if err != nil {
				logger.Err("Failed to upload asset %s: %s", ent.RawFilename, err)
				return
			}

			presentAsset.Append(ent)

		}(ent)
	}

	wg.Wait()

	err = presentAsset.Save()
	if err != nil {
		logger.Fatal("Failed to save asset mapping: %s", err)
	}

	presentPost := new(mapping.Mapping)
	presentPost.SetType(mapping.PostMapping)
	err = presentPost.Load()
	if err != nil {
		logger.Fatal("Failed to load present post mapping: %s", err)
	}

	logger.Info("Dump latest post mapping from S3")

	postObjSlice, err := vaultClient.List(config.X.VaultPost)
	if err != nil {
		logger.Fatal("Failed to dump latest post mapping: %s", err)
	}

	latestPost := mapping.VaultToMapping(postObjSlice, config.X.VaultPost, mapping.PostMapping)

	logger.Info("Update hugo posts")

	postAdd, postDel := presentPost.Diff(latestPost)

	for _, ent := range postDel {
		logger.DelFile(ent.RawFilename, ent.Hash)
		wg.Add(1)

		go func(ent *mapping.Entry) {

			defer wg.Done()

			err := os.Remove(config.X.HugoPostPath + ent.ProcessedFilename)
			if err != nil {
				logger.Err("Failed to remove post %s: %s", ent.RawFilename, err)
				return
			}

			presentPost.Remove(ent.Hash)

		}(ent)
	}

	wg.Wait()

	for _, ent := range postAdd {
		logger.NewFile(ent.RawFilename, ent.Hash)
		wg.Add(1)

		go func(ent *mapping.Entry) {

			defer wg.Done()

			post, err := vaultClient.Get(config.X.VaultPost + ent.RawFilename)
			if err != nil {
				logger.Err("Failed to get post %s: %s", ent.RawFilename, err)
				return
			}

			postStr := postprocess.Process(
				string(post),
				ent.RawFilename,
				config.X.AssetUrl,
				config.X.VaultAsset,
			)

			err = os.WriteFile(config.X.HugoPostPath+ent.ProcessedFilename, []byte(postStr), 0644)
			if err != nil {
				logger.Err("Failed to add post %s: %s", ent.RawFilename, err)
				return
			}

			presentPost.Append(ent)

		}(ent)
	}

	wg.Wait()

	err = presentPost.Save()
	if err != nil {
		logger.Fatal("Failed to save post mapping: %s", err)
	}

	// TODO: hugo build
}
