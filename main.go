package main

import (
	"github.com/kataras/iris/v12"

	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	cnf "github.com/klaudijuskungys/zabbix-exporter-3000/config"
	hdl "github.com/klaudijuskungys/zabbix-exporter-3000/handlers"
)

func main() {
	app := iris.New()

	app.Logger().SetLevel("INFO")

	app.Use(recover.New())
	app.Use(logger.New())

	// prometheus metrics
	m := prometheusMiddleware.New("ze3000", 0.3, 1.2, 5.0)
	hdl.RecordMetrics()
	app.Use(m.ServeHTTP)
	app.Get(cnf.MetricUriPath, iris.FromStd(promhttp.Handler()))

	app.Get("/liveness", func(ctx iris.Context) {
		ctx.WriteString("ok")
	})

	app.Get("/readiness", func(ctx iris.Context) {
		ctx.WriteString("ok")
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.WriteString("zabbix-exporter-3000")
	})

	app.Run(iris.Addr(cnf.MainHostPort), iris.WithoutServerError(iris.ErrServerClosed))
}
