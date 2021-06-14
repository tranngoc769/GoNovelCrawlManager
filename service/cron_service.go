package service

import (
	"fmt"
	"gonovelcrawlmanager/repository"
	"strings"

	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
)

type CronService struct {
	repo repository.NovelRepository
}

func NewCronService() CronService {
	return CronService{
		repo: repository.NewNovelRepository(),
	}
}

const (
	url                   = "https://mirror-h.org/search/country/VN/pages"
	errUnexpectedResponse = "unexpected response: %s"
)

type HTTPClient struct{}

var (
	HttpClient = HTTPClient{}
)

var backoffSchedule = []time.Duration{
	10 * time.Second,
	15 * time.Second,
	20 * time.Second,
	25 * time.Second,
	30 * time.Second,
}

func (c HTTPClient) GetRequest(pathURL string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", pathURL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	c.info(fmt.Sprintf("GET %s -> %d", pathURL, resp.StatusCode))
	if resp.StatusCode != 200 {
		respErr := fmt.Errorf(errUnexpectedResponse, resp.Status)
		return nil, respErr
	}
	return resp, nil
}

func (c HTTPClient) GetRequestWithRetries(url string) (*http.Response, error) {
	var body *http.Response
	var err error
	for _, backoff := range backoffSchedule {
		body, err = c.GetRequest(url)
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Request error: %+v\n", err)
		fmt.Fprintf(os.Stderr, "Retrying in %v\n", backoff)
		time.Sleep(backoff)
	}

	// All retries failed
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c HTTPClient) info(msg string) {
	log.Printf("[client] %s\n", msg)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

type Info struct {
	Attacker string `json:"attacker"`
	Country  string `json:"country"`
	WebUrl   string `json:"web_url"`
	Ip       string `json:"ip"`
	Date     string `json:"date"`
}

func onePage(pathURL string) ([]string, error) {
	response, err := HttpClient.GetRequestWithRetries(pathURL)
	checkError(err)
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)
	infoList := make([]string, 0)
	doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
		tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
			infoList = append(infoList, pHtml.Text())
		})
	})
	context := strings.Join(infoList, "\n")
	_ = context
	return infoList, nil
}

func (service *CronService) RunCron() (int, interface{}) {
	// insertData := model.NovelQueue{
	// 	Url:      "https://",
	// 	Source:   "wika",
	// 	Date:     "2021-06-14 12:12:00",
	// 	IsDelete: 0,
	// }
	resp, err := NovelQueue_Service.GetAllUrlInQueue()
	if err != nil {
		log.Info("RUN CRON ", "RESP : ", resp)
		// for
	}
	// const url = "https://www.wuxiaworld.com/novel/emperors-domination/emperor-chapter-6"
	// onePage(url)
	return 0, nil
}
