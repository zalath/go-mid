package main

import (
	"fmt"
	"net/http"
	"tasktask/src/el"
	"tasktask/src/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

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
	r.GET("/count", count)
	go dohandle()
	go bit()
	fmt.Println("running at 10489");
	r.Run(":10489")
}
func bit() {
		ct := 0
		for {
			ct ++
			el.Bit()
			// fmt.Println("bit send")
			if ct > 10 {
				bit()
			}
			// time.Sleep(time.Microsecond * 1)
			time.Sleep(time.Second * 10)
		}
}

func dohandle() {
		ct := 0
		for {
			ct ++
			el.Handle()
			if ct > 10 {
				dohandle()
			}
			time.Sleep(time.Second * 5)
			// time.Sleep(time.Microsecond * 1)
		}
}

func handle(c *gin.Context) {
	el.Handle()
	c.JSON(http.StatusOK, "fin")
}

func add(c *gin.Context) {
	res := el.New(c)
	c.JSON(http.StatusOK, res)
}

func count(c *gin.Context) {
	res := el.Count()
	c.JSON(http.StatusOK, res)
}