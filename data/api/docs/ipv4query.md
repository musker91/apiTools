> title名称

<view class="api-title">ipv4信息查询</view>

> api描述

<view class="api-desc">获取指定ipv4地址/域名ip信息</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/ipv4query</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/ipv4query?ip=49.95.48.136
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
ip | String | 是 | IP地址/域名

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询成功状态码(0 成功, 非零 失败)
msg | String | 消息
ip | String | ip地址
cityId | Int | 城市Id号
region | String | 区域
country | String | 国家
province | String | 省
city | String | 市
isp | String | isp厂商

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "ip": "49.95.48.136",
        "cityId": 1015,
        "country": "中国",
        "region": "",
        "province": "江苏省",
        "city": "南京市",
        "isp": "电信"
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