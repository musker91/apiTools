> title名称

<view class="api-title">QQ信息查询</view>

> api描述

<view class="api-desc">获取指定QQ号的头像和昵称</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/qqinfo</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/qqinfo?qq=1152490990
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
qq | String | 是 |qq号码

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询状态码(0 成功, 非零 失败)
msg | String | 消息
head_portrait | String | qq头像
nickname | String | qq昵称

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "head_portrait": "http://qlogo3.store.qq.com/qzone/1152490990/1152490990/100",
        "nickname": "Musker"
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