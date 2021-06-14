package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

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

var NovelService interface{}
var NovelQueueService interface{}

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

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/main.html"))
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		username = "tnquang769"
	}
	log.Info("Index View ", "User : ", username)
	is := map[string]interface{}{}
	is["ok"] = "ok"
	tmpl.Execute(w, is)
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
	NovelService := service.NewNovelService()
	NovelQueueService := service.NewNovelQueueService()
	_ = NovelService
	_ = NovelQueueService

	// Cron
	crawlCron := CronService.RunCron
	s2 := gocron.NewScheduler(time.UTC)
	s2.Every(500).Seconds().Do(crawlCron)
	s2.StartAsync()
	defer s2.Clear()
	// End define
	r := NewRouter()
	router := r.PathPrefix("/").Subrouter()
	router.HandleFunc("/{username}", Index)
	router.HandleFunc("/", Index)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	http.ListenAndServe(":"+port, r)
	fmt.Print("Server is running " + port)
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
	default:
		log.SetOutput(os.Stdout)
	}
}
