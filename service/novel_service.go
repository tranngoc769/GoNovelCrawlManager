package service

import (
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/common/response"
	"gonovelcrawlmanager/repository"
)

var Novel_Service NovelService

type NovelService struct {
	repo repository.NovelRepository
}

func NewNovelService() NovelService {
	return NovelService{
		repo: repository.NewNovelRepository(),
	}
}

func (service *NovelService) CreateNovel(entry model.Novel) (interface{}, error) {
	resp, err := repository.NovelRepo.CreateNovel(entry)
	if err != nil {
		return model.Novel{}, err
	}
	return resp, nil
}

func (service *NovelService) DeleteNovel(id string) (int, interface{}) {
	resp, err := repository.NovelRepo.DeleteNovel(id)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}

func (service *NovelService) GetNovelPaging(page int, limit int) ([]model.Novel, error) {
	return repository.NovelRepo.GetNovelPaging(page, limit)
}

func (service *NovelService) CountNovels(search string) (int, error) {
	return repository.NovelRepo.CountNovels(search)
}

func (service *NovelService) IsExistStoryCategory(storyId string) (bool, error) {
	return repository.NovelRepo.IsExistStoryCategory(storyId)
}

func (service *NovelService) CreateStoryCategory(storyId string, cate string) (bool, error) {
	return repository.NovelRepo.CreateStoryCategory(storyId, cate)
}
