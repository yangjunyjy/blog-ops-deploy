# 图片上传 API 文档

## 概述

图片上传服务支持本地存储和远程存储（如阿里云OSS、腾讯云COS、七牛云等）两种方式，可以在配置文件中灵活切换。

## 配置说明

### 配置文件 (config.yaml)

```yaml
upload:
  type: local                # 上传方式: local 或 remote
  urlPrefix: "http://localhost:8081"  # URL前缀，返回前端时自动拼接完整URL

  # 本地上传配置
  local:
    uploadPath: "./uploads"          # 上传目录
    allowedExts:                      # 允许的文件扩展名
      - ".jpg"
      - ".jpeg"
      - ".png"
      - ".gif"
      - ".bmp"
      - ".webp"
      - ".svg"
    maxSize: 10485760                # 最大文件大小(10MB)

  # 远程上传配置
  remote:
    provider: "oss"                 # 服务提供商: oss, cos, qiniu
    endpoint: ""                     # 服务端点
    accessKey: ""                    # 访问密钥
    secretKey: ""                    # 密钥
    bucket: ""                       # 存储桶名称
    region: ""                       # 区域
    cdn: ""                          # CDN域名
```

### 环境变量

也可以通过环境变量配置：

```
UPLOAD_TYPE=local
UPLOAD_URL_PREFIX=http://localhost:8081
UPLOAD_LOCAL_UPLOAD_PATH=./uploads
UPLOAD_LOCAL_MAX_SIZE=10485760
UPLOAD_REMOTE_PROVIDER=oss
UPLOAD_REMOTE_ENDPOINT=
UPLOAD_REMOTE_ACCESS_KEY=
UPLOAD_REMOTE_SECRET_KEY=
UPLOAD_REMOTE_BUCKET=
UPLOAD_REMOTE_REGION=
UPLOAD_REMOTE_CDN=
```

## API 接口

### 上传图片

**接口地址:** `POST /api/v1/admin/upload/image`

**请求头:**
```
Content-Type: multipart/form-data
Authorization: Bearer {token}
```

**请求参数:**
- `file` (file, 必填): 图片文件

**支持的图片格式:**
- JPG/JPEG
- PNG
- GIF
- BMP
- WEBP
- SVG

**文件大小限制:** 最大 10MB

**响应示例:**

成功 (200):
```json
{
  "code": 0,
  "message": "图片上传成功",
  "data": {
    "url": "http://localhost:8081/uploads/images/a1b2c3d4_e5f6g7h8.jpg",
    "name": "example.jpg",
    "size": 123456
  }
}
```

失败 (400):
```json
{
  "code": 400,
  "message": "获取文件失败",
  "error": "..."
}
```

失败 (400):
```json
{
  "code": 400,
  "message": "不支持的文件类型: .txt",
  "error": "..."
}
```

失败 (400):
```json
{
  "code": 400,
  "message": "文件大小超过限制: 最大 10 MB",
  "error": "..."
}
```

## 上传方式说明

### 本地上传 (local)

- 文件存储在服务器本地 `uploadPath` 指定的目录下
- 图片存放在 `uploadPath/images/` 子目录中
- 返回的URL会自动拼接 `urlPrefix` 前缀
- 适用于小型项目或内网部署

### 远程上传 (remote)

支持多种云存储服务：

#### 阿里云 OSS
- `provider`: `oss`
- 需要安装阿里云 OSS SDK: `go get github.com/aliyun/aliyun-oss-go-sdk/oss`
- 配置 `endpoint`, `accessKey`, `secretKey`, `bucket`, `region`

#### 腾讯云 COS
- `provider`: `cos`
- 需要安装腾讯云 COS SDK: `go get github.com/tencentyun/cos-go-sdk-v5`
- 配置 `endpoint`, `secretId`, `secretKey`, `bucket`, `region`

#### 七牛云
- `provider`: `qiniu`
- 需要安装七牛云 SDK: `go get github.com/qiniu/go-sdk/v7`
- 配置 `accessKey`, `secretKey`, `bucket`

**注意:** 远程上传功能需要在 `internal/services/upload_service.go` 中对应的方法中实现完整的SDK调用代码。

## 文件命名规则

上传后的图片文件名格式：
```
{MD5前8位}_{UUID前8位}{原始扩展名}
```

例如：
```
a1b2c3d4_e5f6g7h8.jpg
```

这样可以避免文件名冲突，同时通过MD5前缀可以快速识别重复文件。

## URL 拼接规则

返回的完整URL格式：
```
{urlPrefix}{relativePath}
```

示例：
- `urlPrefix`: `http://localhost:8081`
- `relativePath`: `/uploads/images/a1b2c3d4_e5f6g7h8.jpg`
- 完整URL: `http://localhost:8081/uploads/images/a1b2c3d4_e5f6g7h8.jpg`

## 前端使用示例

```javascript
// 使用 FormData 上传图片
const formData = new FormData();
formData.append('file', fileInput.files[0]);

fetch('http://localhost:8081/api/v1/admin/upload/image', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + token
  },
  body: formData
})
.then(response => response.json())
.then(data => {
  if (data.code === 0) {
    console.log('上传成功:', data.data.url);
    // 使用返回的图片URL
  } else {
    console.error('上传失败:', data.message);
  }
})
.catch(error => {
  console.error('请求失败:', error);
});
```

## 安全说明

1. **文件类型验证:** 只允许上传配置的图片格式
2. **文件大小限制:** 默认最大 10MB
3. **MD5 校验:** 上传时计算文件MD5，可用于去重
4. **唯一文件名:** 避免文件名冲突
5. **访问控制:** 接口需要认证（`admin` 路由组）

## 错误码说明

| HTTP状态码 | 说明 |
|-----------|------|
| 200 | 上传成功 |
| 400 | 请求参数错误（文件获取失败、类型不支持、大小超限） |
| 500 | 服务器内部错误（文件读写失败、目录创建失败等） |

## 静态文件访问

如果使用本地上传，需要在路由中配置静态文件访问，以便访问上传的图片：

```go
// 在 router.go 中添加
router.Static("/uploads", "./uploads")
```

或者通过Nginx/Apache配置反向代理访问静态文件。
