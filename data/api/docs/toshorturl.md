> title名称

<view class="api-title">自定义短链接转换</view>

> api描述

<view class="api-desc">长链接转换为自定义短链接</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/toshorturl</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/toshorturl?url=https://www.cnblogs.com/zhichaoma
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
url | String | 是 | 长链接URL
domain | String | 否 | 短链接域名绑定自己的域名 (默认为系统当前域名)
expireTime | Int | 否 | 设置过期时间, (以分钟为单位, -1代表用不过期)

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 转换成功状态码(0 成功, 非零 失败)
domain | String | 短地址配置的域名
shortUrl | String | 短链接地址
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
        "shortUrl": "http://api.devopsclub.cn/2ndtW1b2Tj"
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