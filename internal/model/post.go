package model

import (
	"mime/multipart"
	"time"
)

type Post struct {
	Id            int
	Name          string
	Text          string
	CreationTime  time.Time
	FormattedTime string
	filepath      string
	UserId        int
	Username      string
	Likes         int
	Dislikes      int
	CategoryId    []int
	Category      []string
	Comment       []Comment
	Categories    string
	// path for potential image
	PhotoPath string
}

type File struct {
	FileGiven multipart.File
	Header    *multipart.FileHeader
}

type Comment struct {
	Id       int
	Text     string
	PostId   int
	UserId   int
	Likes    int
	Dislikes int
	Username string
}
