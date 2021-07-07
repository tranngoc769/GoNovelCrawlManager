package service

import (
	"errors"
	"fmt"
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/repository"
	"io"
	"regexp"
	"strconv"
	"strings"

	"net/http"
	"os"
	"time"

	"github.com/gosimple/slug"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
)

var WUXIA = "wuxiaworld.com"
var NOVEL_FULL = "novelfull.com"
var Re *regexp.Regexp
var IMG_FOLER = "/home/novelhot.net/public_html/public/media/"

func GetChapterIdFromSlug(slug string) int {
	listNumber := Re.FindAllString(slug, -1)
	if len(listNumber) < 1 {
		return 0
	}
	chapter_id := Re.FindAllString(slug, -1)[0]
	id, _ := strconv.Atoi(chapter_id)
	return id
}
func GetSlugFromURL(url string) string {
	spliter := strings.Split(url, "/")
	if len(spliter) > 1 {
		ok := strings.Replace(spliter[len(spliter)-1], ".html", "", -1)
		ok = strings.Replace(ok, ".htm", "", -1)
		ok = strings.Replace(ok, "/", "", -1)
		return ok
	}
	return ""
}

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

func removeEmoji(in string) string {
	var emojiRx = regexp.MustCompile(`[^\x00-\x7F]`)
	var s = emojiRx.ReplaceAllString(in, "")
	return s
}
func GetMetaInfo(doc *goquery.Document) (string, string, string) {
	description := ""
	title := ""
	keywords := ""
	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
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
func GetPreview(url string, source string) (map[string]interface{}, error) {
	preview := map[string]interface{}{}
	response, err := HttpClient.GetRequestWithRetries(url)
	if err != nil {
		log.Error("Crawl Page", "GetRequestWithRetries - ", err)
		return preview, err
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Crawl Page", "NewDocumentFromReader - ", err)
	}
	switch source {
	case WUXIA:
		{
			doc.Find("div[class=novel-left] > a").Each(func(index int, tableHtml *goquery.Selection) {
				tableHtml.Find("img").Each(func(indexTr int, pHtml *goquery.Selection) {
					preview["image"], _ = pHtml.Attr("src")
				})
			})
			doc.Find("div[class=novel-body]").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					preview["info"], _ = s.Html()
				}
			})
		}
	case NOVEL_FULL:
		{
			doc.Find("div[class=book]").Each(func(index int, tableHtml *goquery.Selection) {
				tableHtml.Find("img").Each(func(indexTr int, pHtml *goquery.Selection) {
					preview["image"], _ = pHtml.Attr("src")
					preview["image"] = "https://" + NOVEL_FULL + preview["image"].(string)
				})
			})
			doc.Find("div[class=info]").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					preview["info"], _ = s.Html()
				}
			})
		}
	}
	return preview, nil
}
func GetAuthor(doc *goquery.Document) (string, string) {
	author := ""
	test := doc.Find("div[class=info]")
	_ = test
	doc.Find("div[class=info]").Each(func(_ int, s *goquery.Selection) {
		s.Find("div").Each(func(j int, s2 *goquery.Selection) {
			if j == 0 {
				author = s2.Text()
			}
		})
	})
	return author, slug.Make(author)
}
func GetStoryTitle(doc *goquery.Document, source string) (string, string) {
	title := ""
	switch source {
	case WUXIA:
		doc.Find("h2").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				title = s.Text()
			}
		})
		return title, ""
	case NOVEL_FULL:
		doc.Find("h3[class=title]").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				title = s.Text()
			}
		})
		return title, slug.Make(title)
	default:
		log.Info("Default", "No", "No")
		return "", ""
	}
}

