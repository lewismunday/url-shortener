package main

import (
	"github.com/gin-gonic/gin"
	"url-shortener/routes"
)

func main() {
	r := gin.Default()

	routes.SetupRoutes(r)

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
