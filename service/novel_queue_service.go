package service

import (
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/common/response"
	"gonovelcrawlmanager/repository"
)

type NovelQueueService struct {
	repo repository.NovelRepository
}

func NewNovelQueueService() NovelQueueService {
	return NovelQueueService{
		repo: repository.NewNovelRepository(),
	}
}

func (service *NovelQueueService) CreateNovel(entry model.Novel) (int, interface{}) {
	resp, err := repository.NovelRepo.CreateNovel(entry)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}
func (service *NovelQueueService) UpdateNovel(id string, data map[string]interface{}) (int, interface{}) {
	resp, err := repository.NovelRepo.UpdateNovel(id, data)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}

func (service *NovelQueueService) DeleteNovel(id string) (int, interface{}) {
	resp, err := repository.NovelRepo.DeleteNovel(id)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}
