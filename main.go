package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"gonovelcrawlmanager/common/model"
	imysql "gonovelcrawlmanager/internal/sqldb/mysql"
	mysql "gonovelcrawlmanager/internal/sqldb/mysql/driver"
	"gonovelcrawlmanager/service"

	"github.com/caarlos0/env"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

var CronService interface{}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

type MySqlConfig struct {
	Host         string
	Port         string
	Database     string
	User         string
	Password     string
	Charset      string
	PingInterval int
	MaxOpenConns int
	MaxIdleConns int
}
type DataNove struct {
	Novels []model.Novel
	Page   int
	Pages  []int
}

type DataNoveQueue struct {
	Novels []model.NovelQueue
	Page   int
	Pages  []int
}

func Novel(w http.ResponseWriter, r *http.Request) {
	limit := 4
	vars := mux.Vars(r)
	page := 1
	page_arg := vars["page"]
	if page_arg != "" {
		page, _ = strconv.Atoi(page_arg)
	}
	total_page, _ := service.Novel_Service.CountNovels("")

	log.Error("Main", "Total novels: ", total_page)
	novels, err := service.Novel_Service.GetNovelPaging(page, limit)
	if err != nil {
		log.Error("Main", "Get novels: ", err, total_page)
	}
	data := DataNove{
		Novels: novels,
		Page:   page,
		Pages:  makeRange(1, total_page/limit+1),
	}
	tmpl := template.Must(template.ParseFiles("templates/novel.html"))
	tmpl.Execute(w, data)
}

func Queue(w http.ResponseWriter, r *http.Request) {
	limit := 4
	vars := mux.Vars(r)
	page := 1
	page_arg := vars["page"]
	if page_arg != "" {
		page, _ = strconv.Atoi(page_arg)
	}
	total_page, _ := service.NovelQueue_Service.CountNovels("")
	log.Error("Main", "Total novels: ", total_page)
	novels, err := service.NovelQueue_Service.GetNovelPaging(page, limit)
	if err != nil {
		log.Error("Main", "Get novels: ", err, total_page)
	}
	data := DataNoveQueue{
		Novels: novels,
		Page:   page,
		Pages:  makeRange(1, total_page/limit+1),
	}
	tmpl := template.Must(template.ParseFiles("templates/queue.html"))
	tmpl.Execute(w, data)
}
func AddQueue(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	tmpl := template.Must(template.ParseFiles("templates/addqueue.html"))
	tmpl.Execute(w, data)
}
func AddQueuePost(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	source := r.FormValue("source")
	now := time.Now()
	dt := now.Format("2006-01-02 15:04:05")
	novel := model.NovelQueue{
		Url:      url,
		Source:   source,
		Date:     dt,
		IsDelete: 0,
	}
	data := map[string]interface{}{}
	data["backlink"] = "/add"
	data["Msg"] = "Thêm URL không thành công"
	if url == "" || source == "" {
		tmpl := template.Must(template.ParseFiles("templates/erro.html"))
		tmpl.Execute(w, data)
		return
	}
	code, _ := service.NovelQueue_Service.CreateNovel(novel)
	if code != 200 {
		tmpl := template.Must(template.ParseFiles("templates/erro.html"))
		tmpl.Execute(w, data)
		return
	}
	http.Redirect(w, r, "/add", http.StatusSeeOther)
}
func Test(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	data["backlink"] = "/add"
	tmpl := template.Must(template.ParseFiles("templates/erro.html"))
	tmpl.Execute(w, data)
}
func DeleteQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id_arg := vars["id"]
	if id_arg != "" {
		service.NovelQueue_Service.DeleteNovel(id_arg)
	}
	http.Redirect(w, r, "/queue", http.StatusSeeOther)
}

func DeleteNovel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id_arg := vars["id"]
	if id_arg != "" {
		service.Novel_Service.DeleteNovel(id_arg)
	}
	http.Redirect(w, r, "/novel", http.StatusSeeOther)
}
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	static_dir := "/statics/"
	router.PathPrefix(static_dir).Handler(http.StripPrefix(static_dir, http.FileServer(http.Dir("."+static_dir))))
	return router
}

