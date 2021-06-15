package service

import (
	"errors"
	"fmt"
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/repository"
	"io"
	"strconv"
	"strings"

	"net/http"
	"os"
	"time"

	"github.com/gosimple/slug"

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

func GetMetaInfo(doc *goquery.Document) (string, string, string) {
	description := ""
	title := ""
	keywords := ""
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "description" {
			description, _ = s.Attr("content")
		}
		if name == "title" {
			title, _ = s.Attr("content")
		}
		if name == "keywords" {
			keywords, _ = s.Attr("content")
		}
	})
	return description, title, keywords
}
func GetAuthor(doc *goquery.Document) (string, string) {
	author := ""
	test := doc.Find("div[class=info]")
	_ = test
	doc.Find("div[class=info]").Each(func(i int, s *goquery.Selection) {
		s.Find("div").Each(func(j int, s2 *goquery.Selection) {
			if j == 0 {
				author = s2.Text()
			}
		})
	})
	return author, slug.Make(author)
}
func GetStoryTitle(doc *goquery.Document) (string, string) {
	title := ""
	doc.Find("h3[class=title]").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			title = s.Text()
		}
	})
	return title, slug.Make(title)
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
	novelData := model.Novel{}
	switch source := novel.Source; source {
	case "wuxiaworld.com":
		log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		name, _, _ := GetMetaInfo(doc)
		doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				if pHtml.Text() != "" {
					infoList = append(infoList, pHtml.Text())
				}
			})
		})
		context := strings.Join(infoList, "\n")
		novelData.Content = context
		novelData.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
		novelData.Url = novel.Url
		novelData.Title = name
	case "novelfull.com":
		log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		name, _, _ := GetMetaInfo(doc)
		doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				if pHtml.Text() != "" {
					infoList = append(infoList, pHtml.Text())
				}
			})
		})
		context := strings.Join(infoList, "\n")
		novelData.Content = context
		novelData.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
		novelData.Url = novel.Url
		novelData.Title = name
	default:
		log.Error("Crawl Page", "No Source - ", source)
	}
	if novelData.Content == "" && novelData.Title == "" {
		log.Error("Crawl Page", "RunCron - ", novelData.Url+" is empty !!!")
		return
	}
	Novel_Service.CreateNovel(novelData)
	id := strconv.Itoa(int(novel.ID))
	NovelQueue_Service.DeleteNovel(id)
}
func DownloadFile(fileName string, url string) error {
	//Get the response bytes from the url
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	//Create a empty file

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
func CrawlStory(novel model.NovelQueue) (model.Novel, error) {
	isExist, novel_exist, err := Novel_Service.repo.IsStoryExist(novel.Url)
	if err != nil {
		log.Error("Crawl Page", "CrawlStory - ", err)
	}
	if isExist {
		return novel_exist, nil
	}
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
	novelData := model.Novel{}
	switch source := novel.Source; source {
	case "wuxiaworld.com":
		log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		name, _, _ := GetMetaInfo(doc)
		doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				if pHtml.Text() != "" {
					infoList = append(infoList, pHtml.Text())
				}
			})
		})
		context := strings.Join(infoList, "\n")
		novelData.Content = context
		novelData.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
		novelData.Url = novel.Url
		novelData.Title = name
	case "novelfull.com":
		log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		title, title_slug := GetStoryTitle(doc)
		metaDes, metaTitle, metaKeys := GetMetaInfo(doc)
		author, author_slug := GetAuthor(doc)
		doc.Find("div.desc-text").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				if pHtml.Text() != "" {
					infoList = append(infoList, pHtml.Text())
				}
			})
		})
		desc := strings.TrimSpace(strings.Join(infoList, "\n"))
		infoList = nil
		novelData.Content = desc
		if len(desc) < 160 {
			novelData.Description = desc
		} else {
			novelData.Description = desc[0:160]
		}
		novelData.Title = title
		novelData.Slug = title_slug
		novelData.AuthorName = author
		novelData.AuthorSlug = author_slug
		novelData.IsStatus = 2
		novelData.AccountId = 1
		novelData.Url = novel.Url
		novelData.MetaDescription = metaDes
		novelData.MetaKeyword = metaKeys
		novelData.Metatile = metaTitle
		imagePath := ""
		doc.Find("div[class=book]").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("img").Each(func(indexTr int, pHtml *goquery.Selection) {
				imagePath, _ = pHtml.Attr("src")
			})
		})
		thumbPath := "uploads/" + slug.Make(time.Now().Format("2006-01-02 15:04:05")) + "_" + title_slug + ".jpg"
		imagePath = "https://novelfull.com" + imagePath
		log.Error("CronService", "Crawl Story", imagePath)
		DownloadFile(thumbPath, imagePath)
		novelData.Thumbnail = thumbPath
		_ = thumbPath
	default:
		log.Error("Crawl Page", "No Source - ", source)
	}
	if novelData.Content == "" && novelData.Title == "" {
		log.Error("Crawl Page", "RunCron - ", novelData.Url+" is empty !!!")
		return model.Novel{}, errors.New("empty")
	}
	var newNovel interface{}
	if isExist {
		return novelData, nil
	}
	novelData.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	newNovel, err = Novel_Service.CreateNovel(novelData)
	if err != nil {
		return model.Novel{}, errors.New("empty")
	}
	retNovel := newNovel.(model.Novel)
	isExist, err = Novel_Service.IsExistStoryCategory(strconv.Itoa(int(retNovel.ID)))
	if err != nil {
		return model.Novel{}, errors.New("cannot get category")
	}
	if !isExist {
		Novel_Service.repo.CreateStoryCategory(strconv.Itoa(int(retNovel.ID)), strconv.Itoa(novel.Category))
	}
	return retNovel, nil
}
func GetLastestPage(novel model.Novel, source string) int {
	response, err := HttpClient.GetRequestWithRetries(novel.Url)
	if err != nil {
		log.Error("Crawl Page", "GetRequestWithRetries - ", err)
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Crawl Page", "NewDocumentFromReader - ", err)
	}
	var infoList string
	switch source {
	case "wuxiaworld.com":
		log.Info("Crawl Page", "Source - ", source)
	case "novelfull.com":
		log.Info("Crawl Page", "Source - ", source)
		doc.Find("li[class=last]").Each(func(i int, li *goquery.Selection) {
			li.Find("a").Each(func(i int, a *goquery.Selection) {
				if i == 0 {
					infoList, _ = a.Attr("data-page")
				}
			})
		})
		_ = infoList
		page, err := strconv.Atoi(infoList)
		if err != nil {
			log.Error("GetLastestPage", "Parsing page - ", err)
			return 0
		}
		return page + 1
	default:
		log.Error("Crawl Page", "No Source - ", source)
	}
	return 1
}

