# **A golang template which base gin, designed by yy**

## Installing and Getting started
    
    .
    ├── README.md
    ├── app             # 应用目录
    │   ├── config      # 配置
    │   ├── controller  # 控制器
    │   ├── data_struct # 所有需要的额外结构体
    │   ├── ginServer   # gin engine代理
    │   ├── init        # 初始化所有内容，包括数据库等
    │   ├── libs        # 自定义工具包
    │   ├── models      # 数据库orm
    │   └── router      # 路由配置
    ├── go.mod          # go mod依赖文件
    ├── go.sum
    └── main.go         # 入口文件 

1. Clone the repository.

       git clone git@github.com:guaidashu/gin_template.git

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