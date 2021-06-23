package service

import (
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/common/response"
	"gonovelcrawlmanager/repository"
)

var NovelQueue_Service NovelQueueService

type NovelQueueService struct {
	repo repository.NovelQueueRepository
}

func NewNovelQueueService() NovelQueueService {
	return NovelQueueService{
		repo: repository.NewNovelQueueRepository(),
	}
}

func (service *NovelQueueService) CreateNovel(entry model.NovelQueue) (int, interface{}) {
	resp, err := repository.NovelQueueRepo.CreateNovel(entry)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}
func (service *NovelQueueService) UpdateNovel(id string, data map[string]interface{}) (int, interface{}) {
	resp, err := repository.NovelQueueRepo.UpdateNovel(id, data)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}

func (service *NovelQueueService) DeleteNovel(id string) (int, interface{}) {
	resp, err := repository.NovelQueueRepo.DeleteNovel(id)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}

func (service *NovelQueueService) GetAllUrlInQueue() ([]model.NovelQueue, error) {
	return repository.NovelQueueRepo.GetAllUrlInQueue()
}

func (service *NovelQueueService) GetNovelPaging(page int, limit int) ([]model.NovelQueue, error) {
	return repository.NovelQueueRepo.GetNovelPaging(page, limit)
}

func (service *NovelQueueService) CountNovels(search string) (int, error) {
	return repository.NovelQueueRepo.CountNovels(search)
}

func (service *NovelQueueService) GetAllCategory() ([]map[string]interface{}, error) {
	return repository.NovelQueueRepo.GetAllCategory()
}

func (service *NovelQueueService) IsQueueExist(url string) (bool, error) {
	return repository.NovelQueueRepo.IsQueueExist(url)
}

func (repo *NovelQueueService) GetOtherName() ([]map[string]interface{}, error) {
	return repository.NovelQueueRepo.GetOtherName()
}
