# Gin日志中间件和跨域中间件

## 安装
```shell script
go get -u github.com/maczh/gintool
```

### 使用
> 在Gin引擎初始化后，载入中间件
```go
    import "github.com/maczh/gintool"

	engine := gin.Default()

	//设置接口日志
	engine.Use(gintool.SetRequestLogger())

	//添加跨域处理
	engine.Use(gintool.Cors())
```

**说明：若添加了链路追踪中间件，则engine.Use(mgtrace.TraceId())必须放在第一个加载**

若主配置文件中设置了启用mongodb，且设置了 go.log.req参数(MongoDB中Collection日志表名称参数)，则会自动保存一份日志到MongoDB当中
如：
```yaml
go:
  config:
    used: mongodb,...
  log:
    level: debug
    req: MyappRequestLog
```