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
