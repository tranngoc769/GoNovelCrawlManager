package repository

import (
	"gonovelcrawlmanager/common/model"

	IMySql "gonovelcrawlmanager/internal/sqldb/mysql"
)

type NovelQueueRepository struct {
}

func NewNovelQueueRepository() NovelQueueRepository {
	repo := NovelQueueRepository{}
	return repo
}
func (repo *NovelQueueRepository) CreateNovel(entry model.NovelQueue) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.NovelQueue{}).Create(&entry).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (repo *NovelQueueRepository) UpdateNovel(id string, data map[string]interface{}) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Model(&model.NovelQueue{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (repo *NovelQueueRepository) DeleteNovel(id string) (interface{}, error) {
	err := IMySql.MySqlConnector.GetConn().Where("id = ?", id).Delete(&model.NovelQueue{}).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}
