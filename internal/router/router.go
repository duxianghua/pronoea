package router

import (
	"github.com/duxianghua/pronoea/internal/handler"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// Unified exception handling
	r.Use(handler.Recover)
	r.Use(handler.ErrorHandler)
	// r.Static("/index", "./html/")
	//r.StaticFile("/favicon.ico", "./html/favicon.ico")
	r.Static("/static", "./html/static")
	//r.StaticFS("/static", http.Dir("./html"))
	//r.GET("/", handler.RedirectHomeHandler)
	// r.GET("/", func(ctx *gin.Context) {
	// 	ctx.Redirect(302, "/static/")
	// })
	r.NoRoute(func(ctx *gin.Context) {
		ctx.File("./html/index.html")
	})
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// health apis
	r.GET("/api/ping", (&handler.HealthHandler{}).Ping)

	// http://Host:Port/api/v1/...
	api := r.Group("api")
	{
		v1 := api.Group("v1")
		{
			route_v1(v1)
		}
		v2 := api.Group("v2")
		{
			route_v2(v2)
		}
	}
	api2 := r.Group("/dev-api/api")
	{
		v1 := api2.Group("v1")
		{
			route_v1(v1)
		}
		v2 := api2.Group("v2")
		{
			route_v2(v2)
		}
	}
	api3 := r.Group("/stage-api/api")
	{
		v1 := api3.Group("v1")
		{
			route_v1(v1)
		}
		v2 := api3.Group("v2")
		{
			route_v2(v2)
		}
	}
}

func route_v1(rg *gin.RouterGroup) {
	probeApi(rg)
	contactGroupApi(rg)
	alertsApi(rg)
	ScenariosAPI(rg)
}

func route_v2(rg *gin.RouterGroup) {

}
