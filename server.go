package main

import (
    "os"
    "github.com/kataras/iris"
    "disiko-backend/api"
)

func main() {
    // configuration
    baseApiUrl := "/api/"
    port := ":"+os.Getenv("PORT")

    if (port == ":") {
        port = ":9999"
    }

    net := iris.New()
    net.Post(baseApiUrl+"scraper", api.GetScraper)
    net.Listen(port)
}
