> title名称

<view class="api-title">域名Whois信息查询</view>

> api描述

<view class="api-desc">域名Whois详细信息查询</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/whoisquery</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/whoisquery?domain=devopsclub.cn&type=json&standard=true
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
domain | String | 是 | 域名
type | String | 否 | whois数据返回类型(text 文本串/json json格式数据)
standard | Bool | 否 | 是否按照标准固定格式输出json字段，默认是按原whois信息中的所有字段返回，只对返回json格式有效

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
status | Int | 域名查询状态(0 获取到域名whois信息/1 域名解析失败/2 域名未注册/3 暂不支持此域名后缀查询/4 域名查询失败/5 请求数据错误)
data | Map/String | 域名whois详细数据
msg | String | 消息

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "data": {
            "contactEmail": "aery_mzc9123@163.com",
            "contactPhone": "",
            "dnsNameServer": [
                "dns15.hichina.com",
                "dns16.hichina.com"
            ],
            "domainName": "devopsclub.cn",
            "domainStatus": "ok",
            "expirationTime": "2024-05-18 22:00:04",
            "registrant": "马智超",
            "registrar": "阿里云计算有限公司（万网）",
            "registrarWHOISServer": "",
            "registrationTime": "2019-05-18 22:00:04",
            "updatedDate": ""
        },
        "status": 0
    },
    "msg": ""
}
```

</view>

> 错误码参照

<view class="error-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 请求返回状态码(0 请求成功, 1 请求失败)

</view>

> 示例代码

<view class="code-demo">

```text
暂无示例代码, 问题反馈qq: 1152490990
```

</view>