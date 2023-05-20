/**
 * @Author: DPY
 * @Description: 企业微信客户端模式
 * @File:  client.go
 * @Version: 1.0.0
 * @Date: 2022/2/8 14:49
 */

package wecom

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

const (
	apiHost   = `qyapi.weixin.qq.com`
	basePath  = `/cgi-bin`
	https     = "https"
	pathToken = `/gettoken`
)

func New(corpid, secret string) *Wecom {
	return &Wecom{
		corpID:     corpid,
		corpSecret: secret,
	}
}

type Wecom struct {
	token       string
	tokenExpire time.Time
	corpID      string
	corpSecret  string
	debug       bool
}

func (w *Wecom) Debug() {
	w.debug = true
}

// Check 本地校验对象的可用性
func (w *Wecom) Check() error {
	if w.token == "" {
		return errors.New("token is nil")
	}
	if time.Now().Unix() > w.tokenExpire.Unix() {
		return errors.New("token expired")
	}
	if w.corpID == "" || w.corpSecret == "" {
		return errors.New("object not initialized")
	}
	return nil
}

// CheckAndAuth 如果校验不通过则发起认证
func (w *Wecom) CheckAndAuth() error {
	if w.Check() == nil {
		return nil
	}
	err := w.Auth()
	if err != nil {
		return err
	}
	return nil
}

func (w *Wecom) NewAesKey() string {
	return RandString(32)
}

// Auth 获取并设置企业微信token
func (w *Wecom) Auth() error {
	queryVal := url.Values{}
	queryVal.Add("corpid", w.corpID)
	queryVal.Add("corpsecret", w.corpSecret)

	tokenUrl := &url.URL{
		Scheme: https,
		Host:   apiHost,
		Path:   path.Join(basePath, pathToken),
	}
	tokenUrl.RawQuery = queryVal.Encode()

	header := http.Header{}
	header.Set("Content-type", "application/json")

	req := &http.Request{
		Method: http.MethodGet,
		URL:    tokenUrl,
		Header: header,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var token struct {
		Errcode     int    `json:"errcode"`      // 企业微信错误码，正常时为0
		Errmsg      string `json:"errmsg"`       // 企业微信错误信息，正常时为"ok"
		AccessToken string `json:"access_token"` // Token内容
		ExpiresIn   int    `json:"expires_in"`   // token有效时间（秒）
	}
	err = json.Unmarshal(respBytes, &token)
	if err != nil {
		return err
	}
	if token.Errcode != 0 {
		return errors.New(token.Errmsg)
	}
	w.token = token.AccessToken
	w.tokenExpire = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	return nil
}

// get 通用的get方法,返回body内容或错误
// p - 请求路径
// query - 请求的url-query
func (w *Wecom) get(p string, query url.Values) (map[string]json.RawMessage, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Add("access_token", w.token)
	if w.debug {
		query.Add("debug", "1")
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "https",
			Host:     apiHost,
			Path:     path.Join(basePath, p),
			RawQuery: query.Encode(),
		},
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	mpBytes, err := parseResponseBody(bodyBytes)
	if err != nil {
		return nil, err
	}
	return mpBytes, nil
}

// post 通用post请求,返回body内容或错误
func (w *Wecom) post(p string, b interface{}) (map[string]json.RawMessage, error) {
	query := url.Values{}
	query.Add("access_token", w.token)
	if w.debug {
		query.Add("debug", "1")
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme:   "https",
			Host:     apiHost,
			Path:     path.Join(basePath, p),
			RawQuery: query.Encode(),
		},
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	reqBodyMap := make(map[string]json.RawMessage)
	err = json.Unmarshal(bb, &reqBodyMap)
	if err != nil {
		return nil, err
	}
	//reqBodyMap["access_token"] = []byte(w.token)
	body, _ := json.Marshal(reqBodyMap)
	req.Body = io.NopCloser(bytes.NewReader(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	mpBytes, err := parseResponseBody(bodyBytes)
	if err != nil {
		return nil, err
	}
	return mpBytes, nil
}

func parseResponseBody(body []byte) (map[string]json.RawMessage, error) {
	var r = make(map[string]json.RawMessage)
	err := json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	var ecode Error
	var rb = make(map[string]json.RawMessage)
	for k, v := range r {
		if k == "errcode" {
			ec, _ := strconv.Atoi(string(v))
			ecode.Errcode = ec
			continue
		}
		if k == "errmsg" {
			ecode.Errmsg = string(v)
			continue
		}
		rb[k] = v
	}
	err = ecode.Check()
	if err != nil {
		return nil, err
	}
	return rb, nil
}

func umarshalObject(data map[string]json.RawMessage, obj interface{}) error {
	datab, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(datab, obj)
	if err != nil {
		return err
	}
	return nil
}
