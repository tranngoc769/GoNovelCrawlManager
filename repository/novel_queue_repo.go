package repository

import (
	"gonovelcrawlmanager/common/log"
	"gonovelcrawlmanager/common/model"

	IMySql "gonovelcrawlmanager/internal/sqldb/mysql"
)

type NovelQueueRepository struct {
}

var NovelQueueRepo NovelQueueRepository

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
		log.Error("NovelQueueRepository ", "DeleteNovel", err)
		return nil, err
	}
	return nil, nil
}

func (repo *NovelQueueRepository) GetAllUrlInQueue() ([]model.NovelQueue, error) {
	rows := []model.NovelQueue{}
	resp := IMySql.MySqlConnector.GetConn().Model(&model.NovelQueue{}).Where("is_delete", 0).Order("date").Find(&rows)
	if resp.Error != nil {
		log.Error("NovelQueueRepository ", "GetAllUrlInQueue", resp.Error)
		return []model.NovelQueue{}, resp.Error
	}
	return rows, nil
}

func (repo *NovelQueueRepository) GetNovelPaging(page int, limit int) ([]model.NovelQueue, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	rows := []model.NovelQueue{}
	resp := IMySql.MySqlConnector.GetConn().Model(&model.NovelQueue{}).Select("id,  url, date,  source, is_delete").Where("is_delete", 0).Limit(limit).Offset(offset).Order("date").Find(&rows)
	if resp.Error != nil {
		log.Error("NovelQueueRepository ", "GetNovelPaging", resp.Error)
		return []model.NovelQueue{}, resp.Error
	}
	return rows, nil
}

func (repo *NovelQueueRepository) CountNovels(search string) (int, error) {
	rows := map[string]interface{}{}
	resp := IMySql.MySqlConnector.GetConn().Table("crawl_queue").Where("is_delete", 0).Select("Count(id) as count")
	if search != "" {
		resp = resp.Where("content like %?%", search)
	}
	resp = resp.Take(&rows)
	if resp.Error != nil {
		log.Error("NovelQueueRepository ", "CountNovels", resp.Error)
		return 0, resp.Error
	}
	return int(rows["count"].(int64)), nil
}
