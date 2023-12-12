package post

import (
	"forum/internal/model"
	"log"
	"unicode"
)

func validateGraded(grade model.Grade) error {
	return nil
}

func validatePost(post model.Post) error {
	for _, ch := range post.CategoryId {
		if ch > 4 || ch <= 0 {
			return model.ErrInvalidPostData
		}
	}

	if post.Name == "" {
		return model.ErrInvalidPostData
	}

	if post.Text == "" {
		return model.ErrInvalidPostData
	}

	for _, ch := range post.Text {
		if ch > unicode.MaxASCII {
			return model.ErrInvalidPostData
		}
	}

	seen := make(map[string]bool)
	for _, str := range post.Category {
		if _, ok := seen[str]; ok {
			return model.ErrInvalidPostData
		}
		seen[str] = true
	}

	return nil
}

func getCategoryByName(strings []string) ([]int, error) {
	res := []int{}

	category := map[string]int{
		"Comedy": 1,
		"Horror": 2,
		"Drama":  3,
		"Other":  4,
	}

	for _, str := range strings {
		if num, ok := category[str]; !ok {
			log.Print("djasd")
			return []int{}, model.ErrInvalidData
		} else {
			res = append(res, num)
		}
	}

	return res, nil
}

func getCategoryById(nums []int) ([]string, error) {
	res := []string{}

	category := map[int]string{
		1: "Comedy",
		2: "Horror",
		3: "Drama",
		4: "Other",
	}

	for _, num := range nums {
		if str, ok := category[num]; !ok {
			log.Print("asd")
			return []string{}, model.ErrInvalidData
		} else {
			res = append(res, str)
		}
	}

	return res, nil
}
