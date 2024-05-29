# go-zero工具包

> ***notice：***
>
> 这个包里面的代码会依赖go-zero框架 主要实现go-zero中一些没有实现到到但是开发中经常用的功能
>
> 例如：验证器,常用中间件等等等
>

## 目前已完成

### 1: 输出jsonp


### 2: 验证器

使用方法: 将go-zero中的requests.Parse替换为zeropro下的ParseAndValidate方法，其他完全兼容go-zero(自动集成中文输出)

```go

import "github.com/punpeo/punpeo-lib/zeropro/requests"

if err := requests.ParseAndValidate(r, &req); err != nil {
    httpx.Error(w, err)
    return
}
```

验证规则： [参考这里](https://github.com/go-playground/validator)

**ps: 推荐使用模板生成这这部分代码**

### 3: API 标准输出

统一API接口规则   

[参考这里](https://jz-tech.yuque.com/jz-tech/adtrtg/fuvuci)

**ps: 推荐使用模板生成这这部分代码**

### 4: 模板

支持的自定义模板

1： handler模板：集成验证器,API 标准输出,关于go-zero模板,[参考这里](https://go-zero.dev/cn/docs/goctl/template-cmd/)

使用说明(基于win环境)

```shell
# 初始化版本
PS E：\go-zero\template> template init
Templates are generated in C:\Users\xxx\.goctl\1.4.0, edit on your risk!      

# 替换对应模板文件
# 将 zeropro/goctl/handler.tpl 覆盖 C:\Users\xxx\.goctl\1.4.0\api\handler.tpl
```

生成模板，这里以官方文档这个API项目为例子演示，[官方文档](https://go-zero.dev/cn/docs/advance/api-coding)

```shell
# 编写.api文件

# 生成api项目文件 
PS E:\go-zero\template> goctl api go -api user.api -dir . 
Done.   

# 依赖安装
PS E:\go-zero\template> go mod tidy 


```
2: dockerfile模板： 构建docker镜像调整


### 5: 常用中间件(待完成)

### 6: 常用grpc interceptor(待完成)

### 7：监控部分
[参考这里](./readme-monitor.md)



