# [几个解放双手的 Go 开发利器](https://mp.weixin.qq.com/s/OQ17Y4Knffd8SX-n-fnMsQ)

Original 机器铃砍菜刀 [Golang技术分享](javascript:void(0);) *2021-09-02 08:01*

收录于合集#Go工具与库19个

Go 开发中，我们会构造各种 struct 对象，经常会有 json、数据库表、yaml、toml 等数据结构转 strcut 的需求。这时，我们可以根据字段名和数据类型来将这些数据结构，手动地填充至 Go 代码的 strcut 。但当数据字段很多时，这种方式不但耗时耗力，还容易出现一些低级错误。

针对以上情况，本文推荐几个开箱即用的开发利器，帮助 Gopher 解放双手，拯救时间。

## JSON-to-Go

JSON-to-Go 是一个将 json 数据转换为 Go 结构体的在线服务。

地址：*https://mholt.github.io/json-to-go/*

![Image](https://mmbiz.qpic.cn/mmbiz_png/2EiaKLQksVQI6zXTPxXDrvvvh8gPeKZm4lQI6RQlJZpImPtWfuP6tia6URC0Yxosiaibxl7ZDTIT46BqicROLUzibXmA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## TOML-to-Go

TOML-to-Go 是一个将 toml 数据转换为 Go 结构体的在线服务。

地址：*https://xuri.me/toml-to-go/*

![Image](https://mmbiz.qpic.cn/mmbiz_png/2EiaKLQksVQI6zXTPxXDrvvvh8gPeKZm4fs7cIWDGGckvEtMdhiaIa5xicdJX85IoibicteZ8MjykQdOYPv3AnJ4Ikg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## YAML-to-Go

TOML-to-Go 是一个将 yaml 数据转换为 Go 结构体的在线服务。

地址：*https://zhwt.github.io/yaml-to-go/*

![Image](https://mmbiz.qpic.cn/mmbiz_png/2EiaKLQksVQI6zXTPxXDrvvvh8gPeKZm4ADrZPZO60L9WkvNMnhxt1mCWtt8FP0WFEnMZ0Pficg2KdEt3kRBaYsw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## curl-to-Go

curl-to-Go 是一个将 curl 请求命令和数据格式转换为 Go 相关代码的在线服务。

地址：*https://mholt.github.io/curl-to-go/*

![Image](https://mmbiz.qpic.cn/mmbiz_png/2EiaKLQksVQI6zXTPxXDrvvvh8gPeKZm4fWguE0kUuFuoTNO9RfxSHuODJ6y7AlejfmX9MdfrpFOt82fcBma0JA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

## sql2struct

sql2struct 是一款根据 sql 语句自动生成 Go 结构体的 chrome 插件。

地址：*https://github.com/idoubi/sql2struct*

它的安装非常简单，根据上面地址中给出的使用步骤即可。

当我们需要对某个数据表，例如小菜刀本地数据库中的 `rent` 库，执行以下命令，拿到 SQL 定义语句。

```
mysql> show create table rent\G;
*************************** 1. row ***************************
       Table: rent
Create Table: CREATE TABLE `rent` (
  `name` varchar(100) DEFAULT NULL,
  `price` int(10) DEFAULT NULL,
  `area` varchar(60) DEFAULT NULL,
  `number` varchar(60) DEFAULT NULL,
  `structure` varchar(60) DEFAULT NULL,
  `pay` varchar(60) DEFAULT NULL,
  `orientaion` varchar(60) DEFAULT NULL,
  `floor` varchar(60) DEFAULT NULL,
  `region` varchar(100) DEFAULT NULL,
  `metro` varchar(60) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8
1 row in set (0.00 sec)

ERROR:
No query specified
```

打开 sql2struct 插件，将 SQL 建表语句置入，即可得到对应的 Go 代码 struct 信息。

![Image](https://mmbiz.qpic.cn/mmbiz_png/2EiaKLQksVQI6zXTPxXDrvvvh8gPeKZm4mYlftEYxyqKDkySK6h8ZmYSpiamiaKo7C9NCpRiaIphINwGI9vg2pcjdw/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

当然，我们还可以通过 options 选择多种字段标签，例如上例中，选择的是 gorm 和 json。

![Image](https://mmbiz.qpic.cn/mmbiz_png/2EiaKLQksVQI6zXTPxXDrvvvh8gPeKZm4jib3pz0gmYGCSP8mTqsT27wian90OViamicpvgJwDcUYtyrAnCpk2kVTEg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)