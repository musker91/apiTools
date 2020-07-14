> title名称

<view class="api-title">中文分词</view>

> api描述

<view class="api-desc">将一整段文字拆解为多个词语</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/segcut</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/segcut?text=你若要喜爱你自己的价值，你就得给世界创造价值
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
Text | String | 是 | 要分词的文本

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询状态码(0 成功, 非零 失败)
msg | String | 消息
result | Array | 分词后的词组列表

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "result": [
            "你",
            "若要",
            "喜爱",
            "你",
            "自己",
            "的",
            "价值",
            "，",
            "你",
            "就",
            "得",
            "给",
            "世界",
            "创造",
            "价值"
        ]
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