type Config struct {
	Dir      string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port     string
	LogType  string
	LogLevel string
	LogFile  string
	LogAddr  string
	DB       string
	DBConfig
}

type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	SSLMode         string
	Timeout         int
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime int
}

var config Config

func init() {
	if err := env.Parse(&config); err != nil {
		log.Error("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	cfg := Config{
		Dir:      config.Dir,
		Port:     viper.GetString(`main.port`),
		LogType:  viper.GetString(`main.log_type`),
		LogLevel: viper.GetString(`main.log_level`),
		LogFile:  viper.GetString(`main.log_file`),
		DB:       viper.GetString(`main.db`),
	}
	if cfg.DB == "enabled" {
		cfg.DBConfig = DBConfig{
			Driver:          viper.GetString(`db.driver`),
			Host:            viper.GetString(`db.host`),
			Port:            viper.GetString(`db.port`),
			Username:        viper.GetString(`db.username`),
			Password:        viper.GetString(`db.password`),
			Database:        viper.GetString(`db.name`),
			SSLMode:         viper.GetString(`db.disable`),
			Timeout:         viper.GetInt(`db.timeout`),
			MaxOpenConns:    viper.GetInt(`db.max_open_conns`),
			MaxIdleConns:    viper.GetInt(`db.max_idle_conns`),
			MaxConnLifetime: viper.GetInt(`db.conn_max_lifetime`),
		}
	}
	config = cfg
}

func main() {
	_ = os.Mkdir(filepath.Dir(config.LogFile), 0755)
	file, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	setAppLogger(config, file)

	mysqlconfig := mysql.MySqlConfig{
		Host:         config.DBConfig.Host,
		Database:     config.DBConfig.Database,
		User:         config.DBConfig.Username,
		Password:     config.DBConfig.Password,
		Port:         config.DBConfig.Port,
		Charset:      "utf8",
		PingInterval: config.DBConfig.MaxConnLifetime,
		MaxOpenConns: config.DBConfig.MaxOpenConns,
		MaxIdleConns: config.DBConfig.MaxIdleConns,
	}
	db1 := mysql.NewMySqlConnector(mysqlconfig)
	imysql.MySqlConnector = db1
	imysql.MySqlConnector.Ping()
	mysqlconfig.Database = viper.GetString("db.name")
	db2 := mysql.NewMySqlConnector(mysqlconfig)
	imysql.MySqlGoAutodialConnector = db2
	imysql.MySqlGoAutodialConnector.Ping()
	CronService := service.NewCronService()
	service.Novel_Service = service.NewNovelService()
	service.NovelQueue_Service = service.NewNovelQueueService()

	// Cron
	crawlCron := CronService.RunCron
	s2 := gocron.NewScheduler(time.UTC)
	s2.Every(30).Seconds().Do(crawlCron)
	s2.StartAsync()
	defer s2.Clear()
	// End define
	r := NewRouter()
	router := r.PathPrefix("/").Subrouter()
	//
	router.HandleFunc("/", Novel)
	// Novel
	router.HandleFunc("/novel/page/{page}", Novel)
	router.HandleFunc("/novel", Novel)
	// Queue
	router.HandleFunc("/queue/page/{page}", Queue)
	router.HandleFunc("/queue", Queue)
	port := os.Getenv("PORT")
	// Delete queue
	router.HandleFunc("/queue/delete/{id}", DeleteQueue)
	router.HandleFunc("/queue/delete/", DeleteQueue)
	// Delete novel
	router.HandleFunc("/novel/delete/{id}", DeleteNovel)
	router.HandleFunc("/novel/delete/", DeleteNovel)
	// Add

	router.HandleFunc("/add/", AddQueue)
	router.HandleFunc("/queue_add", AddQueuePost).Methods("POST")
	// Test

	router.HandleFunc("/test", Test)
	//
	if port == "" {
		port = "3001"
	}
	log.Info("Crawl Page", "No Source - ", port)
	http.ListenAndServe(":"+port, r)
}
func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch cfg.LogType {
	case "DEFAULT":
		log.SetOutput(os.Stdout)
	case "FILE":
		if file != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.Error("main ", "Log File "+cfg.LogFile+" error")
			log.SetOutput(os.Stdout)
		}
	default:
		log.SetOutput(os.Stdout)
	}
}
