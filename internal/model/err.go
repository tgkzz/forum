package model

import "errors"

type Err struct {
	Text_err string
	Code_err int
}

var (
	ErrInvalidData              error = errors.New("invalid email or password")
	ErrIncorrectPassword        error = errors.New("incorrect password")
	ErrInvalidId                error = errors.New("incorrect id")
	ErrInvalidPostData          error = errors.New("incorrect post data")
	ErrUnspecifiedId            error = errors.New("specify either PostId or CommentId, not both")
	ErrInvalidCategory          error = errors.New("categories has been selected wrongly")
	ErrInvalidUsername          error = errors.New("username does not exist")
	ErrEmptyComment             error = errors.New("empty comment")
	ErrInvalidComment           error = errors.New("invalid data for comment")
	ErrUsernameIsBusy           error = errors.New("username is already taken")
	ErrEmailIsBusy              error = errors.New("email is already taken")
	ErrInvalidNum               error = errors.New("invalid format: leading zeros or plus signs are not allowed")
	ErrInvalidUsernameCharacter error = errors.New("username consist of invalid characters")
	ErrInvalidExtension         error = errors.New("invalid extension")
	ErrTooLargeFile             error = errors.New("too large file")
)
