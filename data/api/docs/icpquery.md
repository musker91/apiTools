> title名称

<view class="api-title">ICP备案信息查询</view>

> api描述

<view class="api-desc">根据域名查询ICP备案信息</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/icpquery</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/icpquery?url=devopsclub.cn
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
organizer_name | String | 主办单位名称
organizer_nature | String | 主办单位性质
recording_license_number | String | 网站备案/许可证号
site_name | String | 网站名称
site_index_url | String | 网站首页地址
review_time | String | 审核时间

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "organizer_name": "马智超",
        "organizer_nature": "个人",
        "recording_license_number": "京ICP备17029263号-3",
        "site_name": "技术站点",
        "site_index_url": "www.devopsclub.cn",
        "review_time": "2020/3/2 2:11:00"
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