# github_ppva

Windows GitHub加速器

### 原理

获取GitHub距离本机最近cdn服务器ip，添加到本地hosts文件，跳过互联网dns解析，直接本地解析指向目标ip

### 使用

先安装go环境，需1.16版本及以上 git clone代码 记入项目根目录，运行:

```cmd
SET CGO_ENABLED=1
SET GOOS=windows
set GOARCH=amd64
set GO111MODULE=on
go build
```

会生成一个.exe的执行文件，运行即可
