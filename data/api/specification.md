### 格式文档书写规范

```json5
{
   "whoisquery": {                     // api名称，要与配置的路由名称一致
    "docFile": "whois.json",           // docs数据存储文件
    "countKey": "whoisquery",          // redis统计请求次数key名
    "description": "查询域名whois信息",  // 描述信息
    "titleName": "域名Whois查询",       // api标题
    "enable": true,                    // 是否启用
    "mainten": false                   // 是否维护中
   }
}
```