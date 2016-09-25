package scraper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/antonholmquist/jason"
)

// Data struct for response
type Data struct {
	Name     string
	URL      string
	ImageURL string
	Price    int64
	Location string
	Seller   string
	Source   string
}

// GetAllData return all array of Data struct from selected marketplace
func GetAllData(q string) (allData []Data) {
	chan1 := make(chan []Data, 1)
	chan2 := make(chan []Data, 1)
	chan3 := make(chan []Data, 1)
	chan4 := make(chan []Data, 1)

	go func() {
		chan1 <- GetTokopedia(q)
		chan2 <- GetBukalapak(q)
		chan3 <- GetLazada(q)
		chan4 <- GetBlibli(q)

		close(chan1)
		close(chan2)
		close(chan3)
		close(chan4)
	}()

	allData = append(allData, chanToSlice(chan1)...)
	allData = append(allData, chanToSlice(chan2)...)
	allData = append(allData, chanToSlice(chan3)...)
	allData = append(allData, chanToSlice(chan4)...)

	return allData
}

func chanToSlice(c chan []Data) (data []Data) {
	for v := range c {
		data = v
	}

	return data
}

func parser(url, catalogClass, sellerClass, priceClass, nameClass, sourceName, urlClass, imageURLClass, locationClass string) (data []Data) {
	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Print(err)
		return
	}

	doc.Find(catalogClass).Each(func(i int, s *goquery.Selection) {

		seller := strings.Replace(s.Find(sellerClass).Text(), "\n", "", -1)
		location := strings.Replace(s.Find(locationClass).Text(), "\n", "", -1)
		price := strings.Replace(s.Find(priceClass).Text(), "\n", "", -1)
		name := strings.Replace(s.Find(nameClass).Text(), "\n", "", -1)
		source := sourceName
		url, _ := s.Find(urlClass).Attr("href")
		imageURL, _ := s.Find(imageURLClass).Attr("src")
		data = append(data, Data{
			name,
			url,
			imageURL,
			formatPrice(price),
			location,
			seller,
			source,
		})

	})

	return
}

// GetTokopedia return item list from tokopedia.com
func GetTokopedia(q string) (data []Data) {
	url := "https://ace.tokopedia.com/search/v1/product?st=product&q=" + q + "&source=search&device=desktop&full_domain=www.tokopedia.com&scheme=https&page=1&fshop=1&rows=25&unique_id=6403ca6f11e44b3cbb5828ba30893d1c&start=0&ob=23&full_domain=www.tokopedia.com"

	doc, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}

	defer doc.Body.Close()

	body, _ := ioutil.ReadAll(doc.Body)

	if isJSON(string(body)) == true {
		v, _ := jason.NewObjectFromBytes(body)
		items, _ := v.GetObjectArray("data")
		source := "tokopedia"

		for _, item := range items {
			name, _ := item.GetString("name")
			url, _ := item.GetString("uri")
			imageURL, _ := item.GetString("image_uri")
			price, _ := item.GetString("price")
			shop, _ := item.GetObject("shop")
			location, _ := shop.GetString("location")
			seller, _ := shop.GetString("name")

			data = append(data, Data{
				name,
				url,
				imageURL,
				formatPrice(price),
				location,
				seller,
				source,
			})
		}
	}

	return
}

// GetBukalapak return item list from bukalapak.com
func GetBukalapak(q string) (data []Data) {
	data = parser(
		"https://www.bukalapak.com/products?utf8=%E2%9C%93&search%5Bkeywords%5D="+q,
		".product-card",
		".user__name",
		".product-price amount",
		".product__name",
		"bukalapak",
		".product__name",
		".product-media__img",
		".user-city")
	return
}

// GetBlibli return item list from blibli.com
func GetBlibli(q string) (data []Data) {
	data = parser(
		"https://www.blibli.com/search?s="+q,
		"#catalogProductListContentDiv .large-4",
		".user__name",
		".new-price-text",
		".product-title",
		"blibli",
		".single-product",
		".lazy",
		".user-city")
	return
}

// GetLazada return item list from lazada.co.id
func GetLazada(q string) (data []Data) {
	data = parser(
		"http://www.lazada.co.id/catalog/?q="+q,
		".product-card",
		".user__name",
		".product-card__price",
		".product-card__name",
		"lazada",
		"a",
		"img",
		".user-city")
	return
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func formatPrice(price string) int64 {
	reg := regexp.MustCompile(`[^0-9]`)
	strPrice := reg.ReplaceAllString(price, "")
	intPrice, err := strconv.ParseInt(strPrice, 10, 64)

	if err != nil {
		return 0
	}

	return intPrice
}
