package main

import (
    "os"
    "github.com/kataras/iris"
    "github.com/disiko/disiko-backend/api"
)

func main() {
    // configuration
    baseApiUrl := "/api/"
    port := ":"+os.Getenv("PORT")

    net := iris.New()
    net.Post(baseApiUrl+"scraper", api.GetScraper)
    net.Listen(port)
}

