package service

import (
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/common/response"
	"gonovelcrawlmanager/repository"
)

type NovelService struct {
	repo repository.NovelRepository
}

func NewNovelService() NovelService {
	return NovelService{
		repo: repository.NewNovelRepository(),
	}
}

func (service *NovelService) CreateNovel(entry model.Novel) (int, interface{}) {
	resp, err := repository.NovelRepo.CreateNovel(entry)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}
func (service *NovelService) UpdateNovel(id string, data map[string]interface{}) (int, interface{}) {
	resp, err := repository.NovelRepo.UpdateNovel(id, data)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}

func (service *NovelService) DeleteNovel(id string) (int, interface{}) {
	resp, err := repository.NovelRepo.DeleteNovel(id)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}
