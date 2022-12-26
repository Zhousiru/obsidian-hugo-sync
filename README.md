<div align="center">

# Obsidian Hugo Sync
同步 S3-compatible 对象存储上的 Obsidian Vault 到 Hugo 站点

⚠️ Unstable

</div>

## Features

- 同步 Obsidian 资源文件夹，上传到资源 Bucket。其中，受支持的图片会被转换为 WebP 格式
- 同步文章
- 文章转换
  - Front Matter 处理
  - 转换 Wikilink
  - 图片链接更改到资源 Bucket

## 前置条件

- Obsidian Vault 在 S3-compatible 对象存储上<br />
  例如使用 [Remotely Save](https://github.com/remotely-save/remotely-save) 的 S3 同步方式
- 一个公有读 S3 Bucket 专用于存储图片
- 一个可以运行 Obsidian Hugo Sync 的环境<br />
  例如 GitHub Actions / VPS
- Hugo 开启 HTML 支持

## Setup

1. 安装 `libwebp` 以及命令行工具，确保 `cwebp` 以及 `gif2webp` 可用<br />
   ```
   sudo apt install -y webp
   ```
2. 从源码编译
   ```
   git clone https://github.com/Zhousiru/obsidian-hugo-sync.git
   ```
   ```
   cd obsidian-hugo-sync && go build
   ```

3. 重命名 `config.json.sample` 到 `config.json`，填入配置

4. 加入 Cron Job，定时运行

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
| vaultPost  | 仓库 Bucket 中的文章文件夹（以 `/` 结尾）<br />只有此文件夹下的文章会被同步和处理 |
| vaultAsset | 仓库 Bucket 中的资源文件夹（以 `/` 结尾）<br />此文件夹下的文件会被同步到资源 Bucket<br />受支持的图片将会转换为 WebP 格式 |
| assetUrl | 访问资源 Bucket 的 URL（以 `/` 结尾），用于替换文章中图片链接 |
| assetCacheControl | 资源 Bucket 中图片等资源的缓存策略<br />如 `public, max-age=31536000` |
| hugo.sitePath | Hugo 站点路径 |
| hugo.postPath | Hugo 存放文章的路径 |
| hugo.cmd | Hugo 构建命令，如不需要构建则留空 |

## Front Matter 处理

在 Obsidian 文章的 Front Matter 中：

- `title` 中的 `{{auto}}` 会自动替换为 `title: <Obsidian 文章标题>`
- `{{asset feature-image.jpg}}` 会被替换为 `https://<资源 Bucket URL>/feature-image.webp`

## Obsidian Markdown 语法拓展

*Based on RegEx, not AST 少数情况下会误伤 :D*

### 图片
- 使用 CommonMark：`![...](<URL>)`<br />
  指定大小：`![<Size>](...)`<br />
  指定替代文本：`![alt <Alt>](...)`<br />
  同时指定大小和替代文本：`![alt <Alt>|<Size>](...)`

- 使用 Wikilink：`![[<URL>]]`<br />
  指定大小：`![[<URL>|<Size>]]`<br />
  指定替代文本：`![[<URL>|alt <Alt>]]`<br />
  同时指定大小和替代文本：`![[<URL>|alt <Alt>|<Size>]]`

- 说明
  - `<URL>` 可为 `<资源文件夹>/<图片文件名>` 或 `<图片文件名>`
  - `<Size>` 可为 `<Width>x<Height>` 或 `<Width>`
  - 使用网络图片：图片 URL 以 `http(s)://` 开头即可

### 内部文章链接

- 使用 CommonMark：`![](<URL>)`<br />
  指定显示文本：`![<显示文本>](<URL>)`

- 使用 Wikilink：`[[<URL>]]`<br />
  指定显示文本：`![[<URL>|<显示文本>]]`

- 说明
  - `<URL>` 可为 `<文章文件夹>/<文章标题>` 或 `<文章标题>`
  - 链接会在当前页面打开

### 外部链接

- 使用 CommonMark：`![](<URL>)`<br />
  指定显示文本：`![<显示文本>](<URL>)`

- 说明
  - `<URL>` 需要以 `http(s)://` 开头
  - 链接会在新页面打开

## 关于 WebP 转换

`vaultAsset` 下的所有文件都会被上传到 Asset Bucket，其中拓展名为 `png`, `jpg`, `jpeg`, `tiff`, `tif`, `gif` 的图片会被转换为 WebP 格式

同时，引入这些图片时的 URL 会被修改

## 关于同步

使用两个文件分别记录摘要、文件名的关系

格式为：

```
md5(<Filename of Raw File> + <ETag>)|<Filename of Raw File>|<Filename of Processed File>
```

例如：

```
470a2f67b129a34e09bacf25d20e5a72|test-image.jpg|test-image.webp
```

```
32832606886ee83d907a82e8e63b85a0|test-post.md|test-post.md
```

通过检测键的增删，生成或删除对应的文件，达到同步的目的

只有 `data/post_mapping` 中描述的文章会受影响，因此可以方便地与其他发布方式（如 Git）协同使用
