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
	tic := time.NewTicker(time.Second * 10)
	go func () {
		for range tic.C {
			time.AfterFunc(10 * time.Second, func () {el.Bit()})
		}
	}()
}

func dohandle() {
	tic := time.NewTicker(time.Second * 5)
	go func() {
		for range tic.C {
			time.AfterFunc(5 * time.Second, func () {el.Handle()})
		}
	}()
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