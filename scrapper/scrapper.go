package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	title    string
	company  string
	location string
	desc     string
	link     string
}

func Scrape(term string) {
	var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=" + term
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL)
	fmt.Println(totalPages)

	for i := 1; i <= totalPages; i++ {
		//extractedJobs := getPage(i)
		//jobs = append(jobs, extractedJobs...)
		go getPage(i, baseURL, c)
	}
	for i := 1; i <= totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}
	//fmt.Println(jobs)
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageUrl := url + "&recruitPage=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".item_recruit")

	searchCards.Each(func(i int, card *goquery.Selection) {
		//job := extractJob(card)
		//jobs = append(jobs, job)
		go extractJob(card, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	//return jobs
	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	link, _ := card.Find(".item_recruit>.area_job>.job_tit>a").Attr("href")
	title := CleanString(card.Find(".item_recruit>.area_job>.job_tit>a").Text())
	company := CleanString(card.Find(".item_recruit>.area_corp > .corp_name > a").Text())
	location := CleanString(card.Find(".item_recruit>.area_job>.job_condition>span>a").Text())
	lang := CleanString(card.Find(".item_recruit>.area_job>.job_sector>b>a").Text())
	desc := ""

	link = "saramin.co.kr" + link

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
	c <- extractedJob{title: title,
		company:  company,
		location: location,
		desc:     desc,
		link:     link,
	}
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
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

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Title", "Company", "Location", "Desc", "Link"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.title, job.company, job.location, job.desc, job.link}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
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

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
