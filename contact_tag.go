// 通讯录 - 标签管理

package wecom

import (
	"encoding/binary"
	"encoding/json"
	"net/url"
	"strconv"
)

// TagCreate 创建标签
// https://developer.work.weixin.qq.com/document/path/90210
// 创建的标签属于该应用，只有该应用的secret才可以增删成员
// 注意，标签总数不能超过3000个
// tagname - 标签名
// tagid - 可选，非负id整数
// 返回：tagid或错误
func (w *Wecom) TagCreate(tagname string, tagid ...int) (int, error) {
	var a = map[string]string{
		"tagname": tagname,
	}
	if len(tagid) > 0 {
		a["tagid"] = strconv.Itoa(tagid[0])
	}
	body, err := w.post("tag/create", a)
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint32(body["tagid"])), nil
}

// TagUpdate 更新标签名字
// https://developer.work.weixin.qq.com/document/path/90211
func (w *Wecom) TagUpdate(tagname string, tagid int) error {
	var a = map[string]string{
		"tagid":   strconv.Itoa(tagid),
		"tagname": tagname,
	}
	_, err := w.post("tag/update", a)
	if err != nil {
		return err
	}
	return nil
}

// TagDelete 删除标签
// https://developer.work.weixin.qq.com/document/path/90212
func (w *Wecom) TagDelete(tagid int) error {
	var query = url.Values{}
	query.Add("tagid", strconv.Itoa(tagid))
	_, err := w.get("tag/delete", query)
	if err != nil {
		return err
	}
	return nil
}

// TagGetUser 根据标签获取成员列表
// https://developer.work.weixin.qq.com/document/path/90213
// 返回标签名、标签中包含的部门ID列表、用户列表、错误
func (w *Wecom) TagGetUser(tagid int) (string, []int, []User, error) {
	var query = url.Values{}
	query.Add("tagid", strconv.Itoa(tagid))
	b, err := w.get("tag/get", query)
	if err != nil {
		return "", nil, nil, err
	}
	tagname := string(b["tagname"])
	partyList := b["partylist"]
	userlistBoby := b["userlist"]
	var plist []int
	var userList []User
	err = json.Unmarshal(partyList, &plist)
	if err != nil {
		return "", nil, nil, err
	}
	err = json.Unmarshal(userlistBoby, &userList)
	if err != nil {
		return "", nil, nil, err
	}
	return tagname, plist, userList, nil
}

// TagAddUsers 添加标签成员
// https://developer.work.weixin.qq.com/document/path/90214
// userlist - (可选)企业成员ID列表，注意：userlist、partylist不能同时为空，单次请求个数不超过1000
// partylist - (可选)企业部门ID列表，注意：userlist、partylist不能同时为空，单次请求个数不超过100
// 返回值：
// invalidlist - 失败的用户id，"usr1|usr2|usr"
// invalidparty - 失败的部门id，[2,4]
// error - 如果非空，则整体失败
func (w *Wecom) TagAddUsers(tagid int, userlist []string, partylist []int) (string, []int, error) {
	var a = map[string]any{
		"tagid": tagid,
	}
	if len(userlist) > 0 {
		a["userlist"] = userlist
	}
	if len(partylist) > 0 {
		a["partylist"] = partylist
	}
	body, err := w.post("tag/addtagusers", a)
	if err != nil {
		return "", nil, err
	}
	invalidList := string(body["invalidlist"])
	invalidpartyBytes := body["invalidparty"]
	var invalidparty []int
	err = json.Unmarshal(invalidpartyBytes, &invalidparty)
	if err != nil {
		return "", nil, err
	}
	return invalidList, invalidparty, nil
}

// TagDelUsers 删除标签成员
// https://developer.work.weixin.qq.com/document/path/90215
// userlist - (可选)企业成员ID列表，注意：userlist、partylist不能同时为空，单次请求个数不超过1000
// partylist - (可选)企业部门ID列表，注意：userlist、partylist不能同时为空，单次请求个数不超过100
// 返回值：
// invalidlist - 失败的用户id，"usr1|usr2|usr"
// invalidparty - 失败的部门id，[2,4]
// error - 如果非空，则整体失败
func (w *Wecom) TagDelUsers(tagid int, userlist []string, partylist []int) (string, []int, error) {
	var a = map[string]any{
		"tagid": tagid,
	}
	if len(userlist) > 0 {
		a["userlist"] = userlist
	}
	if len(partylist) > 0 {
		a["partylist"] = partylist
	}
	body, err := w.post("tag/deltagusers", a)
	if err != nil {
		return "", nil, err
	}
	invalidList := string(body["invalidlist"])
	invalidpartyBytes := body["invalidparty"]
	var invalidparty []int
	err = json.Unmarshal(invalidpartyBytes, &invalidparty)
	if err != nil {
		return "", nil, err
	}
	return invalidList, invalidparty, nil
}

// TagList 获取标签列表
// https://developer.work.weixin.qq.com/document/path/90216
func (w *Wecom) TagList() ([]Tag, error) {
	body, err := w.get("tag/list", nil)
	if err != nil {
		return nil, err
	}
	taglistBytes := body["taglist"]
	var taglist []Tag
	err = json.Unmarshal(taglistBytes, &taglist)
	if err != nil {
		return nil, err
	}
	return taglist, nil
}
