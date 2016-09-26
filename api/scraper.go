package api

import (
	"disiko-backend/lib/scraper"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/kataras/iris"
)

// FilteredRequest is a struct for GetFiltered body
type FilteredRequest struct {
	Query     string `form:"q"`
	Order     string `form:"order"`
	Direction string `form:"dir"`
}

// GetFiltered return filtered items with specific user params
// Exmaple {baseAPIURL}search?q=psp&order=name&dir=asc
func GetFiltered(ctx *iris.Context) {
	filteredReq := FilteredRequest{}
	err := ctx.ReadForm(&filteredReq)

	if err != nil {
		resBadRequest(ctx)
		return
	}

	req := structs.Map(filteredReq)

	query := strings.Replace(req["Query"].(string), " ", "+", -1)
	order := "Name"
	direction := "asc"

	if req["Order"] != "" {
		order = req["Order"].(string)
	}

	if req["Direction"] == "desc" {
		direction = "desc"
	}

	fmt.Println(query, order, direction)

	allData := scraper.GetAllData(query)

	if len(allData) == 0 {
		resNoContent(ctx)
		return
	}

	ctx.Render("application/json", allData)
}

// GetFeatured return featured items from selected marketplace
func GetFeatured(ctx *iris.Context) {
	ctx.Render("application/json", iris.Map{"foo": "bar"})
}

func resBadRequest(ctx *iris.Context) {
	ctx.Render("application/json", iris.Map{"status": 400, "error": "Bad request"})
}

func resNoContent(ctx *iris.Context) {
	ctx.Render("application/json", iris.Map{"status": 204, "message": "No content"})
}
