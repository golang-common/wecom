package wecom

import "net/url"

// UserIdGetFromCode 通过code获取用户ID
func (w *Wecom) UserIdGetFromCode(code string) (string, error) {
	var query = url.Values{}
	query.Add("code", code)
	body, err := w.get("", query)
	if err != nil {
		return "", err
	}
	return string(body["userid"]), nil
}
