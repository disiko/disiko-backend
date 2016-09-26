package main

import (
	"github.com/disiko/disiko-backend/api"
	"os"

	"github.com/kataras/iris"
)

func main() {
	// configuration
	baseAPIURL := "/api/"
	port := ":" + os.Getenv("PORT")

	if port == ":" {
		port = ":9999"
	}

	net := iris.New()

	net.Get(baseAPIURL+"featured", api.GetFeatured)
	net.Get(baseAPIURL+"search", api.GetFiltered)

	net.Listen(port)
}
