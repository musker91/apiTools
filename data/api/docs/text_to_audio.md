> title名称

<view class="api-title">文字转语音</view>

> api描述

<view class="api-desc">将文字转换为语音文件</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/text_to_audio</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/text_to_audio?lang=zh&charset=UTF-8&speed=2&text=我是要转换的文字
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
lang | String | 是 | 语言(zh 中文, en 英文)
charset | String | 是 | 文字编码(UTF-8, GBK)
speed | Int | 是 | 语速，可以是1-9的数字，数字越大，语速越快。
text | String | 是 | 要转化的文字

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询状态码(0 成功, 非零 失败)
msg | String | 消息

</view>

> 返回示例

<view class="api-reponse-demo">

```text
// 转换成功直接返回文件
// 转换失败返回如下格式json
{
    "code": 1,
    "msg": "text conversion audio failed"
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