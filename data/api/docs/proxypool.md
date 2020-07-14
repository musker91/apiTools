> title名称

<view class="api-title">全球免费代理IP库</view>

> api描述

<view class="api-desc">全球免费代理IP库，高可用IP，精心筛选优质IP（注:代理IP采集于网络，仅供个人学习使用。请勿用于非法途径，违者后果自负！）</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/proxypool</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/proxypool
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
page | Int | 否 | 当前页码(每页15条数据)
protocol | String | 否 | 代理的协议类型(http/https)
anonymity | String | 否 | 匿名类型(透明/高匿)
country | String | 否 | 所在国家
address | String | 否 | 所在地区
isp | String | 否 | 运营商
order_by | String | 否 | 排序字段 (speed:响应速度,verify_time:校验时间)
order_rule |String | 否 | 排序规则(desc:降序 asc:升序)

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 获取成功状态码(0 成功, 非零 失败)
msg | String | 消息
ip | String | ip地址
port | String | 端口号
anonymity | String | 匿名类型
protocol | String | 协议类型
country | String | 所在国家
address | String | 所在地区
isp | String | isp厂商
speed | Int | 响应速度(毫秒, 为0表示最后一次校验无响应)
verify_time | DateTime | 最后校验时间
pages | Int | 总页数

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "data": [
            {
                "ip": "116.196.85.150",
                "port": "3128",
                "anonymity": "高匿",
                "protocol": "https",
                "country": "中国",
                "address": "中国 北京 北京",
                "isp": "电信",
                "speed": 369,
                "verify_time": "2020/02/28 00:35:38"
            },
            {
                "ip": "117.27.152.236",
                "port": "1080",
                "anonymity": "高匿",
                "protocol": "http",
                "country": "中国",
                "address": "中国 福建 福州",
                "isp": "电信",
                "speed": 217,
                "verify_time": "2020/02/28 00:35:38"
            }
        ],
        "pages": 7
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