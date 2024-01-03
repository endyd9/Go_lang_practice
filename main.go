package main

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
	link string
	title string
	condition string
	limit string
}

var baseUrl string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python&recruitPageCount=100"

func main() {
	var jobs []extractedJob

	c := make(chan []extractedJob)
	writeC := make(chan bool)

	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	go writeJobs(jobs, writeC)
	<- writeC
	fmt.Println("Done, extracted '", len(jobs), "' jobs")
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection){
		pages = s.Find("a").Length()
	})

	return pages
}

func getPage(page int, mainC chan<- []extractedJob){
	var jobs []extractedJob
	c := make(chan extractedJob)

	pageURL := baseUrl + "&recruitPage=" + strconv.Itoa(page+1)

	res, err := http.Get(pageURL)
	fmt.Println("Requesting", pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCard := doc.Find(".item_recruit")

	searchCard.Each(func(i int, card *goquery.Selection){
		go extractJob(card, c)
	})
	
	for i := 0; i < searchCard.Length(); i++{
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("value")
	title := cleanString(card.Find(".job_tit>a").Text())
	condition := cleanString(card.Find(".job_condition").Text())
	limit := cleanString(card.Find(".date").Text())
	c <- extractedJob{
		link:"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx="+id,
		title: title,
		condition: condition,
		limit: limit,
	}
}

func checkErr(err error){
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request Failed whit: ", res.StatusCode)
	}
}

func cleanString(str string)string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func writeJobs(jobs []extractedJob, c chan<- bool){
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"LINK", "TITLE", "CONDITION", "LIMIT"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.link, job.title, job.condition, job.limit}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
	c <- true
}
