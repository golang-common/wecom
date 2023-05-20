// 异步导出接口

// 涉及到的encoding_aeskey说明：Base64编码后的加密密钥。长度固定为43，从a-z, A-Z, 0-9共62个字符中选取，
// 是AESKey的Base64编码。解码后即为32字节长的AESKey。加密方式采用AES-256-CBC方式，数据采用PKCS#7填充至32字节的倍数；
// IV初始向量大小为16字节，取AESKey前16字节，详见：https://datatracker.ietf.org/doc/html/rfc2315

package wecom

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// AsyncExportUser 导出成员
// AESKey=Base64_Decode(encoding_aeskey + “=”)
// aeskey - (必选)用于解密结果的aeskey
// blockSize - 10^4 ~ 10^6之间，默认10^6
func (w *Wecom) AsyncExportUser(aeskey string, blockSize ...int) (string, error) {
	var a = map[string]string{
		"encoding_aeskey": base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(aeskey)),
	}
	if len(blockSize) > 0 {
		a["block_size"] = strconv.Itoa(blockSize[0])
	}
	body, err := w.post("export/simple_user", a)
	if err != nil {
		return "", err
	}
	jobid := strings.Trim(string(body["jobid"]), `"`)
	return jobid, nil
}

// AsyncExportUserDetail 导出成员详细信息
// AESKey=Base64_Decode(encoding_aeskey + “=”)
// aeskey - (必选)用于解密结果的aeskey
// blockSize - 10^4 ~ 10^6之间，默认10^6
func (w *Wecom) AsyncExportUserDetail(aeskey string, blockSize ...int) (string, error) {
	var a = map[string]string{
		"encoding_aeskey": base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(aeskey)),
	}
	if len(blockSize) > 0 {
		a["block_size"] = strconv.Itoa(blockSize[0])
	}
	body, err := w.post("export/user", a)
	if err != nil {
		return "", err
	}
	jobid := strings.Trim(string(body["jobid"]), `"`)
	return jobid, nil
}

// AsyncExportDepartment 导出部门
// AESKey=Base64_Decode(encoding_aeskey + “=”)
// aeskey - (必选)用于解密结果的aeskey
// blockSize - 10^4 ~ 10^6之间，默认10^6
func (w *Wecom) AsyncExportDepartment(aeskey string, blockSize ...int) (string, error) {
	var a = map[string]string{
		"encoding_aeskey": base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(aeskey)),
	}
	if len(blockSize) > 0 {
		a["block_size"] = strconv.Itoa(blockSize[0])
	}
	body, err := w.post("export/department", a)
	if err != nil {
		return "", err
	}
	jobid := strings.Trim(string(body["jobid"]), `"`)
	return jobid, nil
}

// AsyncExportTagMember 导出标签成员
// AESKey=Base64_Decode(encoding_aeskey + “=”)
// aeskey - (必选)用于解密结果的aeskey
// blockSize - 10^4 ~ 10^6之间，默认10^6
// tagid - 标签id
func (w *Wecom) AsyncExportTagMember(tagid, aeskey string, blockSize ...int) (string, error) {
	var a = map[string]string{
		"encoding_aeskey": base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(aeskey)),
		"tagid":           tagid,
	}
	if len(blockSize) > 0 {
		a["block_size"] = strconv.Itoa(blockSize[0])
	}
	body, err := w.post("export/taguser", a)
	if err != nil {
		return "", err
	}
	jobid := strings.Trim(string(body["jobid"]), `"`)
	return jobid, nil
}

// AsyncExportGetResult 根据jobid，获取导出任务结果
// 结果导出后，还需要下载导出文件
func (w *Wecom) AsyncExportGetResult(jobid string) (*ExportResult, error) {
	query := url.Values{}
	query.Add("jobid", jobid)
	body, err := w.get("export/get_result", query)
	if err != nil {
		return nil, err
	}
	var r ExportResult
	err = umarshalObject(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// AsyncExportDownloadResult 从导出结果的url中下载并解密数据
// aeskey - 用户数据解密的key
// 数据解密后的值，可能为几种情况
// 1:用户信息/详细信息 - []User
// 2:部门信息 - []Department
// 3:标签成员信息 - *TagMemberList
func (w *Wecom) AsyncExportDownloadResult(aeskey string, dataUrl ExportUrl) ([]byte, error) {
	resp, err := http.Get(dataUrl.Url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m := md5.New()
	n, err := io.Copy(m, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if int(n) != dataUrl.Size {
		return nil, errors.New("data size mismatch")
	}
	dataMd5 := fmt.Sprintf("%x", m.Sum(nil))
	if dataMd5 != dataUrl.Md5 {
		return nil, errors.New("data md5 mismatch")
	}
	decrypted, err := AesDecryptCBC(body, []byte(aeskey), []byte(aeskey[:16]))
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
