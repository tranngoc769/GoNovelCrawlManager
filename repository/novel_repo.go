package repository

import (
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
		return nil, err
	}
	return nil, nil
}

func (repo *NovelRepository) UpdateNovel(id string, data map[string]interface{}) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (repo *NovelRepository) DeleteNovel(id string) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Where("id = ?", id).Update("is_delete", 1).Error
	if err != nil {
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
	resp := IMySql.MySqlConnector.GetConn().Model(&model.Novel{}).Where("is_delete", 0).Select("id, caption, url, date, SUBSTRING(content, 1,10) as content, is_delete").Limit(limit).Offset(offset).Order("date").Find(&rows)
	if resp.Error != nil {
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
		return 0, resp.Error
	}
	return int(rows["count"].(int64)), nil
}
