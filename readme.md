#   项目代码结构说明文档

#### main.go：程序入口，连接数据库，打开日志，配置路由。

#### svc文件
    - handler.go:定义处理函数。
    - newhandler.go:处理函数的主体。
    - validate.go:检查参数是否合法 或 用户是否有权限访问。
    - user.go:用户接口的业务逻辑。
    - manage.go:管理员接口的业务逻辑。
#### utils文件
    - tool.go:一些工具，如密码加盐、获得当前时间等。
    - log.go:封装的log工具，log分三个等级，info、error和fatal。
#### sql文件
    - create_table.sql:项目的数据库表设计。
#### model文件
    - request.go:请求参数。
    - response.go:返回参数。
    - do.go:与数据库表结构对应。
    - bo.go
#### dao文件
    - init.go:初始化mysql。
    - mysql.go:操作mysql的代码。
##### constant文件
    - constant.go:常量 
#### config
    - server.yml:配置文件
#### codes
    - codes.go:自定义错误码