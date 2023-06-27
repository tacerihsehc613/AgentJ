package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	title    string
	company  string
	location string
	desc     string
}

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=NODEJS"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	fmt.Println(totalPages)

	for i := 1; i <= totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}
	fmt.Println(jobs)
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageUrl := baseURL + "&recruitPage=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".item_recruit")

	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	title := cleanString(card.Find(".item_recruit>.area_job>.job_tit>a").Text())
	company := cleanString(card.Find(".item_recruit>.area_corp > .corp_name > a").Text())
	location := cleanString(card.Find(".item_recruit>.area_job>.job_condition>span>a").Text())
	lang := cleanString(card.Find(".item_recruit>.area_job>.job_sector>b>a").Text())
	desc := ""

	var aTags []string
	card.Find(".item_recruit > .area_job > .job_sector > a").Each(func(i int, s *goquery.Selection) {
		aTags = append(aTags, s.Text())
	})

	// Concatenate the text values of the <a> tags with a comma separator
	desc += strings.Join(aTags, ", ")
	if lang != "" {
		desc = lang + ", " + desc
	}
	//fmt.Println(title, company, location, desc)
	return extractedJob{title: title,
		company:  company,
		location: location,
		desc:     desc}
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	//fmt.Println(doc)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
