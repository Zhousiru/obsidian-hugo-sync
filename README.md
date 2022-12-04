# obsidian-hugo-sync
同步 S3-compatible 对象存储上的 Obsidian Vault 到 Hugo 站点

## 流程

1. 拉取仓库 Bucket 资源文件夹，Diff
2. 处理图片
   - 删除资源 Bucket 中过期图片
   - 转换为 WebP 格式
   - 上传到资源 Bucket 中
3. 拉取仓库 Bucket 文章文件夹，Diff
4. 生成 Hugo 适用的 Markdown 文件
   - 删除过期文章
   - 生成 Hugo Front Matter
   - 修改图片链接
5. Run `hugo`

## Diff

- 建立文件，记录对应关系：
  仓库 Bucket 资源文件夹中原图片的 `ETag` <--> 资源 Bucket 中处理过的图片的 `filename`
  仓库 Bucket 文章文件夹中原文章的 `ETag` <--> 站点目录下处理过的文章的 `filename`
- 当本次获取到的 `ETag` 多出：
  1. 处理此文件（文章 / 图片），记录处理后的 `filename`
  2. 新增此条对应关系
- 当本次获取到的 `ETag` 少于：
  1. 删除少的 `ETag` 对应的文件（文章 / 图片）
  2. 删除此条对应关系

## 配置

重命名 `config.json.sample` 到 `config.json`

| 项目                          | 描述                                 |
| :--------------------------- | :----------------------------------- |
| s3.vault.endpoint            | 仓库 Bucket 的 Endpoint               |
| s3.vault.region              | 仓库 Bucket 的 Region                 |
| s3.vault.accessKeyId         | 仓库 Bucket 的 Access Key ID          |
| s3.vault.secretAccessKey     | 仓库 Bucket 的 Secret Access Key      |
| s3.vault.bucket              | 仓库 Bucket 的 Bucket                 |
| s3.asset.endpoint            | 资源 Bucket 的 Endpoint               |
| s3.asset.region              | 资源 Bucket 的 Region                 |
| s3.asset.accessKeyId         | 资源 Bucket 的 Access Key ID          |
| s3.asset.secretAccessKey     | 资源 Bucket 的 Secret Access Key      |
| s3.asset.bucket              | 资源 Bucket 的 Bucket                 |
| vaultPost  | 仓库 Bucket 中的文章文件夹<br />只有此文件夹下的文章会被同步 |
| vaultAsset | 仓库 Bucket 中的资源文件夹<br />此文件夹下的图片会被转换为 WebP 格式<br />并同步到资源 Bucket |
| sitePath | Hugo 站点路径 |
| hugoCmd | Hugo 构建命令 |
