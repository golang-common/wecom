package wecom

import (
	"encoding/json"
)

// User 用户信息
type User struct {
	UserID           string   `json:"userid" xml:"UserID"`                                // 成员UserID。对应管理端的帐号，企业内必须唯一。不区分大小写，长度为1~64个字节
	Name             string   `json:"name" xml:"Name"`                                    // 成员名称
	Department       []int    `json:"department,omitempty" xml:"Department"`              // 成员所属部门id列表,仅返回该应用有查看权限的部门id
	OpenUserid       string   `json:"open_userid,omitempty" xml:"OpenUserid"`             // 全局唯一。对于同一个服务商，不同应用获取到企业内同一个成员的open_userid是相同的，最多64个字节。仅第三方应用可获取
	Order            []int    `json:"order,omitempty" xml:"Order"`                        // 部门内的排序值，默认为0。数量必须和department一致，数值越大排序越前面
	Position         string   `json:"position,omitempty" xml:"Position"`                  // 职务信息
	Mobile           string   `json:"mobile,omitempty" xml:"Mobile"`                      // 手机号码
	Gender           string   `json:"gender,omitempty" xml:"Gender"`                      // 性别。0表示未定义，1表示男性，2表示女性
	Email            string   `json:"email,omitempty" xml:"Email"`                        // 邮箱
	BizMail          string   `json:"biz_mail,omitempty" xml:"BizMail"`                   // 企业邮箱
	IsLeaderInDept   []int    `json:"is_leader_in_dept,omitempty" xml:"IsLeaderInDept"`   // 表示在所在的部门内是否为部门负责人
	DirectLeader     []string `json:"direct_leader,omitempty" xml:"DirectLeader"`         // 直属上级UserID，最多有五个直属上级
	Avatar           string   `json:"avatar,omitempty" xml:"Avatar"`                      // 头像url
	ThumbAvatar      string   `json:"thumb_avatar,omitempty" xml:"ThumbAvatar"`           // 头像缩略图url
	Telephone        string   `json:"telephone,omitempty" xml:"Telephone"`                // 座机
	Alias            string   `json:"alias,omitempty" xml:"Alias"`                        // 别名
	Status           int      `json:"status,omitempty" xml:"Status"`                      // 激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业
	Address          string   `json:"address,omitempty" xml:"Address"`                    // 地址
	HideMobile       int      `json:"hide_mobile,omitempty" xml:"HideMobile"`             // 是否隐藏手机号 0=未隐藏
	EnglishName      string   `json:"english_name,omitempty" xml:"EnglishName"`           // 英文名
	MainDepartment   int      `json:"main_department,omitempty" xml:"MainDepartment"`     // 主部门，仅当应用对主部门有查看权限时返回
	QrCode           string   `json:"qr_code,omitempty" xml:"QrCode"`                     // 员工个人二维码，扫描可添加为外部联系人(注意返回的是一个url，可在浏览器上打开该url以展示二维码)
	ExternalPosition string   `json:"external_position,omitempty" xml:"ExternalPosition"` // 对外职务，如果设置了该值，则以此作为对外展示的职务

	ExternalProfile *UserExternalProfile `json:"external_profile,omitempty"` // 成员对外属性
}

// UserExternalProfile 用户外部属性
type UserExternalProfile struct {
	ExternalCorpName string `json:"external_corp_name,omitempty"`

	WechatChannels []UserWechatChannel `json:"wechat_channels,omitempty"`
	ExternalAttr   []UserExternalAttr  `json:"external_attr,omitempty"`
}

// UserWechatChannel 视频号属性
type UserWechatChannel struct {
	Nickname string `json:"nickname,omitempty"` // 视频号名称
	Status   int    `json:"status,omitempty"`   // 对外展示视频号状态。0表示企业视频号已被确认，可正常使用，1表示企业视频号待确认
}

// UserExternalAttr 扩展属性
type UserExternalAttr struct {
	Type int    `json:"type"` // 属性类型: 0-文本 1-网页 2-小程序
	Name string `json:"name"` // 属性名称： 需要先确保在管理端有创建该属性，否则会忽略
	Text struct {
		Value string `json:"value"` // 文本属性内容，长度限制32个UTF8字符
	} `json:"text"` // 文本类型的属性
}

