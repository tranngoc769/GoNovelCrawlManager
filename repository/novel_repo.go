package repository

import (
	"gonovelcrawlmanager/common/log"
	"gonovelcrawlmanager/common/model"

	IMySql "gonovelcrawlmanager/internal/sqldb/mysql"
)

type NovelRepository struct {
}

func NewNovelRepository() NovelRepository {
	repo := NovelRepository{}
	return repo
}

var NovelRepo NovelRepository

func (repo *NovelRepository) CreateNovel(entry model.Novel) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Create(&entry).Error
	if err != nil {
		log.Error("Novel Repository ", "CreateNovel", err)
		return nil, err
	}
	return entry, nil
}

func (repo *NovelRepository) UpdateNovel(id string, data map[string]interface{}) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		log.Error("Novel Repository ", "UpdateNovel", err)
		return nil, err
	}
	return nil, nil
}

func (repo *NovelRepository) DeleteNovel(id string) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Where("id = ?", id).Update("is_delete", 1).Error
	if err != nil {
		log.Error("Novel Repository ", "DeleteNovel", err)
		return nil, err
	}
	return nil, nil
}

func (repo *NovelRepository) GetNovelPaging(page int, limit int) ([]model.Novel, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	rows := []model.Novel{}
	resp := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Where("is_delete", 0).Select("id,title,description,crawler_href as url").Limit(limit).Offset(offset).Order("date").Find(&rows)
	if resp.Error != nil {
		log.Error("Novel Repository ", "GetNovelPaging", resp.Error)
		return []model.Novel{}, resp.Error
	}
	return rows, nil
}

func (repo *NovelRepository) CountNovels(search string) (int, error) {
	rows := map[string]interface{}{}
	resp := IMySql.MySqlConnector.GetConn().Table("novel").Where("is_delete", 0).Select("Count(id) as count")
	if search != "" {
		resp = resp.Where("content like %?%", search)
	}
	resp = resp.Take(&rows)
	if resp.Error != nil {
		log.Error("Novel Repository ", "CountNovels", resp.Error)
		return 0, resp.Error
	}
	return int(rows["count"].(int64)), nil
}

// Novel
func (repo *NovelRepository) IsStoryExist(url string) (bool, model.Novel, error) {
	rows := model.Novel{}
	resp := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Select("*").Where("crawler_href = ?", url).Limit(1).Take(&rows)
	if resp.RowsAffected < 1 {
		return false, rows, resp.Error
	}
	return true, rows, nil
}

func (repo *NovelRepository) IsExistStoryCategory(storyId string) (bool, error) {
	rows := []map[string]interface{}{}
	resp := IMySql.MySqlConnector.GetConn().Table("st_story_category").Where("story_id", storyId).Limit(1).Find(&rows)
	if resp.Error != nil {
		return false, resp.Error
	}
	if len(rows) > 0 {
		return true, nil
	}
	return false, nil
}

func (repo *NovelRepository) CreateStoryCategory(storyId string, category_id string) (bool, error) {
	resp := IMySql.MySqlConnector.GetConn().Table("st_story_category").Create([]map[string]interface{}{
		{"story_id": storyId, "category_id": category_id},
	})
	if resp.Error != nil {
		return false, resp.Error
	}
	return false, nil
}
func (repo *NovelRepository) GetStoryChapter(storyId string) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{}
	resp := IMySql.MySqlConnector.GetConn().Table("st_chapter").Select("slug", "chapter").Where("story_id", storyId).Find(&rows)
	if resp.Error != nil {
		return rows, resp.Error
	}
	return rows, nil
}
