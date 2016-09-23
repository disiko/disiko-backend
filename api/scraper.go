package api

import (
    "strings"
    "github.com/kataras/iris"
    "github.com/disiko/disiko-backend/lib/scraper"
)

func GetScraper(ctx *iris.Context){
    q:= ctx.PostValue("q")
    q = strings.Replace(q, " ", "+", -1)
    ctx.Render("application/json", scraper.GetTokopedia(q))
}
