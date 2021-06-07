# **A golang template which base gin, designed by yy**

## Installing and Getting started

      ├── app
      │   ├── config                   # 配置层
      │   ├── controller               # controller 层
      │   ├── data_struct              # 结构体存储层
      │   ├── enum                     # 枚举层
      │   ├── ginServer                # router初始化层
      │   ├── init                     # 初始化所有依赖层
      │   ├── libs                     # 小工具层
      │   ├── middlewares              # 中间件
      │   ├── miniprogram              # 小程序层
      │   ├── models                   # mysql 层
      │   ├── redis                    # redis初始化层
      │   ├── router                   # 路由层
      │   └── services                 # 逻辑层
      ├── docker-compose.yml           # docker-compose 启动文件
      ├── go.mod
      ├── go.sum
      ├── kibana.yml                   # kibana 日志分析工具配置文件
      ├── logs                         # 日志
      ├── main.go                      # 入口文件
      ├── README.md
      ├── statics                      # 静态资源文件夹
      │   └── excel
      ├── test                         # 测试文件夹
      └── volumes                      # docker的依赖文件夹(用于持久化docker数据)
          ├── elasticsearch            # es 文件夹
          ├── elasticsearch2           # es 文件夹
          ├── elasticsearch3           # es 文件夹
          └── kafka                    # kafka 文件夹

PS: Please create the directory that called volumes in root path.

      └── volumes                      # docker的依赖文件夹(用于持久化docker数据)
          ├── elasticsearch            # es 文件夹
          ├── elasticsearch2           # es 文件夹
          ├── elasticsearch3           # es 文件夹
          └── kafka                    # kafka 文件夹


1. Clone the repository.

       git clone git@github.com:guaidashu/law_article_find.git

2. Add some code in gin.go and context.go

   (1) Add code on line 105 in gin.go

        import(
            ...
            "reflect"
        )
        
        type Engine struct {
            ...
            AutoRouterGroup  map[string]map[string]string
            AutoRouterController map[string]reflect.Type
        }
        
        func (engine *Engine) AddToAutoRouterGroup(controllerName, method, methodName string) {
            if engine.AutoRouterGroup == nil {
                engine.AutoRouterGroup = make(map[string]map[string]string)
            }
            if engine.AutoRouterGroup[controllerName] == nil{
                engine.AutoRouterGroup[controllerName] = make(map[string]string)
            }
            engine.AutoRouterGroup[controllerName][method] = methodName
        }
        
        func (engine *Engine) AddToAutoRouterController(controllerName string, controller *reflect.Type) {
            if engine.AutoRouterController == nil {
                engine.AutoRouterController = make(map[string]reflect.Type)
            }
            engine.AutoRouterController[controllerName] = *controller
        }

   (2) Add code on anywhere in context.go

        import(
            ...
            "reflect"
        )
        
        func (c *Context) GetAutoRouterGroup() map[string]map[string]string {
            return c.engine.AutoRouterGroup
        }
        
        func (c *Context) GetAutoRouterController() *map[string]reflect.Type {
            return &c.engine.AutoRouterController
        }

## Usage

None

## FAQ

Contact to me with email "1023767856@qq.com" or "song42960@gmail.com"

## Running Tests

Add files to /test and run it.

## Finally Thanks 