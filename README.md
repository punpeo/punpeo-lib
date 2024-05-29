## 介绍
 pun的go公用仓库，旨在减少微服务中各项目重复代码的复制，便于各项目的升级。


## 环境要求
go 1.18.8

## 私有包结构规范
```
.
├── event   //公用事件：广播类事件，主要是常量和结构体
├── extends //拓展模块：如七牛云、腾讯云、微信支付sdk等 第三方服务模块
├── log     //日志相关：统一接入elk规范
├── mq      //消息队列
│   └── rmq //rabbit mq
├── zeropro //gozero框架个性化升级：如 配置相关或接入etcd/prometheus/jaeger/zrpc等的升级
├── rest    //http协议request和response相关
├── stores  //数据存储相关：遵循单例和连接池模型
│   ├── cache //缓存中间件:如 redis
│   ├── db  //数据库 如：mysql
│   └── es  //其他存储单独目录 如 es、 doris
├── task    //任务脚本相关：go-queue，go-cron接入等
└── utils   //常用函数/效率相关：如 时间、excel、图片、排序/数据结构处理、通用算法、token加解密、序列化、php对接相关等等

```


## 安装
```
go env -w GOPRIVATE="github.com/punpeo/punpeo-lib"
go get -u github.com/punpeo/punpeo-lib
```

## 本地调试
```
replace (
gitlab.github.com/punpeo/punpeo-lib v0.0.0-20221010024834-5a4bad007892 => ../go-lib
)
```

## 审核流程
详见：https://jz-tech.yuque.com/jz-tech/adtrtg/trt2l2ox8o3gfhnc

## 文档
```
go get -v golang.org/x/tools/cmd/godoc
godoc -http ":8090"
```

## 版本更新记录
<a href="./CHANGELOG.md">点击查看</a>
