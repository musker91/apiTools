> title名称

<view class="api-title">短链接生成</view>

> api描述

<view class="api-desc">长链接转化为短链, 支持 (t.cn/url.cn ...)</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/tooffshorturl</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/tooffshorturl?url=https://api.devopsclub.cn&type=0
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
url | String | 是 | 长链接
type | Int | 否 | 转换的短链接类型(0:url.cn/1:t.cn)

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询状态码(0 成功, 非零 失败)
msg | String | 消息
longurl | String | 长链接
shorturl | String | 短链接

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "longurl": "https://api.devopsclub.cn",
        "shorturl": "https://t.cn/A6zm07gm"
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