func GetChapterInPages(url string, source string, page int) []string {
	if page != 0 {
		url = url + "?page=" + strconv.Itoa(page) + "&per-page=50"
	}
	response, err := HttpClient.GetRequestWithRetries(url)
	if err != nil {
		log.Error("GetChapterInPages", "GetChapterInPages - ", err)
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("GetChapterInPages", "NewDocumentFromReader - ", err)
	}
	var listChapters []string
	switch source {
	case "wuxiaworld.com":
		log.Info("GetChapterInPages", "Source - ", source)
	case "novelfull.com":
		log.Info("GetChapterInPages", "Source - ", source)
		doc.Find("ul[class=list-chapter]").Each(func(i int, li *goquery.Selection) {
			li.Find("a").Each(func(i int, a *goquery.Selection) {
				target, err := a.Attr("href")
				if !err {
					log.Error("GetChapterInPages", "NewDocumentFromReader - ", err)
				} else {
					listChapters = append(listChapters, "https://novelfull.com"+target)
				}
			})
		})
	default:
		log.Error("Crawl Page", "No Source - ", source)
	}
	return listChapters
}

func CrawlByPage(page int, novel model.Novel, source string) {

	listChapter := GetChapterInPages(novel.Url, source, page)
	_ = listChapter
}
func (service *CronService) RunCron() (int, interface{}) {
	log.Error("Crawl Page", "RunCron - ", "Running")
	resp, err := NovelQueue_Service.GetAllUrlInQueue()
	for _, novel := range resp {
		story, err := CrawlStory(novel)
		if err != nil {
			log.Error("CronService ", "CrawlStory", err)
		} else {
			page := GetLastestPage(story, novel.Source)
			CrawlByPage(page, story, novel.Source)
		}

	}
	if err != nil {
		log.Error("CronService ", "RunCron", err)
	}
	return 0, nil
}
