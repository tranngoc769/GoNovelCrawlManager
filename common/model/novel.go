package model

type Novel struct {
	ID       uint   `gorm:"column:id;primarykey" json:"id"`
	Content  string `gorm:"column:content;type:longtext;not null" json:"content"`
	Url      string `gorm:"column:url;type:varchar(255);unique" json:"url"`
	Date     string `gorm:"column:date;type:longtext;" json:"date"`
	IsDelete int    `gorm:"column:is_delete;type:bit" json:"is_delete"`
}

// TableName sets the insert table name for this struct type
func (l *Novel) TableName() string {
	return "novel"
}

//

type NovelQueue struct {
	ID       uint   `gorm:"column:id;primarykey" json:"id"`
	Url      string `gorm:"column:url;type:varchar(255);unique" json:"url"`
	Source   string `gorm:"column:source;type:varchar(255);" json:"source"`
	Date     string `gorm:"column:date;type:longtext;" json:"date"`
	IsDelete int    `gorm:"column:is_delete;type:bit" json:"is_delete"`
}

// TableName sets the insert table name for this struct type
func (l *NovelQueue) TableName() string {
	return "crawl_queue"
}
