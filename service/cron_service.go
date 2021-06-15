package service

import (
	"fmt"
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/repository"
	"strconv"
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

type Info struct {
	Attacker string `json:"attacker"`
	Country  string `json:"country"`
	WebUrl   string `json:"web_url"`
	Ip       string `json:"ip"`
	Date     string `json:"date"`
}

func GetDescription(doc *goquery.Document) string {
	description := ""
	doc.Find("meta[name=description]").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "description" {
			description, _ = s.Attr("content")
		}
	})
	return strings.TrimSpace(strings.Replace(description, "online free from your Mobile, Table, PC... Novel Updates Daily", "", -1))
}
func CrawlPage(novel model.NovelQueue) {
	response, err := HttpClient.GetRequestWithRetries(novel.Url)
	if err != nil {
		log.Error("Crawl Page", "GetRequestWithRetries - ", err)
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Crawl Page", "NewDocumentFromReader - ", err)
	}
	infoList := make([]string, 0)
	switch source := novel.Source; source {
	case "wuxiaworld.com":
		novelData := model.Novel{}
		log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		name := GetDescription(doc)
		doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				if pHtml.Text() != "" {
					infoList = append(infoList, pHtml.Text())
				}
			})
		})
		context := strings.Join(infoList, "\n")
		novelData.Content = context
		novelData.Date = time.Now().Format("2006-01-02 15:04:05")
		novelData.IsDelete = 0
		novelData.Url = novel.Url
		novelData.Name = name
		Novel_Service.CreateNovel(novelData)
		id := strconv.Itoa(int(novel.ID))
		NovelQueue_Service.DeleteNovel(id)
		_ = novelData
	case "novelfull.com":
		novelData := model.Novel{}
		log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		name := GetDescription(doc)
		doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				if pHtml.Text() != "" {
					infoList = append(infoList, pHtml.Text())
				}
			})
		})
		context := strings.Join(infoList, "\n")
		novelData.Content = context
		novelData.Date = time.Now().Format("2006-01-02 15:04:05")
		novelData.IsDelete = 0
		novelData.Url = novel.Url
		novelData.Name = name
		Novel_Service.CreateNovel(novelData)
		id := strconv.Itoa(int(novel.ID))
		NovelQueue_Service.DeleteNovel(id)
		_ = novelData
	default:
		log.Error("Crawl Page", "No Source - ", source)
	}
}

func (service *CronService) RunCron() (int, interface{}) {
	resp, err := NovelQueue_Service.GetAllUrlInQueue()
	for _, novel := range resp {
		CrawlPage(novel)
	}
	if err != nil {
		log.Error("CronService ", "RunCron", err)
	}
	return 0, nil
}
