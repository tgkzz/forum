package post

import (
	"forum/internal/model"
	"forum/internal/repository/post"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// worked in database, problem in post
type PostService struct {
	repo post.Post
}

type Poster interface {
	GetAllPost() ([]model.Post, error)
	GetPostById(id string) (model.Post, error)
	GetCategoryByName(strings []string) ([]int, error)
	CreatePost(post model.Post, file model.File) error
	CreateComment(comment model.Comment) error
	AddGrade(grade model.Grade) error
}

func NewPostService(repo post.Post) *PostService {
	return &PostService{
		repo: repo,
	}
}

// TODO: write logic to validate grade
func (p *PostService) AddGrade(grade model.Grade) error {
	// if err := validateGrade(grade); err != nil {
	// 	return err
	// }

	switch {
	case grade.PostId != 0 && grade.CommentId == 0:
		if err := p.repo.AddGradeToPost(grade); err != nil {
			log.Printf("Error adding or updating grade: %v", err)
			return err
		}
	case grade.CommentId != 0 && grade.PostId != 0:
		if err := p.repo.AddGradeToComment(grade); err != nil {
			return err
		}
	default:
		return model.ErrUnspecifiedId
	}

	return nil
}

func (p *PostService) GetPostById(id string) (model.Post, error) {
	num, err := strconv.Atoi(id)
	if err != nil {
		return model.Post{}, model.ErrInvalidId
	}

	post, err := p.repo.GetPostById(num)
	if err != nil {
		return post, err
	}

	post.FormattedTime = post.CreationTime.Format("2006-01-02 15:04:05")

	switch {
	case len(post.Category) == 0:
		post.Category, err = getCategoryById(post.CategoryId)
		if err != nil {
			return post, err
		}
	case len(post.CategoryId) == 0:
		post.CategoryId, err = getCategoryByName(post.Category)
		if err != nil {
			return post, err
		}
	default:
		return post, err
	}

	return post, nil
}

func (p *PostService) GetCategoryByName(strings []string) ([]int, error) {
	if len(strings) == 0 {
		return []int{4}, nil
	}
	res := []int{}

	category := map[string]int{
		"Comedy": 1,
		"Horror": 2,
		"Drama":  3,
		"Other":  4,
	}

	for _, str := range strings {
		if num, ok := category[str]; !ok {
			return []int{}, model.ErrInvalidPostData
		} else {
			res = append(res, num)
		}
	}

	return res, nil
}

func (p *PostService) GetAllPost() ([]model.Post, error) {
	posts, err := p.repo.GetAllPost()

	for _, post := range posts {
		post.FormattedTime = post.CreationTime.Format("2006-01-02 15:04:05")
	}

	return posts, err
}

func (p *PostService) CreatePost(post model.Post, file model.File) error {
	if err := validatePost(post); err != nil {
		return err
	}

	if err := validateFile(file); err != nil {
		return err
	}

	imagePath := "./template/images/" + file.Header.Filename
	if _, err := os.Stat("./template/images"); os.IsNotExist(err) {
		os.Mkdir("./template/images", os.ModePerm)
	}

	dst, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file.FileGiven); err != nil {
		return err
	}

	post.PhotoPath = strings.TrimLeft(imagePath, ".")

	return p.repo.CreatePost(post)
}

func (p *PostService) CreateComment(comment model.Comment) error {
	if err := validateComment(comment); err != nil {
		return err
	}

	return p.repo.CreateComment(comment)
}
