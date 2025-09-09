package routes

import "github.com/gin-gonic/gin"

// 每个模块提供注册函数类型
type ModuleRegister func(rg *gin.RouterGroup)

// 模块注册表
var moduleRegistry = make(map[string]ModuleRegister)

// 注册模块
func RegisterModule(name string, registerFn ModuleRegister) {
	moduleRegistry[name] = registerFn
}

// 聚合所有模块路由
func RegisterAllModules(r *gin.Engine) {
	apiGroup := r.Group("/api")
	for _, registerFn := range moduleRegistry {
		registerFn(apiGroup)
	}
}
