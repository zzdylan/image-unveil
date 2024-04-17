# 图像分割处理服务
  
这个项目使用 Alibaba Cloud 的图像分割 API 来自动处理图片，将图片中的主体对象从背景中分割出来。这对于图像编辑、增强现实应用或任何需要图像预处理的场景都非常有用。
  
## 开始
  
以下指南将帮助你在本地机器上部署和运行此项目，用于开发和测试目的。
  
### 先决条件
  
在开始之前，确保你的开发环境满足以下要求：
  
- Go 语言环境（1.16 或更高版本）
- git（用于克隆仓库）
  
  # 检查 Go 版本


```
go  version
```

  
# 应该显示 Go version go1.16 或更高
  
  
# 克隆仓库
```
git  clone  https://github.com/zzdylan/image-unveil.git
```
  
# 进入项目目录
```
cd  image-unveil
```
  
# 安装依赖
```
go  mod  tidy
```

# 复制.env.example为.env
```
cp .env.example .env
```
  
# 设置环境变量
```
vim .env
```

# 开通服务
使用前请务必在https://vision.aliyun.com/imageseg进行开通，0.0020元/次
  
# 设置阿里云的AccessKeyId和AccessKeySecret
```
ACCESS_KEY_ID=your_access_key_id
ACCESS_KEY_SECRET=your_access_key_secret
```
    
`
  
# 运行程序
```
go  run  main.go
```
  
  
# 构建与部署
```
#linux
go build -o app
#windows
go build -o app.exe
```