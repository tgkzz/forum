package model

import "errors"

type Err struct {
	Text_err string
	Code_err int
}

var (
	ErrInvalidData       error = errors.New("invalid email or password")
	ErrIncorrectPassword error = errors.New("incorrect password")
	ErrInvalidId         error = errors.New("incorrect id")
	ErrInvalidPostData   error = errors.New("incorrect post data")
	ErrAlreadyGraded     error = errors.New("post already graded by user")
	ErrNoValue           error = errors.New("No value in grade or comment")
	ErrUnspecifiedId     error = errors.New("specify either PostId or CommentId, not both")
	ErrNoCategory        error = errors.New("no categories provided")
	ErrInvalidUsername   error = errors.New("username does not exist")
)
