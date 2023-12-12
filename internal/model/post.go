package model

import "time"

type Post struct {
	Id            int
	Name          string
	Text          string
	CreationTime  time.Time
	FormattedTime string
	UserId        int
	Username      string
	Likes         int
	Dislikes      int
	CategoryId    []int
	Category      []string
	Comment       []Comment
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
