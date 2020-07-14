### 这是什么
ApiTools是一个集成各种开放Api的web项目，为大家提供免费的、常用的Api功能

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

> 手动编译(交叉编译)
$ GOOS=linux GOARCH=amd64 go build -o apiTools main.go
```
2.发布

> 将编译好的二进制文件和配置文件和基本数据打包

```text
config/             // 配置文件目录需要打包 
data/               // 数据目录需要打包  
views/              // html文件
static/             // 静态文件   
apiTools            // 可执行程序文件
```

---
官网地址: [https://api.devopsclub.cn](https://api.devopsclub.cn)