// Get Chapter Title
func GetChapterTile(source string, doc *goquery.Document) (string, string) {
	title := ""
	switch source {
	case WUXIA:
		log.Info("Crawl Page", "Source - ", source)

	case NOVEL_FULL:
		title = doc.Find("h3").First().Text()
		if title == "" {
			title = doc.Find("span[class=chapter-text]").First().Text()
		}
	default:
		_ = "ok"
	}
	return title, slug.Make(title)
}
func CrawlChapter(source string, url string) (model.Chapter, error) {
	response, err := HttpClient.GetRequestWithRetries(url)
	if err != nil {
		log.Error("CrawlChapter", "GetRequestWithRetries - ", err)
		return model.Chapter{}, err
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Crawl Page", "NewDocumentFromReader - ", err)
		return model.Chapter{}, err
	}
	infoList := make([]string, 0)
	chapterData := model.Chapter{}
	switch source {
	case WUXIA:
		// log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		name, _, _ := GetMetaInfo(doc)
		doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				ret, _ := pHtml.Html()
				if pHtml.Text() != "" {
					infoList = append(infoList, "<p>"+ret+"</p>")
				}
			})
		})
		context := strings.Join(infoList, "\n")
		chapterData.Content = context
		chapterData.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
		chapterData.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
		chapterData.Title = name
	case NOVEL_FULL:
		// log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		_, _, _ = GetMetaInfo(doc)
		title, _ := GetChapterTile(source, doc)
		doc.Find("div#chapter-content").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				ret, _ := pHtml.Html()
				if pHtml.Text() != "" {
					infoList = append(infoList, "<p>"+ret+"</p>")
				}
			})
		})
		context := strings.Join(infoList, "\n")
		chapterData.Content = context
		chapterData.Title = title
	default:
		log.Error("Crawl Page", "No Source - ", source)
	}
	chapterData.IsStatus = 1
	chapterData.AccountId = 8
	chapterData.IsRobot = 1
	chapterData.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	chapterData.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
	if chapterData.Content == "" && chapterData.Title == "" {
		log.Error("Crawl Page", "RunCron - ", url+" is empty !!!")
		return model.Chapter{}, err
	}
	return chapterData, nil
	// Chapter_Service.CreateChapter(chapterData)
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
func CrawlStory(novel model.NovelQueue) (bool, model.Novel, []string, error) {
	urlList := []string{}
	var onGo bool
	response, err := HttpClient.GetRequestWithRetries(novel.Url)
	if err != nil {
		log.Error("Crawl Page", "GetRequestWithRetries - ", err)
		return false, model.Novel{}, urlList, err
	}
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Crawl Page", "NewDocumentFromReader - ", err)
		return false, model.Novel{}, urlList, err
	}
	var infoList []string
	switch source := novel.Source; source {
	case WUXIA:
		onGo = true
	case NOVEL_FULL:

		doc.Find("div[class=info]").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				infoList = append(infoList, s.Text())
			}
		})
		desc := strings.TrimSpace(strings.Join(infoList, "\n"))
		if strings.Contains(desc, "Status:Completed") {
			_ = desc
			onGo = false
			// return false, model.Novel{}, urlList, nil
		} else {
			onGo = true
		}
	}
	infoList = nil
	novel_slug := GetSlugFromURL(novel.Url)
	isExist, novel_exist, _ := Novel_Service.repo.IsStoryExist(novel.Url, novel_slug)
	if isExist {
		switch source := novel.Source; source {
		case WUXIA:
			doc.Find("li.chapter-item>a").Each(func(index int, tableHtml *goquery.Selection) {
				test, _ := tableHtml.Attr("href")
				urlList = append(urlList, "https://www.wuxiaworld.com"+test)
			})
		case NOVEL_FULL:
			_ = onGo
		}
		return onGo, novel_exist, urlList, nil
	}
	novelData := model.Novel{}
	switch source := novel.Source; source {
	case WUXIA:
		// log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		// GET URL LIST
		doc.Find("li.chapter-item>a").Each(func(index int, tableHtml *goquery.Selection) {
			test, _ := tableHtml.Attr("href")
			urlList = append(urlList, "https://www.wuxiaworld.com"+test)
		})
		metaDes, metaKeys, metaTitle := GetMetaInfo(doc)
		novelData.MetaDescription = metaDes
		novelData.MetaKeyword = metaKeys
		novelData.Metatile = metaTitle

		infoList = nil
		doc.Find("div.fr-view").Each(func(index int, tableHtml *goquery.Selection) {
			// test := tableHtml.Text()
			ret, _ := tableHtml.Html()
			infoList = append(infoList, ret)
		})

		context := strings.Join(infoList, "\n")
		infoList = nil
		novelData.Content = strings.TrimSpace(context)
		doc.Find("div.novel-body >dd").Each(func(index int, tableHtml *goquery.Selection) {
			test := tableHtml.Text()
			infoList = append(infoList, test)
		})
		title, _ := GetStoryTitle(doc, source)
		novelData.Title = title
		title_slug := GetSlugFromURL(novel.Url)
		novelData.Slug = title_slug
		novelData.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
		novelData.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
		novelData.Url = novel.Url
		author := strings.Join(infoList, "-")
		novelData.AuthorName = author
		novelData.AuthorSlug = slug.Make(author)
		imagePath := ""
		doc.Find("div.novel-left > a").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("img").Each(func(indexTr int, pHtml *goquery.Selection) {
				imagePath, _ = pHtml.Attr("src")
			})
		})
		_ = imagePath

		thumbPath := IMG_FOLER + title_slug + ".jpg"
		log.Error("CronService", "Crawl Story", imagePath)
		DownloadFile(thumbPath, imagePath)
		novelData.Thumbnail = "/" + title_slug + ".jpg"
		_ = thumbPath
	case NOVEL_FULL:
		// log.Info("Crawl Page", "Source - ", source)
		// Get Caption
		title, title_slug := GetStoryTitle(doc, source)
		metaDes, metaTitle, metaKeys := GetMetaInfo(doc)
		author, author_slug := GetAuthor(doc)
		doc.Find("div.desc-text").Each(func(index int, tableHtml *goquery.Selection) {
			tableHtml.Find("p").Each(func(indexTr int, pHtml *goquery.Selection) {
				ret, _ := pHtml.Html()
				if pHtml.Text() != "" {
					infoList = append(infoList, "<p>"+ret+"</p>")
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
		novelData.IsStatus = 1
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
		thumbPath := IMG_FOLER + title_slug + ".jpg"
		imagePath = "https://novelfull.com" + imagePath
		log.Error("CronService", "Crawl Story", imagePath)
		DownloadFile(thumbPath, imagePath)
		novelData.Thumbnail = "/" + title_slug + ".jpg"
		_ = thumbPath
	default:
		log.Error("Crawl Page", "No Source - ", source)
	}
	if novelData.Content == "" && novelData.Title == "" {
		log.Error("Crawl Page", "RunCron - ", novelData.Url+" is empty !!!")
		return onGo, model.Novel{}, urlList, errors.New("empty")
	}
	var newNovel interface{}
	if isExist {
		return onGo, novelData, urlList, nil
	}
	novelData.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	novelData.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
	newNovel, err = Novel_Service.CreateNovel(novelData)
	if err != nil {
		return onGo, model.Novel{}, urlList, errors.New("empty")
	}
	retNovel := newNovel.(model.Novel)
	isExist, err = Novel_Service.IsExistStoryCategory(strconv.Itoa(int(retNovel.ID)))
	if err != nil {
		return onGo, model.Novel{}, urlList, errors.New("cannot get category")
	}
	if !isExist {
		caterory_list := strings.Split(novel.Category, ",")
		for _, ct := range caterory_list {
			Novel_Service.repo.CreateStoryCategory(strconv.Itoa(int(retNovel.ID)), ct)
		}
	}
	return onGo, retNovel, urlList, nil
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
	case WUXIA:
		log.Info("Crawl Page", "Source - ", source)
	case NOVEL_FULL:
		// log.Info("Crawl Page", "Source - ", source)
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
	case WUXIA:
		log.Info("GetChapterInPages", "Source - ", source)
	case NOVEL_FULL:
		// log.Info("GetChapterInPages", "Source - ", source)
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

func CrawlByList(urlList []string, novel model.Novel, source string, onGoing bool, queueID string, isComplete bool) {
	if len(urlList) == 0 {
		return
	}
	for i := len(urlList) - 1; i >= 0; i-- {
		url := urlList[i]
		slug := GetSlugFromURL(url)
		if slug == "" {
			log.Error("CrawlByPageRev", "GetSlugFromURL - ", "Slug is empty")
			return
		}
		id := GetChapterIdFromSlug(slug)
		if id == 0 {
			log.Error("CrawlByPageRev", "GetChapterIdFromSlug - ", "Chapterid is empty")
		}
		isExist, _, err := Chapter_Service.repo.IsChapterExist(slug, id)
		if err != nil {
			log.Error("CrawlByPageRev", "IsChapterExist - ", err)
		}
		if isExist {
			log.Info("CrawlByPageRev", "Chapter is exist! Break", "")
		} else {
			chapterDate, err := CrawlChapter(source, url)
			if err != nil {
				log.Info("CrawlByPageRev", "IsChapterExist - ", err)
			}
			chapterDate.Slug = slug
			chapterDate.Chapter = uint(id)
			chapterDate.StoryId = novel.ID
			chapterDate.StorySlug = novel.Slug
			chapterDate.Content = removeEmoji(chapterDate.Content)
			_, err = Chapter_Service.repo.CreateChapter(chapterDate)
			if err != nil {
				log.Error("CrawlByPageRev", "CreateChapter - ", err)
				// return
			}
		}
	}
}
func CrawlByPageRev(page int, novel model.Novel, source string, onGoing bool, queueID string, isComplete bool) {
	if page == 0 || page == -1 {
		if !onGoing {
			NovelQueue_Service.DeleteNovel(queueID)
			return
		}
		NovelQueue_Service.MakeQueueComplete(queueID)
		return
	}
	listChapter := GetChapterInPages(novel.Url, source, page)
	// Ex : page 10 : Chapter 50 --> 40
	for i := len(listChapter) - 1; i >= 0; i-- {
		url := listChapter[i]
		slug := GetSlugFromURL(url)
		if slug == "" {
			log.Error("CrawlByPageRev", "GetSlugFromURL - ", "Slug is empty")
			return
		}
		id := GetChapterIdFromSlug(slug)
		if id == 0 {
			log.Error("CrawlByPageRev", "GetChapterIdFromSlug - ", "Chapterid is empty")
		}
		isExist, _, err := Chapter_Service.repo.IsChapterExist(slug, id)
		if err != nil {
			log.Error("CrawlByPageRev", "IsChapterExist - ", err)
			return
		}
		if isExist {
			if isComplete {
				log.Warn("CrawlByPageRev", "Chapter is exist and story complete! Stop", "")
				return
			}
			break
		} else {
			chapterDate, err := CrawlChapter(source, listChapter[i])
			if err != nil {
				log.Info("CrawlByPageRev", "IsChapterExist - ", err)
			}
			chapterDate.Slug = slug
			chapterDate.Chapter = uint(id)
			chapterDate.StoryId = novel.ID
			chapterDate.StorySlug = novel.Slug
			chapterDate.Content = removeEmoji(chapterDate.Content)
			_, err = Chapter_Service.repo.CreateChapter(chapterDate)
			if err != nil {
				log.Error("CrawlByPageRev", "CreateChapter - ", err)
				// return
			}
		}
	}
	CrawlByPageRev((page - 1), novel, source, onGoing, queueID, isComplete)
}
func GetAllChapterURL(page int, novel model.Novel, source string) (map[string]interface{}, []string) {
	data := map[string]interface{}{}
	list_slug := []string{}
	_ = data
	inc := 0
	for i := 1; i <= page; i++ {
		listChapter := GetChapterInPages(novel.Url, source, i)
		for j := 0; j < len(listChapter); j++ {
			url := listChapter[j]
			slug := GetSlugFromURL(url)
			data[slug] = url
			list_slug = append(list_slug, slug)
			inc += 1
		}
	}
	return data, list_slug
}
func RemoveSlugFromList(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
func (service *CronService) RunCron() (int, interface{}) {
	// log.Error("Crawl Page", "RunCron - ", "Running")
	resp, err := NovelQueue_Service.GetAllUrlInQueue()
	for _, novel := range resp {
		onGoing, story, urlList, err := CrawlStory(novel)
		if err != nil {
			log.Error("CronService ", "CrawlStory", err)
		} else {
			isComplete, _ := NovelQueue_Service.IsMakeCompleted(novel.Url)
			switch novel.Source {
			case WUXIA:
				if isComplete {
					CrawlByList(urlList, story, novel.Source, onGoing, strconv.Itoa(int(novel.ID)), isComplete)
				} else {
					listChapter_db, _ := Novel_Service.GetStoryChapter(strconv.Itoa(int(story.ID)))
					_ = listChapter_db
					listChapter_site := urlList
					for i := 0; i < len(listChapter_db); i++ {
						story.Url = strings.Replace(story.Url, "https://wuxiaworld.com", "https://www.wuxiaworld.com", -1)
						tmp := story.Url + "/" + listChapter_db[i]["slug"].(string)
						ix := indexOf(tmp, listChapter_site)
						if ix > -1 {
							listChapter_site = RemoveSlugFromList(listChapter_site, ix)
						}
					}
					// Crawl list diff
					for i := 0; i < len(listChapter_site); i++ {
						sl := GetSlugFromURL(listChapter_site[i])
						chap_slug := sl
						id := GetChapterIdFromSlug(chap_slug)
						chapterDate, err := CrawlChapter(WUXIA, listChapter_site[i])
						if err != nil {
							log.Info("CrawlByPageRev", "IsChapterExist - ", err)
						}
						chapterDate.Slug = chap_slug
						chapterDate.Chapter = uint(id)
						chapterDate.StoryId = story.ID
						chapterDate.StorySlug = story.Slug
						chapterDate.Content = removeEmoji(chapterDate.Content)
						_, err = Chapter_Service.repo.CreateChapter(chapterDate)
						if err != nil {
							log.Error("CrawlByPageRev", "CreateChapter - ", err)
						}
					}
				}
			case NOVEL_FULL:
				// true --> only get update chap
				page := GetLastestPage(story, novel.Source)
				if isComplete {
					CrawlByPageRev(page, story, novel.Source, onGoing, strconv.Itoa(int(novel.ID)), isComplete)
				} else {
					listChapter_site, list_slug := GetAllChapterURL(page, story, novel.Source)
					listChapter_db, _ := Novel_Service.GetStoryChapter(strconv.Itoa(int(story.ID)))
					_ = listChapter_db
					_ = listChapter_site
					for i := 0; i < len(listChapter_db); i++ {
						tmp := listChapter_db[i]["slug"].(string)
						ix := indexOf(tmp, list_slug)
						if ix > -1 {
							list_slug = RemoveSlugFromList(list_slug, ix)
						}
					}
					// Crawl list diff
					for i := 0; i < len(list_slug); i++ {
						chap_slug := list_slug[i]
						id := GetChapterIdFromSlug(chap_slug)
						chapterDate, err := CrawlChapter(NOVEL_FULL, listChapter_site[chap_slug].(string))
						if err != nil {
							log.Info("CrawlByPageRev", "IsChapterExist - ", err)
						}
						chapterDate.Slug = chap_slug
						chapterDate.Chapter = uint(id)
						chapterDate.StoryId = story.ID
						chapterDate.StorySlug = story.Slug
						chapterDate.Content = removeEmoji(chapterDate.Content)
						_, err = Chapter_Service.repo.CreateChapter(chapterDate)
						if err != nil {
							log.Error("CrawlByPageRev", "CreateChapter - ", err)
						}
					}
					log.Info("CronService ", "CrawlStory Chapter Finish", len(list_slug))
					NovelQueue_Service.MakeQueueComplete(strconv.Itoa(int(novel.ID)))
					if !onGoing {
						NovelQueue_Service.DeleteNovel(strconv.Itoa(int(novel.ID)))
					}
				}
			}
		}
	}
	if err != nil {
		log.Error("CronService ", "RunCron", err)
	}
	return 0, nil
}

func (service *CronService) TEST() {
	url := "https://novelfull.com/im-really-a-superstar/chapter-403-beating-up-three-people-consecutively.html"
	slug := GetSlugFromURL(url)
	id := GetChapterIdFromSlug(slug)
	src := "novelfull.com"
	chapterDate, err := CrawlChapter(src, url)
	if err != nil {
		log.Info("CrawlByPageRev", "IsChapterExist - ", err)
	}

	chapterDate.Slug = slug
	chapterDate.Chapter = uint(id)
	chapterDate.StoryId = 1
	chapterDate.StorySlug = "test"
	// chapterDate.Content = removeEmoji(chapterDate.Content)
	_, _ = Chapter_Service.repo.CreateChapter(chapterDate)
}
