> title名称

<view class="api-title">域名报红检测</view>

> api描述

<view class="api-desc">检测域名是否在微信、QQ中报红</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/dcheck</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/dcheck?url=baidu.com
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
url | String | 是 | 域名

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询状态码(0 成功, 非零 失败)
msg | String | 消息
domain | String | 域名
wx | String | 域名在微信中状态(danger 危险, unknown 未知, safe 安全)
qq | String | 域名在QQ中状态(danger 危险, unknown 未知, safe 安全)

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "domain": "baidu.com",
        "wx": "safe",
        "qq": "safe"
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