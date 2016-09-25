package api

import (
    "strings"
    "github.com/kataras/iris"
    "disiko-backend/lib/scraper"
)

func GetScraper(ctx *iris.Context) {
    q := ctx.PostValue("q")
    q = strings.Replace(q, " ", "+", -1)

    allData := scraper.GetAllData(q)

    ctx.Render("application/json", allData)
}
