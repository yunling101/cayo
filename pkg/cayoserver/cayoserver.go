package cayoserver

import (
	"crypto/tls"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/auth"
	"github.com/yunling101/cayo/pkg/cayoserver/client"
	"github.com/yunling101/cayo/pkg/cayoserver/controller/alarm"
	"github.com/yunling101/cayo/pkg/cayoserver/controller/console"
	"github.com/yunling101/cayo/pkg/cayoserver/controller/node"
	"github.com/yunling101/cayo/pkg/cayoserver/controller/settings"
	"github.com/yunling101/cayo/pkg/cayoserver/controller/task"
	"github.com/yunling101/cayo/pkg/cayoserver/controller/web"
	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/propb"
	"github.com/yunling101/cayo/public"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// NewRpcStart
func NewRpcStart(cert, key string) {
	if global.Config().RpcServer.Enable {
		certs, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			log.Fatalf("failed to load key pair: %s", err)
		}

		listen, err := net.Listen("tcp", global.Config().RpcServer.Listen)
		if err != nil {
			log.Fatalf("listen to net: %v", err)
		}
		opts := []grpc.ServerOption{
			grpc.Creds(credentials.NewServerTLSFromCert(&certs)),
		}
		server := grpc.NewServer(opts...)
		propb.RegisterMonitorServer(server, &client.RpcController{})

		reflection.Register(server)
		log.Printf("rpc listen to %s", global.Config().RpcServer.Listen)
		if err := server.Serve(listen); err != nil {
			log.Fatalf("listen to serve: %v", err)
		}
	}
}

// NewHTTPStart
func NewHTTPStart() {
	if global.Config().WebServer.Enable {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		router.Use(CorsHandler())
		router.Use(auth.LoginSessions())

		// 静态资源
		// https://github.com/gin-contrib/static
		router.StaticFS(public.StaticFS("/static"))
		router.StaticFS(public.StaticFS("/images"))
		router.SetHTMLTemplate(template.Must(template.New("").ParseFS(public.FS, "build/index.html")))

		// 前台首页
		router.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{})
		})
		// 登录注销
		controller := web.WebController{}
		router.Use(auth.RequestBody())
		{
			router.POST("/login", controller.Login)
			router.GET("/logout", controller.LogOut)
		}
		// API接口
		apiV1 := router.Group("/v1")
		apiV1.Use(auth.LoginRequired())
		apiV1.Use(auth.RequestBody())
		{
			apiV1.GET("/get_user_owner", console.ConsoleController{}.Owner)
			apiV1.GET("/m/summary", console.ConsoleController{}.Summary)
			apiV1.GET("/m/week/stat", console.ConsoleController{}.WeekStat)
			apiV1.GET("/m/week/ratio", console.ConsoleController{}.WeekRatio)
			apiV1.GET("/m/top/provider", console.ConsoleController{}.TopProvider)
			apiV1.GET("/m/top/node", console.ConsoleController{}.TopNode)
			apiV1.GET("/get_task_type", task.TaskController{}.TaskType)
			apiV1.GET("/get_dns_type", task.TaskController{}.DnsType)
			apiV1.GET("/task/list", task.TaskController{}.List)
			apiV1.POST("/task/add", task.TaskController{}.Add)
			apiV1.DELETE("/task/delete/:id", task.TaskController{}.Delete)
			apiV1.POST("/task/state", task.TaskController{}.State)
			apiV1.GET("/task/details/summary/:id", task.MonitorController{}.SummaryArea)
			apiV1.GET("/task/details/metric", task.MonitorController{}.SummaryMetric)
			apiV1.GET("/task/details/rules", task.MonitorController{}.SummaryRules)
			apiV1.GET("/alarm/task", alarm.AlarmController{}.AlarmTask)
			apiV1.GET("/alarm/metric", alarm.AlarmController{}.MetricList)
			apiV1.GET("/alarm/rule", alarm.AlarmController{}.RuleList)
			apiV1.POST("/alarm/rule/add", alarm.AlarmController{}.RuleAdd)
			apiV1.POST("/alarm/rule/state", alarm.AlarmController{}.RuleState)
			apiV1.DELETE("/alarm/rule/delete/:id", alarm.AlarmController{}.RuleDelete)
			apiV1.GET("/get_notify_channel", alarm.AlarmController{}.NotifyChannel)
			apiV1.GET("/alarm/contact", alarm.AlarmController{}.ContactList)
			apiV1.POST("/alarm/contact/add", alarm.AlarmController{}.ContactAdd)
			apiV1.POST("/alarm/contact/stat", alarm.AlarmController{}.ContactStat)
			apiV1.POST("/alarm/contact/reset", alarm.AlarmController{}.ContactReset)
			apiV1.DELETE("/alarm/contact/delete/:id", alarm.AlarmController{}.ContactDelete)
			apiV1.GET("/notify/contact", alarm.AlarmController{}.NotifyContact)
			apiV1.GET("/node/probe", node.NodeController{}.Probe)
			apiV1.GET("/get_operator_list", node.NodeController{}.OperatorList)
			apiV1.GET("/node/list", node.NodeController{}.List)
			apiV1.POST("/node/add", node.NodeController{}.Add)
			apiV1.POST("/node_selected_status", node.NodeController{}.SelectedStatus)
			apiV1.DELETE("/node/delete/:id", node.NodeController{}.Delete)
			apiV1.POST("/modify_user_settings", settings.SettingsController{}.ModifyUserSettings)
			apiV1.POST("/modify_user_password", settings.SettingsController{}.ModifyUserPassword)
			apiV1.GET("/get_user_settings", settings.SettingsController{}.GetUserSettings)
		}

		s := &http.Server{
			Addr:           global.Config().WebServer.Listen,
			Handler:        router,
			MaxHeaderBytes: 1 << 32,
		}
		log.Printf("http listen to %s", global.Config().WebServer.Listen)
		s.ListenAndServe()
	}
}

func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许访问所有域
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:3000") // 要么一个值，要么是*
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, token")
		c.Header("Access-Control-Max-Age", "172800")

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Set("content-type", "application/json")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
