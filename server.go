package main

import (
    "github.com/kataras/iris"
    "github.com/disiko/disiko-backend/api"
)

func main() {
    // configuration
    baseApiUrl := "/api/"
    port := ":8080"

    net := iris.New()
    net.Post(baseApiUrl+"scraper", api.GetScraper)
    net.Listen(port)
}

