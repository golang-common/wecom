package wecom

import (
	"errors"
	"fmt"
)

type Error struct {
	Errcode int    `json:"errcode" xml:"ErrCode"`
	Errmsg  string `json:"errmsg" xml:"ErrMsg"`
}

func (e Error) Check() error {
	if e.Errcode == 0 {
		return nil
	}
	if e.Errmsg != "" {
		return errors.New(e.Errmsg)
	}
	return errors.New(fmt.Sprintf("request error with code %d", e.Errcode))
}
