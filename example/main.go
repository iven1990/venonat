package main

import (
	"github.com/iven1990/enheng/venonat"
	"net/http"
	"fmt"
)

func main()  {
	e := venonat.New()
	if e == nil {
		fmt.Println("engine is nil")
	} else {
		fmt.Println("engine is not nil")

	}
	e.GET("/hello", func(c *venonat.Context) {
		fmt.Println("hello ok")
		c.JSON(http.StatusOK, venonat.H{
			"hello": "hello",
		})
	})
	e.Run(":9889")
}
