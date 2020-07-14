## http proxy ip pool

- 依赖项目：
    `https://github.com/jhao104/proxy_pool/`
   
- docker 部署

```bash
$ cd data/proxy_pool
$ git clone https://github.com/jhao104/proxy_pool.git
$ cd proxy_pool
$ docker build -t proxy_pool:version .
$ docker rm -f proxy_pool
$ docker run -d --env DB_CONN=redis://:password@ip:port/db -p 15010:5010 --name proxy_pool proxy_pool:version
```

- 配置信息
```text
监听端口: 15010
redis: 宿主机: 6379
```
