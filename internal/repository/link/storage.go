package linkStorage

import "errors"

var (
	ErrAliasNotFound     = errors.New("alias not found")
	ErrAliasAlreadyExist = errors.New("alias already exist")
)
