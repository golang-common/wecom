// 通讯录 - 成员管理

package wecom

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// UserCreate 创建成员
// https://developer.work.weixin.qq.com/document/path/90195
// 必填参数：userid name department
// 如果有extattr，则必须传：type name text.value
// 如果有extattr，则必须传：web.url web.title
func (w *Wecom) UserCreate(user User) error {
	_, err := w.post("user/create", user)
	if err != nil {
		return err
	}
	return nil
}

// UserGet 读取成员，根据userid查询成员详细信息
// https://developer.work.weixin.qq.com/document/path/90196
func (w *Wecom) UserGet(userid string) (*User, error) {
	var query = url.Values{}
	query.Add("userid", userid)
	body, err := w.get("user/get", query)
	if err != nil {
		return nil, err
	}
	var r User
	err = umarshalObject(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// UserListGetByDepartment 根据部门id获取成员列表
// https://developer.work.weixin.qq.com/document/path/90200
// 如需获取该部门及其子部门的所有成员，需先获取该部门下的子部门，然后再获取子部门下的部门成员，逐层递归获取。
// 接口返回：userid、name、department、open_userid
func (w *Wecom) UserListGetByDepartment(departmentId int) ([]User, error) {
	var query = url.Values{}
	query.Add("department_id", strconv.Itoa(departmentId))
	body, err := w.get("user/simplelist", query)
	if err != nil {
		return nil, err
	}
	var r []User
	err = umarshalObject(body, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// UserListGetDetailByDepartment 根据部门id获取成员列表的详细信息
// https://developer.work.weixin.qq.com/document/path/90201
// 如需获取该部门及其子部门的所有成员，需先获取该部门下的子部门，然后再获取子部门下的部门成员，逐层递归获取。
func (w *Wecom) UserListGetDetailByDepartment(departmentId int) ([]User, error) {
	var query = url.Values{}
	query.Add("department_id", strconv.Itoa(departmentId))
	body, err := w.get("user/list", query)
	if err != nil {
		return nil, err
	}
	var r []User
	err = umarshalObject(body, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// UserUpdate 更新成员
// https://developer.work.weixin.qq.com/document/path/90197
// user必须携带userid
// 如果携带mobile：若成员已激活企业微信，则需成员自行修改（此情况下该参数被忽略，但不会报错）
// 如果携带email：若是绑定了腾讯企业邮箱的企业微信，则需要在腾讯企业邮箱中修改邮箱（此情况下该参数被忽略，但不会报错）
func (w *Wecom) UserUpdate(user User) error {
	_, err := w.post("user/update", user)
	if err != nil {
		return err
	}
	return nil
}

// UserDelete 删除成员
// https://developer.work.weixin.qq.com/document/path/90198
// userid = "zhangsan"
func (w *Wecom) UserDelete(userid string) error {
	var query = url.Values{}
	query.Add("userid", userid)
	_, err := w.get("user/delete", query)
	if err != nil {
		return err
	}
	return nil
}

// UserBatchDelete 批量删除成员
// https://developer.work.weixin.qq.com/document/path/90199
// useridlist = ["zhangsan", "lisi"]
func (w *Wecom) UserBatchDelete(useridList []string) error {
	var r = map[string][]string{
		"useridlist": useridList,
	}
	_, err := w.post("user/batchdelete", r)
	if err != nil {
		return err
	}
	return nil
}

// UserConvertToOpenID userid转openid
// https://developer.work.weixin.qq.com/document/path/90202
// 该接口使用场景为企业支付，在使用企业红包和向员工付款时，需要自行将企业微信的userid转成openid
func (w *Wecom) UserConvertToOpenID(userid string) (string, error) {
	var r = map[string]string{
		"userid": userid,
	}
	body, err := w.post("user/convert_to_openid", r)
	if err != nil {
		return "", err
	}
	return string(body["openid"]), nil
}

// UserConvertToUserID openid转userid
// https://developer.work.weixin.qq.com/document/path/90202
// 该接口主要应用于使用企业支付之后的结果查询。
// 开发者需要知道某个结果事件的openid对应企业微信内成员的信息时，可以通过调用该接口进行转换查询。
func (w *Wecom) UserConvertToUserID(openid string) (string, error) {
	var r = map[string]string{
		"openid": openid,
	}
	body, err := w.post("user/convert_to_userid", r)
	if err != nil {
		return "", err
	}
	return string(body["userid"]), nil
}

// UserAuthsucc 企业二次验证接口
// 操作步骤 ：
// 1- 在我的企业，安全与保密中开启二次验证，并设置验证的企业URL
// 2- 当成员登录企业微信或关注微信插件（原企业号）进入企业时，会自动跳转到企业的验证页面。
// 在跳转到企业的验证页面时，会带上如下参数：code=CODE。
// 3- 企业收到code后，使用“通讯录同步助手”调用接口"根据code获取成员信息"获取成员的userid"
// 4- 如果成员是首次加入企业，企业获取到userid，并验证了成员信息后，调用如下接口即可让成员成功加入企业
func (w *Wecom) UserAuthsucc(userid string) error {
	var query = url.Values{}
	query.Add("userid", userid)
	_, err := w.get("user/authsucc", query)
	if err != nil {
		return err
	}
	return nil
}

// UserBatchInvite 批量邀请成员
// https://developer.work.weixin.qq.com/document/path/90975
// 企业可通过接口批量邀请成员使用企业微信，邀请后将通过短信或邮件下发通知
func (w *Wecom) UserBatchInvite(useridList, partyList, tagList []string) (*UserInvalidList, error) {
	var a = map[string][]string{}
	if len(useridList) > 0 {
		a["user"] = useridList
	}
	if len(partyList) > 0 {
		a["party"] = partyList
	}
	if len(tagList) > 0 {
		a["tag"] = tagList
	}
	body, err := w.post("batch/invite", a)
	if err != nil {
		return nil, err
	}
	var r UserInvalidList
	err = umarshalObject(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// UserGetJoinQrcode 获取加入企业二维码, 返回二维码链接
// https://developer.work.weixin.qq.com/document/path/91714
// size_type用于控制二维码图片尺寸，1=171*171, 2=399*399, 3=741*741, 4=2052*2052
func (w *Wecom) UserGetJoinQrcode(sizeType ...int) (string, error) {
	var query = url.Values{}
	if len(sizeType) > 0 {
		query.Add("size_type", strconv.Itoa(sizeType[0]))
	}
	body, err := w.get("corp/get_join_qrcode", query)
	if err != nil {
		return "", err
	}
	return string(body["join_qrcode"]), nil
}

// UserGetIDByMobile 通过手机号获取其所对应的userid
// https://developer.work.weixin.qq.com/document/path/95402
// mobile - 用户手机号
func (w *Wecom) UserGetIDByMobile(mobile string) (string, error) {
	var a = map[string]string{
		"mobile": mobile,
	}
	body, err := w.post("user/getuserid", a)
	if err != nil {
		return "", err
	}
	return string(body["userid"]), nil
}

// UserGetIDByEmail 通过email获取用户ID
// https://developer.work.weixin.qq.com/document/path/95895
// email - 用户电子邮箱地址
// email_type - 1=企业邮箱，2=个人邮箱
func (w *Wecom) UserGetIDByEmail(email string, emailType ...int) (string, error) {
	var a = map[string]string{
		"email": email,
	}
	if len(emailType) > 0 {
		a["email_type"] = strconv.Itoa(emailType[0])
	}
	body, err := w.post("user/get_userid_by_email", a)
	if err != nil {
		return "", err
	}
	return string(body["userid"]), nil
}

// UserGetIDList 获取企业成员的userid与对应的部门ID列表
// https://developer.work.weixin.qq.com/document/path/96067
// cursor - 用于分页查询的游标，字符串类型，由上一次调用返回，首次调用不填
// limit - 分页，预期请求的数据量，取值范围 1 ~ 10000
// 返回值 - 新的游标，用户列表(仅包含userid和department)，错误
func (w *Wecom) UserGetIDList(cursor string, limit int) (string, []User, error) {
	var a = map[string]string{
		"cursor": cursor,
		"limit":  strconv.Itoa(limit),
	}
	body, err := w.post("user/list_id", a)
	if err != nil {
		return "", nil, err
	}
	nextCursor := string(body["next_cursor"])
	userListBytes := body["dept_user"]
	var r []User
	err = json.Unmarshal(userListBytes, &r)
	if err != nil {
		return "", nil, err
	}
	return nextCursor, r, nil
}
