> title名称

<view class="api-title">短视频解析</view>

> api描述

<view class="api-desc">短视频去水印/获取背景音乐/封面图片/视频标题(目前支持 抖音/皮皮虾/微视/最右/火山...)，更多短视频支持正在开发中</view>

> Api接口地址

<view class="api-url">https://api.devopsclub.cn/api/svp</view>

> 返回格式

<view class="api-reponse-format">JSON</view>

> 请求方式

<view class="api-request-method">GET/POST</view>

> 请求示例

<view class="api-request-demo">

```text
https://api.devopsclub.cn/api/svp?url=https://v.douyin.com/Efcp8a/
```

</view>

> 请求参数说明

<view class="request-param">

字段名称 | 类型 | 必填 | 说明
--- | --- | --- | ---
url | String | 是 | 短视频分享的url地址

</view>

> 返回参数说明

<view class="reponse-param">

字段名称 | 类型 | 说明
--- | --- | ---
code | Int | 查询状态码(0 成功, 非零 失败)
msg | String | 消息
desc | String | 短视频描述信息, 部分短视频平台无法获取视频描述信息
pic | String | 短视频封面图片链接
video | String | 短视频无水印视频地址
music | String | 短视频背景音乐, 大部分短视频平台无法获取单独音乐链接

</view>

> 返回示例

<view class="api-reponse-demo">

```json
{
    "code": 0,
    "data": {
        "desc": "#洪真英 #微胖才是极品",
        "pic": "https://p3.pstatp.com/large/tos-cn-p-0015/f2aee4622b8a46d094468cece4de2628_1589078966.jpg",
        "video": "http://v5-dy.ixigua.com/cdd4e976dba08b62698cb8e7d6bd4110/5eb832a8/video/tos/cn/tos-cn-ve-15/dcf0400a94924cf6b622b9de0f590a48/?a=1128&br=0&bt=1832&cr=0&cs=0&dr=0&ds=6&er=&l=202005102358080100140400382B4689C5&lr=&qs=0&rc=am9kZnQ6OWVwdDMzOmkzM0ApZGk4NTRkNWU5NztlaGc2OGcucmdlXjZka29fLS1fLS9zcy0vYV5iX19hMWBfYC9iLmA6Yw%3D%3D&vl=&vr=",
        "music": "http://p9-dy.byteimg.com/obj/ies-music/d2f8b67bd822ad8d3ee03d4999702176.mp3"
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