// UserInvalidList 邀请成员的报错信息
type UserInvalidList struct {
	InvalidUser  []string `json:"invaliduser"`
	InvalidParty []string `json:"invalidparty"`
	InvalidTag   []string `json:"invalidtag"`
}

// Department 部门详细信息
type Department struct {
	Id       int      `json:"id" xml:"Id"`                              // 部门ID
	Name     string   `json:"name" xml:"Name"`                          // 部门名称
	NameEn   string   `json:"name_en" xml:"NameEn"`                     // 部门英文名
	Leaders  []string `json:"department_leader" xml:"DepartmentLeader"` // 部门负责人Userid, ["zhangsan", "lisi"]
	ParentId int      `json:"parentid" xml:"ParentId"`                  // 上级部门id，根部门为1
	Order    int      `json:"order" xml:"Order"`                        // 在父部门中的次序值，order值大的排序靠前。值范围是[0, 2^32)
}

// Tag 通讯录的标签结构
type Tag struct {
	TagId   int    `json:"tagid" xml:"TagId"`
	Tagname string `json:"tagname"`
}

// TagMemberList tag成员信息结构
type TagMemberList struct {
	TagName   string `json:"tagname"`
	UserList  []User `json:"userlist"`
	PartyList []int  `json:"partylist"`
}

// Import 异步导入接口
type Import struct {
	MediaID  string `json:"media_id"`            // 上传的csv文件的media_id, media_id在素材管理接口上传csv后获取
	ToInvite bool   `json:"to_invite,omitempty"` // 是否邀请新建的成员使用企业微信（将通过微信服务通知或短信或邮件下发邀请，每天自动下发一次，最多持续3个工作日），默认值为true。

	Callback *ImportCallback `json:"callback,omitempty"` // 回调信息。如填写该项则任务完成后，通过callback推送事件给企业。具体请参考应用回调模式中的相应选项
}

// ImportCallback 异步导入时配置的web回调
type ImportCallback struct {
	Url            string `json:"url,omitempty"`            // 企业应用接收企业微信推送请求的访问协议和地址，支持http或https协议
	Token          string `json:"token,omitempty"`          // 用于生成签名
	Encodingaeskey string `json:"encodingaeskey,omitempty"` // 用于消息体的加密，是AES密钥的Base64编码
}

// ImportResult 通过任务ID查询导入结果
type ImportResult struct {
	Status     int    `json:"status"`     // 任务状态，整型，1表示任务开始，2表示任务进行中，3表示任务已完成
	Type       string `json:"type"`       // 操作类型，字节串，目前分别有：1. sync_user(增量更新成员) 2. replace_user(全量覆盖成员)3. replace_party(全量覆盖部门)
	Total      int    `json:"total"`      // 任务运行总条数
	Percentage int    `json:"percentage"` // 目前运行百分比，当任务完成时为100

	Result json.RawMessage `json:"result"` // 详细的处理结果,需要上层应用后续解析
}

// ImportUserResult 用户的导入结果
type ImportUserResult struct {
	UserID string `json:"userid"` // 成员userid
	Error
}

// ImportPartyResult 组的导入结果
type ImportPartyResult struct {
	Action  int `json:"action"`  // 操作类型(按位或), 1=新建部门, 2=更改部门名称, 4=移动部门, 8=修改部门排序
	PartyID int `json:"partyid"` // 部门ID
	Error
}

// ExportResult 通过任务ID查询的导出结果
type ExportResult struct {
	Status   int         `json:"status"`    // 任务状态:0-未处理，1-处理中，2-完成，3-异常失败
	DataList []ExportUrl `json:"data_list"` // 数据文件列表
}

// ExportUrl 导出结果资源所指向的url，用户需要根据该url下载导出文件并解密
type ExportUrl struct {
	Url  string `json:"url"`  // 数据下载链接,支持指定Range头部分段下载。有效期2个小时
	Size int    `json:"size"` // 密文数据大小
	Md5  string `json:"md5"`  // 密文数据md5
}
