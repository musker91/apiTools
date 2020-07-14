> title名称

<view class="api-title">银行卡信息查询</view>

> api描述

<view class="api-desc">获取银行卡信息，发卡行、卡类型</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/bankcard</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/bankcard?bakCard=6222600260001072444
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
bakCard | String | 是 | 银行卡号

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询状态码(0 成功, 非零 失败)
msg | String | 消息
bank_name | String | 银行名称
bank_name_en | String | 银行英文名称
card_type | String | 银行卡类型
card_number | String | 银行卡号


</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "card_type": "太平洋借记卡",
        "bank_name": "交通银行",
        "bank_name_en": "COMM",
        "card_number": "6222600260001072444"
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