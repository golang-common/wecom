package wecom

// - MsgType = event
//   - Event = batch_job_result
//     - JobType = export_user(导出成员详情)
//	   - JobType = export_simple_user(导出成员)
//	   - JobType = export_department(导出部门)
//	   - JobType = export_tag(导出标签成员)
//	   - JobType = sync_user(增量更新成员完成通知)
// 	   - JobType = replace_user(全量覆盖成员完成通知)
//     - JobType = invite_user(邀请成员关注完成通知)
//     - JobType = replace_party(全量覆盖部门完成通知)
//   - Event = change_contact
//     - ChangeType = create_user(新增成员事件)
//     - ChangeType = update_user(更新成员事件)
//     - ChangeType = delete_user(删除用户事件)
//     - ChangeType = create_party(新增部门事件)
//     - ChangeType = update_party(更新部门事件)
//     - ChangeType = delete_party(删除部门事件)
//     - ChangeType = update_tag(标签成员变更事件)

// CallbackEvent
type CallbackEvent struct {
	ToUserName   string `xml:"ToUserName"`   // 企业微信CorpID
	FromUserName string `xml:"FromUserName"` // 消息的产生这，企业微信发出的为sys
	CreateTime   int    `xml:"CreateTime"`   // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType"`      // 消息的类型
	Event        string `xml:"Event"`        // 事件的类型
}

// BatchJobResult 导出任务完成通知
// 对应 Callback.Event = "batch_job_result"
// -> BatchJob []BatchJob
type BatchJob struct {
	JobId   string `xml:"JobId"`
	JobType string `xml:"JobType"`
	Error
}
