package scraper

import (
  "log"
  "strings"
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

func GetAllData(q string) (allData [][]Data) {
    allData = append(allData, GetTokopedia(q))
    allData = append(allData, GetBukalapak(q))
    allData = append(allData, GetBlibli(q))
    allData = append(allData, GetLazada(q))

    return allData
}

func Parser(url, catalogClass, sellerClass, priceClass, nameClass, sourceName, urlClass, imageUrlClass, locationClass string) (data []Data){
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
        url,_ := s.Find(urlClass).Attr("href")
        imageUrl,_ := s.Find(imageUrlClass).Attr("src")
        data = append(data, Data {
            name,
            url,
            imageUrl,
            price,
            location,
            seller,
            source,
        })

    })

    return
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

    return 
}

func GetBukalapak(q string) (data []Data) {
    data = Parser(
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

func GetBlibli(q string) (data []Data) {
    data = Parser(
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

func GetLazada(q string) (data []Data) {
    data = Parser(
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
