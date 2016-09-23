package scraper

import (
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "github.com/PuerkitoBio/goquery"
  "github.com/antonholmquist/jason"
)

type Data struct {
    name string
    url string
    imageUrl string
    price string
    location string
    seller string
    source string
}

func parser() (data []map[string]string){
    doc, err := goquery.NewDocument("http://metalsucks.net")

    if err != nil {
        log.Print(err)
        return
    }

    doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {

        band := s.Find("a").Text()
        title := s.Find("i").Text()
        data = append(data, map[string]string{"band": band, "title": title})
    })

    return
}

func GetAllData(q string) (allData [][]Data) {
    allData = append(allData, GetTokopedia(q))

    return allData
}

func GetTokopedia(q string) (data []Data) {
    url:= "https://ace.tokopedia.com/search/v1/product?st=product&q=" + q + "&source=search&device=desktop&full_domain=www.tokopedia.com&scheme=https&page=1&fshop=1&rows=25&unique_id=6403ca6f11e44b3cbb5828ba30893d1c&start=0&ob=23&full_domain=www.tokopedia.com"

    doc, err := http.Get(url)
    if err != nil {
        log.Print(err)
        return
    }

    defer doc.Body.Close()

    body, err := ioutil.ReadAll(doc.Body)

    if isJSON(string(body)) == true {
        v, _ := jason.NewObjectFromBytes(body)
        items, _ := v.GetObjectArray("data")
        source := "tokopedia"

        for _, item := range items {
            name, _ := item.GetString("name")
            url, _ := item.GetString("uri")
            imageUrl, _ := item.GetString("image_uri")
            price, _ := item.GetString("price")
            shop, _ := item.GetObject("shop")
            location,_ := shop.GetString("location")
            seller,_ := shop.GetString("name")

            data = append(data, Data {
                name,
                url,
                imageUrl,
                price,
                location,
                seller,
                source,
            })
        }
    }

    return data
}

func GetBukalapak(q string) (data []Data) {
    return data
}

func GetBlibli(q string) (data []map[string]string) {
    return data
}

func GetLazada(q string) (data []map[string]string) {
    return data
}

func isJSON(s string) bool {
    var js map[string]interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}

func parseBukalapak() (data []Data) {
    return
}

func parseBlibli() (data []Data) {
    return
}

func parseLazada() (data []Data) {
    return
}
