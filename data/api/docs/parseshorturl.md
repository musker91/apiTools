> title名称

<view class="api-title">自定义短链接还原</view>

> api描述

<view class="api-desc">自定义短链链接解析为长链接</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/parseshorturl</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/parseshorturl?shortUrl=https://api.devopsclub.cn/2ndtW1b2Tj
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
shortUrl | String | 是 | 短链接地址

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 转换成功状态码(0 成功, 非零 失败)
domain | String | 短地址配置的域名
longUrl | String | 原长连接地址
msg | String | 消息

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "code": 0,
        "domain": "api.devopsclub.cn",
        "longUrl": "https://www.cnblogs.com/zhichaoma",
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