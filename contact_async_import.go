// 异步导入接口
// 异步导入步骤大致如下：
// 1- 下载导入csv模版
// 2- 上传导入的csv文件（使用素材上传接口），获取media_id
// 3- 携带media_id,调用接口,出发导入任务，获取任务ID
// 4- 通过任务ID获取导入结果

package wecom

import (
	"net/url"
)

// SyncImportUpdateUser 增量更新成员
// https://developer.work.weixin.qq.com/document/path/90980
func (w *Wecom) SyncImportUpdateUser(syncImport Import) (string, error) {
	body, err := w.post("batch/syncuser", syncImport)
	if err != nil {
		return "", err
	}
	return string(body["jobid"]), nil
}

// SyncImportReplaceUser 全量覆盖成员
// https://developer.work.weixin.qq.com/document/path/90981
func (w *Wecom) SyncImportReplaceUser(syncImport Import) (string, error) {
	body, err := w.post("batch/replaceuser", syncImport)
	if err != nil {
		return "", err
	}
	return string(body["jobid"]), nil
}

// SyncImportReplaceParty 全量覆盖部门
// https://developer.work.weixin.qq.com/document/path/90982
// 接口会忽略to_invite字段
func (w *Wecom) SyncImportReplaceParty(syncImport Import) (string, error) {
	syncImport.ToInvite = false
	body, err := w.post("batch/replaceparty", syncImport)
	if err != nil {
		return "", err
	}
	return string(body["jobid"]), nil
}

// SyncImportGetResult 查询提交过的历史任务
// https://developer.work.weixin.qq.com/document/path/90983
// 返回值中的result字段，有两种可能性：[]SyncImportUserResult, []SyncImportPartyResult
func (w *Wecom) SyncImportGetResult(jobid string) (*ImportResult, error) {
	query := url.Values{}
	query.Add("jobid", jobid)
	body, err := w.get("batch/getresult", query)
	if err != nil {
		return nil, err
	}
	var r ImportResult
	err = umarshalObject(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
