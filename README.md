#### ApiTools是一个集成各种常用工具，为开发者提供免费、便捷的API接口，方便开发者可以通过API方式快速获取数据

> QQ交流群: 565070560

### 开发

*技术栈*:
```bash
开发语言: go 1.13
web框架: gin 
关系型数据库: mysql 5.7
nosql数据库: redis 5.x
```

### 部署

> 需要具备go语言环境

1.编译代码

```text
$ git clone https://github.com/spdir/apiTools.git
$ go env -w GOPROXY=https://goproxy.io,direct
$ cd apiTools

> make编译
$ make clean       // 清理
$ make build       // 编译二进制

> 编译并打包
$ make pack     // 本地环境打包
$ make pack-linux-amd64   // linux amd64 环境交叉编译打包


> 手动编译(交叉编译)
$ GOOS=linux GOARCH=amd64 go build -o apiTools main.go
```

2.构建docker image

```text
1) 修改 `config` 目录下配置文件
或者不修改，通过file映射、docker ENV 去设置，程序会自动读取

2) 进行编译和打包
$ make pack-linux-amd64

3) 构建docker image
$ export apiTools_version=0.1
$ docker build -t apitools:$apiTools_version .
```

---
官网地址: [https://api.devopsclub.cn](https://api.devopsclub.cn)

