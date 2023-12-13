package filter

import (
	"forum/internal/model"
	"forum/internal/repository/filter"
)

type FilterService struct {
	repo filter.Filter
}

type Filterer interface {
	GetUserPosts(userId int) ([]model.Post, error)
	FilterByCategory(categories []int) ([]model.Post, error)
	FilterByLikes(userId int) ([]model.Post, error)
}

func NewFilterService(repository filter.Filter) *FilterService {
	return &FilterService{
		repo: repository,
	}
}

func (f *FilterService) GetUserPosts(userId int) ([]model.Post, error) {
	posts, err := f.repo.GetUserPostsById(userId)

	for _, post := range posts {
		post.FormattedTime = post.CreationTime.Format("2006-01-02 15:04:05")
	}

	return posts, err
}

func (f *FilterService) FilterByCategory(categories []int) ([]model.Post, error) {
	return f.repo.GetPostsByCategory(categories)
}

func (f *FilterService) FilterByLikes(userId int) ([]model.Post, error) {
	return f.repo.GetUsersLikedPost(userId)
}
