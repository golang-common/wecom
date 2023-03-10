/**
 * @Author: DPY
 * @Description:
 * @File:  oauth.go
 * @Version: 1.0.0
 * @Date: 2022/2/9 11:16
 */

package wecom

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	OauthHost    = `open.work.weixin.qq.com`
	qrPath       = "/wwopen/sso/qrConnect"
	userInfoPath = "/user/getuserinfo"
)

type AuthRequest struct {
	Appid       string `json:"appid" url:"appid"`                     // 企业微信的CorpID，在企业微信管理端查看
	Agentid     string `json:"agentid" url:"agentid"`                 // 授权方的网页应用ID，在具体的网页应用中查看
	RedirectUri string `json:"redirect_uri" url:"redirect_uri"`       // 重定向地址，需要进行UrlEncode
	State       string `json:"state,omitempty" url:"state,omitempty"` // 用于保持请求和回调的状态，授权请求后原样带回给企业。该参数可用于防止csrf攻击（跨站请求伪造攻击），建议企业带上该参数，可设置为简单的随机数加session进行校验
	Lang        string `json:"lang,omitempty" url:"lang,omitempty"`   // 自定义语言，支持zh、en；lang为空则从Headers读取Accept-Language，默认值为zh
}

// OauthGinLogin 执行重定向到企业微信
func OauthGinLogin(req AuthRequest, c *gin.Context) error {
	u, err := oauthRequestForm(req)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusFound, u.String())
	return nil
}

// OauthGetUserinfo 企业微信认证成功后的回调
// 返回用户ID,错误
func OauthGetUserinfo(token, code string) (string, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   https,
			Host:     apiHost,
			Path:     path.Join(basePath, userInfoPath),
			RawQuery: url.Values{"access_token": []string{token}, "code": []string{code}}.Encode(),
		},
	}
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return "", err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var r struct {
		Errcode int    `json:"errcode"` // 企业微信错误码，正常时为0
		Errmsg  string `json:"errmsg"`  // 企业微信错误信息，正常时为"ok"
		UserId  string `json:"UserId"`
	}
	err = json.Unmarshal(bodyBytes, &r)
	if err != nil {
		return "", err
	}
	if r.Errcode != 0 || r.Errmsg != "ok" || r.UserId == "" {
		if r.Errmsg == "" {
			return "", errors.New("获取用户信息失败")
		}
		return "", errors.New(r.Errmsg)
	}
	return r.UserId, nil
}

func oauthRequestForm(req AuthRequest) (*url.URL, error) {
	if req.Appid == "" {
		return nil, errors.New("appid为空")
	}
	if req.Agentid == "" {
		return nil, errors.New("agentid为空")
	}
	if req.RedirectUri == "" {
		return nil, errors.New("redirect_uri为空")
	}
	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}
	r := url.URL{
		Scheme:   "https",
		Host:     OauthHost,
		Path:     qrPath,
		RawQuery: v.Encode(),
	}
	return &r, nil
}
