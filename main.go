package main

import (
	"fmt"
	"net/http"
	"tasktask/src/el"
	"tasktask/src/middleware"
	"time"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/pprof"
	dbt "tasktask/src/sqlitem"
)

// cd /www/wwwroot/qrv-mid.appbsl.cn/extend/gomid/
// cd /www/wwwroot/qrv-mid.appbsl.cn/extend/gomm/
// go tool pprof http://localhost:10489/debug/pprof/heap

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/add", add)
	r.GET("/h", handle)
	go dohandle()
	go bit()
	fmt.Println("running at 10489");
	pprof.Register(r)
	r.Run(":10489")
}
func bit() {
	for {
		el.Bit()
		time.Sleep(time.Second * 10)
		// time.Sleep(time.Microsecond * 1)
		runtime.GC()
	}
}

func dohandle() {
	c := new(dbt.Con)
	c.Opendb()
	for {
		el.Handle(c)
		time.Sleep(time.Second * 5)
		// time.Sleep(time.Microsecond * 1)
		runtime.GC()
	}
}

func handle(c *gin.Context) {
	con := new(dbt.Con)
	con.Opendb()
	el.Handle(con)
	c.JSON(http.StatusOK, "fin")
}

func add(c *gin.Context) {
	res := el.New(c)
	c.JSON(http.StatusOK, res)
}
