package api

import (
    "strings"
    "github.com/kataras/iris"
    "disiko-backend/lib/scraper"
    "fmt"
)

func GetScraper(ctx *iris.Context) {
    q := ctx.PostValue("q")
    q = strings.Replace(q, " ", "+", -1)

    allData := scraper.GetAllData(q)

    fmt.Println(allData)
    // TODO Filter all data

    ctx.Render("application/json", iris.Map{"foo": "bar"})
}
