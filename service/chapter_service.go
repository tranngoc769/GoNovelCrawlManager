package service

import (
	"gonovelcrawlmanager/common/model"
	"gonovelcrawlmanager/common/response"
	"gonovelcrawlmanager/repository"
)

var Chapter_Service ChapterService

type ChapterService struct {
	repo repository.ChapterRepository
}

func NewChapterService() ChapterService {
	return ChapterService{
		repo: repository.NewChapterRepository(),
	}
}

func (service *ChapterService) CreateChapter(entry model.Chapter) (interface{}, error) {
	resp, err := repository.ChapterRepo.CreateChapter(entry)
	if err != nil {
		return model.Chapter{}, err
	}
	return resp, nil
}

func (service *ChapterService) DeleteChapter(id string) (int, interface{}) {
	resp, err := repository.ChapterRepo.DeleteChapter(id)
	if err != nil {
		return response.NotFound()
	}
	return response.NewOKResponse(resp)
}

func (service *ChapterService) CountChapters(search string) (int, error) {
	return repository.ChapterRepo.CountChapters(search)
}
