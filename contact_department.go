// 通讯录 - 部门管理

package wecom

import (
	"encoding/binary"
	"encoding/json"
	"net/url"
	"strconv"
)

// DepartmentCreate 创建部门,返回部门ID或错误
// https://developer.work.weixin.qq.com/document/path/90205
// - 必须携带： name, parentid(根部门id为1)
// - 可选携带id，若不填则自动生成
func (w *Wecom) DepartmentCreate(dpmt Department) (int, error) {
	body, err := w.post("department/create", dpmt)
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint32(body["id"])), nil
}

// DepartmentUpdate 更新部门
// https://developer.work.weixin.qq.com/document/path/90206
// - 必须携带：id
func (w *Wecom) DepartmentUpdate(dpmt Department) error {
	_, err := w.post("department/update", dpmt)
	return err
}

// DepartmentDelete 删除部门
// https://developer.work.weixin.qq.com/document/path/90207
func (w *Wecom) DepartmentDelete(id int) error {
	var query = url.Values{}
	query.Add("id", strconv.Itoa(id))
	_, err := w.get("department/delete", query)
	if err != nil {
		return err
	}
	return nil
}

// DepartmentListGet 获取单个部门详情
// https://developer.work.weixin.qq.com/document/path/95351
func (w *Wecom) DepartmentGet(id int) (*Department, error) {
	var query = url.Values{}
	query.Add("id", strconv.Itoa(id))
	body, err := w.get("department/delete", query)
	if err != nil {
		return nil, err
	}
	var r Department
	err = json.Unmarshal(body["department"], &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// DepartmentListGet 获取部门列表
// https://developer.work.weixin.qq.com/document/path/90208
// id - 部门ID，如果传入，则获取id对应部门及其下属部门列表;否则获取全量组织架构
func (w *Wecom) DepartmentListGet(id ...int) ([]Department, error) {
	var query = url.Values{}
	if len(id) > 0 {
		query.Add("id", strconv.Itoa(id[0]))
	}
	b, err := w.get("department/list", query)
	if err != nil {
		return nil, err
	}
	var r []Department
	err = json.Unmarshal(b["department"], &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DepartmentSubIDListGet 获取子部门ID列表
// https://developer.work.weixin.qq.com/document/path/95350
// id - 部门id。获取指定部门及其下的子部门（以及子部门的子部门等等，递归）。 如果不填，默认获取全量组织架构
// 返回值中只会包含 id、parentid、order字段
func (w *Wecom) DepartmentSubIDListGet(id ...int) ([]Department, error) {
	var query = url.Values{}
	if len(id) > 0 {
		query.Add("id", strconv.Itoa(id[0]))
	}
	b, err := w.get("department/simplelist", query)
	if err != nil {
		return nil, err
	}
	var r []Department
	err = json.Unmarshal(b["department_id"], &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
