Release Notes

## 1.0.7
修复：
* 1、修复rabbitmq消费库参数poolsize设置为1时队列不消费的问题

## 1.0.6（2023年10月17日）
新增功能特性：
* 1、rabbitmq消费库（适用于go-zero框架）
  - 1.1、新增ants协程池支持&协程池信息统计
  - 1.2、支持手动ack模式
  - 1.3、新增channel&queue异常自动重连机制&异常消息通知
  - 1.4、新增消费panic捕获&日志打印
* 2、rabbitmq发送库（适用于go-zero框架）
  - 1.1、支持延迟队列投递

## 1.0.5 (2023年8月24日)
修复：
* 1、修复rest/restyclient请求完之后没有关闭资源，导致在高并发场景，链接释放慢导致泄露问题

新增：
* 1、添加safe goroutine, 并且能够recover panic, 会打印对应堆栈信息
* 2、添加group包，提供一组任务的goroutine同步，错误，取消功能

##  1.0.4 (2023年8月16日)

修复：
* 1、修复使用封装的excel第三方库，操作文件写入后，没有关闭文件释放资源，导致协程泄露。
* 2、stores/db/model/models.go 优化"分页"函数，兼容page和page_size小于等于0的情况。

##  1.0.3 (2023年4月20日)

修复：
* 1、修复因数据库空闲连接时间超时导致连接失效后，请求数据库时抛出"driver: bad connection"警告。

##  1.0.2 (2023年4月04日)

新增功能特性：
* 1、rest
    - 1.1、xerr 封装自定义xerr通用业务错误码
    - 1.2、restyclient 基于go-resty，封装通用http GET、POST请求
    - 1.3、result 封装自定义通用请求与返回结构体
    - 1.4、interceptor 封装自定义服务（rpc）拦截器
* 2、stores/db 数据库
    - 2.1、model 封装通用模型，以及基于gorm通用的Scopes方法
    - 2.2、xgorm 接入gorm使用
* 3、utils
    - 3.1、utils 新增通用函数
    - 3.2、timeutil 封装时间格式转换
* 4、extends
    - 4.1、qwrobot 接入企微机器人消息推送
    - 4.2、cos 腾讯云存储功能（新增功能：上传文件流/字节流、获取临时密钥）

##  1.0.1 (2023年3月9日)

新增功能特性：
* 1、extends/qcloud 接入腾讯云存储功能（包括：上传对象、下载对象、获取对象访问URL、获取对象预签名URL、获取对象预签名URL【自定义域名】）
* 2、utils/excel 基于excelize，封装excel操作（流式导出）


##  1.0.0 (2022年12月9日)

初始版本 功能特性：
* 1、extends/qiniu    七牛token获取
* 2、mq/rmq           rabbitmq链接入队
* 3、task/gocronnode  gocron 伪装节点服务
* 4、zeropro          依赖go-zero框架 实现一些框架功能的升级(jsonp,验证器,模板,错误日志入库)
* 5、utils          
    - 6.1、php序列化
    - 6.2、app登录态security_key加解密

///////
    


