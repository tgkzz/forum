package post

import (
	"fmt"
	"forum/internal/model"
	"io"
	"log"
	"path/filepath"
	"strings"
	"unicode"
)

func validateComment(comment model.Comment) error {
	trimmedText := strings.TrimSpace(comment.Text)
	if trimmedText == "" {
		return model.ErrEmptyComment
	}

	if len(trimmedText) < 5 {
		return model.ErrInvalidComment
	}

	return nil
}

func validateFile(file model.File) error {
	_, err := file.FileGiven.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	size, err := file.FileGiven.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	_, err = file.FileGiven.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	if size > 20<<20 {
		return model.ErrTooLargeFile
	}

	filename := file.Header.Filename
	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".svg": true}
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return model.ErrInvalidExtension
	}

	return nil
}

func validatePost(post model.Post) error {
	for _, ch := range post.CategoryId {
		if ch > 4 || ch <= 0 {
			fmt.Println("")

			return model.ErrInvalidPostData
		}
	}

	log.Print(len(post.Name))

	if len(post.Name) < 5 || len(post.Name) > 50 {
		return model.ErrInvalidPostData
	}

	if len(post.Text) < 5 || len(post.Text) > 200 {
		return model.ErrInvalidPostData
	}

	post.Name = strings.TrimRight(post.Name, " ")
	post.Text = strings.TrimRight(post.Text, " ")

	if post.Name == "" {
		return model.ErrInvalidPostData
	}

	if post.Text == "" {
		return model.ErrInvalidPostData
	}

	for _, ch := range post.Text {
		if ch > unicode.MaxASCII {
			fmt.Println("3")

			return model.ErrInvalidPostData
		}
	}

	seen := make(map[string]bool)
	for _, str := range post.Category {
		if _, ok := seen[str]; ok {
			fmt.Println("4")

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
