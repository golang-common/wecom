package wecom

import (
	"encoding/json"
	"net/url"
)

// UserIdGetFromCode 通过code获取用户ID
func (w *Wecom) UserIdGetFromCode(code string) (string, error) {
	var query = url.Values{}
	query.Add("code", code)
	body, err := w.get("auth/getuserinfo", query)
	if err != nil {
		return "", err
	}
	var userid string
	err = json.Unmarshal(body["userid"], &userid)
	if err != nil {
		return "", err
	}
	return userid, nil
}
