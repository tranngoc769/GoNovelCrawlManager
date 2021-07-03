package repository

import (
	"gonovelcrawlmanager/common/log"
	"gonovelcrawlmanager/common/model"

	IMySql "gonovelcrawlmanager/internal/sqldb/mysql"
)

type ChapterRepository struct {
}

func NewChapterRepository() ChapterRepository {
	repo := ChapterRepository{}
	return repo
}

var ChapterRepo ChapterRepository

func (repo *ChapterRepository) CreateChapter(entry model.Chapter) (interface{}, error) {
	log.Info("Story "+entry.StorySlug, entry.Slug, entry.Title)
	err := IMySql.MySqlConnector.GetConn().Model(&model.Chapter{}).Create(&entry).Error
	if err != nil {
		log.Error("Chapter Repository ", "CreateChapter", err)
		return nil, err
	}
	return entry, nil
}

func (repo *ChapterRepository) UpdateChapter(id string, data map[string]interface{}) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.Chapter{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		log.Error("Chapter Repository ", "UpdateChapter", err)
		return nil, err
	}
	return nil, nil
}

func (repo *ChapterRepository) DeleteChapter(id string) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.Chapter{}).Where("id = ?", id).Update("is_delete", 1).Error
	if err != nil {
		log.Error("Chapter Repository ", "DeleteChapter", err)
		return nil, err
	}
	return nil, nil
}

func (repo *ChapterRepository) CountChapters(search string) (int, error) {
	rows := map[string]interface{}{}
	resp := IMySql.MySqlConnector.GetConn().Table("Chapter").Where("is_delete", 0).Select("Count(id) as count")
	if search != "" {
		resp = resp.Where("content like %?%", search)
	}
	resp = resp.Take(&rows)
	if resp.Error != nil {
		log.Error("Chapter Repository ", "CountChapters", resp.Error)
		return 0, resp.Error
	}
	return int(rows["count"].(int64)), nil
}

// Chapter
func (repo *ChapterRepository) IsChapterExist(slug string, chapter_id int) (bool, model.Chapter, error) {
	rows := model.Chapter{}
	resp := IMySql.MySqlConnector.GetConn().Model(&model.Chapter{}).Where("slug = ? and chapter = ?", slug, chapter_id).Select("id").Limit(1).Take(&rows)
	if resp.RowsAffected < 1 {
		return false, rows, nil
	}
	return true, rows, nil
}
