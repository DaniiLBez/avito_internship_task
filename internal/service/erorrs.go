package service

import "fmt"

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exist")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")

	ErrSlugAlreadyExists = fmt.Errorf("slug already exist")
	ErrCannotCreateSlug  = fmt.Errorf("cannot create slug")
	ErrSlugNotFound      = fmt.Errorf("slug not found")
	ErrCannotDeleteSlug  = fmt.Errorf("cannot delete slug")
	ErrCannotGetSlug     = fmt.Errorf("cannot get slug")

	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse tocken")
)
