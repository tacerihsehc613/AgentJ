package main

import (
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/tacerihsehc613/job-scraper/scrapper"
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(fileName, fileName)
}
func main() {
	//scrapper.Scrape("NODEJS")
	e := echo.New()
	e.Static("/static", "static")
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
