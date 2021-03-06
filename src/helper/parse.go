package helper

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
)

const TARGET_URL = "http://www.dmm.co.jp"
type Result struct {
	rank        string
	url         string
	name        string
	img         string
	information string
}

func(r *Result) setUrl(url string){
	r.url = TARGET_URL + url
}

func Extract(res http.Response) []Result {
	doc, err := goquery.NewDocumentFromResponse(&res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	var results []Result
	doc.Find("#w td.bd-b").Each(func(i int, s *goquery.Selection) {
		val := s.Find("p > a")

		var result Result
		result.rank = s.Find("span.rank").Text()
		url, _ := val.Attr("href")
		result.setUrl(url)
		result.name = val.Text()
		img, _ := s.Find("img").Attr("src")
		result.img = img
		result.information = s.Find(".data").Text()

		results = append(results, result)
	})
	return results
}
