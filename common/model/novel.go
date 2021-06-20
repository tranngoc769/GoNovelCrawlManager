package model

type Novel struct {
	ID              uint   `gorm:"column:id;primarykey;"`
	Title           string `gorm:"column:title;type:varchar(255);"`
	OtherName       string `gorm:"column:other_name;type:varchar(255);"`
	Content         string `gorm:"column:content;type:text;"`
	Thumbnail       string `gorm:"column:thumbnail;type:varchar(255);"`
	IsStatus        int    `gorm:"column:is_status;type:int;"`
	AccountId       int    `gorm:"column:account_id;type:int;"`
	Slug            string `gorm:"column:slug;type:varchar(255);"`
	Description     string `gorm:"column:description;type:varchar(255);"`
	AuthorName      string `gorm:"column:author_name;type:varchar(150);"`
	AuthorSlug      string `gorm:"column:author_slug;type:varchar(150);"`
	Metatile        string `gorm:"column:meta_title;type:varchar(255);"`
	MetaDescription string `gorm:"column:meta_description;type:varchar(255);"`
	MetaKeyword     string `gorm:"column:meta_keyword;type:varchar(255);"`
	CreatedTime     string `gorm:"column:created_time;type:longtext;"`
	UpdatedTime     string `gorm:"column:updated_time;type:longtext;"`
	Url             string `gorm:"column:crawler_href;type:varchar(255);"`
}

// TableName sets the insert table name for this struct type
func (l *Novel) TableName() string {
	return "st_story"
}

type Chapter struct {
	ID          uint   `gorm:"column:id;primarykey;"`
	Title       string `gorm:"column:title;type:int;"`
	StoryId     uint   `gorm:"column:story_id;type:varchar(255);"`
	StorySlug   string `gorm:"column:story_slug;type::varchar(255);"`
	AccountId   int    `gorm:"column:account_id;type:int;"`
	IsStatus    int    `gorm:"column:is_status;type:int;"`
	IsRobot     int    `gorm:"column:is_robot;type:int;"`
	Slug        string `gorm:"column:slug;type:varchar(255);"`
	Content     string `gorm:"column:content;type:text;"`
	Chapter     uint   `gorm:"column:chapter;type:int;"`
	CreatedTime string `gorm:"column:created_time;type:longtext;"`
	UpdatedTime string `gorm:"column:updated_time;type:longtext;"`
}

// TableName sets the insert table name for this struct type
func (l *Chapter) TableName() string {
	return "st_chapter"
}

type NovelQueue struct {
	ID       uint   `gorm:"column:id;primarykey" json:"id"`
	Url      string `gorm:"column:url;type:varchar(255);unique" json:"url"`
	Source   string `gorm:"column:source;type:varchar(255);" json:"source"`
	Date     string `gorm:"column:date;type:longtext;" json:"date"`
	IsDelete int    `gorm:"column:is_delete;type:bit" json:"is_delete"`
	Category string `gorm:"column:category_id;type:int" json:"category_id"`
}

// TableName sets the insert table name for this struct type
func (l *NovelQueue) TableName() string {
	return "crawl_queue"
}

type StoryCategory struct {
	StoryId    string `gorm:"column:story_id;type:int;"`
	CategoryId string `gorm:"column:category_id;type:int;"`
}

// TableName sets the insert table name for this struct type
func (l *StoryCategory) TableName() string {
	return "st_story_category"
}
