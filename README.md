# obsidian-hugo-sync
同步 S3-compatible 对象存储上的 Obsidian Vault 到 Hugo 站点

🚧 UNDER CONSTRUCTION 🚧

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

## Mapping

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

虽然目前对文章的处理不会改变文件名，但是为了保持一致性，还是两个文件名都写进去

通过检测键的增删，生成或删除对应的文件，达到同步的目的

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
| assetUrl | 访问资源 Bucket 的 URL，用于替换文章中图片 `src` |
| sitePath | Hugo 站点路径 |
| hugoCmd | Hugo 构建命令 |

## Obsidian 文章处理

### Front Matter

- 在 Obsidian 文章的 Front Matter 中使用 `title: {{auto}}`，会自动替换为 `title: <Obsidian 文章标题>`

### Obsidian Markdown 转换

Based on RegEx, not AST :D

**不支持 WikiLinks**

- 引入图片：`![](<资源文件夹>/<图片文件名>)` 或 `![](<图片文件名>)`<br />
  指定大小：`![<Width>x<Height>](...)`<br />
  指定替代文本：`![alt <替代文本>](...)`<br />
  使用网络图片：图片 URL 以 `http(s)://` 开头即可
  将会被转换为 HTML 标签

### 图片格式转换

`vaultAsset` 下的所有文件都会被上传到 Asset Bucket，其中拓展名为 `png`, `jpg`, `jpeg`, `tiff`, `tif` 的图片会被转换为 WebP 格式

同时，引入这些图片时的 URL 会被修改
