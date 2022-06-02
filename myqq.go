package myqq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elancom/go-util/str"
	"io"
	"net/http"
	"net/url"
)

type MsgType int

const (
	MsgTypeFriend       MsgType = 1 // 好友
	MsgTypeGroup                = 2 // 群
	MsgTypeGroupTalk            = 3 // 讨论组
	MsgTypeGroupTmp             = 4 // 群临时会话
	MsgTypeGroupTalkTmp         = 5 // 讨论组临时会话
	MsgTypeOnlineTmp            = 6 // 在线临时会话
)

type ApiCreator func() *ApiReq

var (
	ApiGetQQList       = NewApiProc("Api_GetQQList")       // 取框架QQ号
	ApiGetOnlineQQList = NewApiProc("Api_GetOnlineQQlist") // 取框架在线QQ号
	ApiSendMsg         = NewApiProc("Api_SendMsg")         // 发送消息
	ApiSearchGroup     = NewApiProc("Api_SearchGroup")     // 发送消息
)

func NewApiProc(name string) ApiCreator {
	return func() *ApiReq { return &ApiReq{name: name} }
}

type MQ struct {
	Robot     string `json:"MQ_Robot"`     // 用于判定哪个QQ接收到该消息
	Type      int    `json:"MQ_type"`      // 接收到消息类型，该类型可在[常量列表]中查询具体定义
	TypeSub   int    `json:"MQ_type_sub"`  // 此参数在不同情况下，有不同的定义
	FromID    string `json:"MQ_fromID"`    // 此消息的来源，如：群号、讨论组ID、临时会话QQ、好友QQ等
	FromQQ    string `json:"MQ_fromQQ"`    // 主动发送这条消息的QQ，踢人时为踢人管理员QQ
	PassiveQQ string `json:"MQ_passiveQQ"` // 被动触发的QQ，如某人被踢出群，则此参数为被踢出人QQ
	Msg       string `json:"MQ_msg"`       // （此参数将被URL UTF8编码，您收到后需要解码处理）此参数有多重含义，常见为：对方发送的消息内容，但当消息类型为 某人申请入群，则为入群申请理由,当消息类型为收到财付通转账、收到群聊红包、收到私聊红包时为原始json，详情见[特殊消息]章节
	MsgSeq    string `json:"MQ_msgSeq"`    // 撤回别人或者机器人自己发的消息时需要用到
	MsgID     string `json:"MQ_msgID"`     // 撤回别人或者机器人自己发的消息时需要用到
	MsgData   string `json:"MQ_msgData"`   // UDP收到的原始信息，特殊情况下会返回JSON结构（入群事件时，这里为该事件data）
	Timestamp string `json:"MQ_timestamp"` // 对方发送该消息的时间戳，引用回复消息时需要用到
}

func (m *MQ) BodyUnescape() string {
	unescape, err := url.QueryUnescape(m.Msg)
	if err != nil {
		return ""
	}
	return unescape
}

func NewQQGSearch() *QQGSearch {
	return &QQGSearch{}
}

type QQGSearch struct {
	Code         string `json:"code"`         // 群号
	Gid          string `json:"gid"`          // 群号
	Class        string `json:"class"`        // 类别
	ClassId      int    `json:"classId"`      // 类别ID
	Tags         string `json:"tags"`         // 分类标签
	Features     string `json:"features"`     // 状态特征
	Labels       string `json:"labels"`       // 群标签
	MaxMemberNum int    `json:"maxMemberNum"` // 最大会员数
	MemberNum    int    `json:"memberNum"`    // 会员数
	Memo         string `json:"memo"`         // 群描述
	Name         string `json:"name"`         // 群名称
	Image        string `json:"image"`        // 群头像
	CityId       int    `json:"cityId"`       // 城市ID
	Province     string `json:"province"`     // 省份
	City         string `json:"City"`         // 城市
	Level        int    `json:"level"`        // 等级
	Activity     int    `json:"activity"`     // 活动情况?
}

type ApiResp struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}

func (r *ApiResp) Get(key string) any {
	return r.Data.(map[string]any)[key]
}

func (r *ApiResp) GetString(key string) string {
	return r.Data.(map[string]any)[key].(string)
}

func (r *ApiResp) GetMap(key string) map[string]any {
	return r.Data.(map[string]any)[key].(map[string]any)
}

type ApiReq struct {
	url  string
	name string
	ps   map[string]any
}

func (a *ApiReq) Param(name string, v any) *ApiReq {
	if a.ps == nil {
		a.ps = make(map[string]any, 5)
	}
	a.ps[name] = v
	return a
}

func (a *ApiReq) Exec(url string, p map[string]any) (*ApiResp, error) {
	body := make(map[string]any, 3)
	body["function"] = a.name
	body["token"] = "666"
	params := a.ps
	if params == nil {
		params = make(map[string]any, 0)
	}
	if p != nil {
		for k, v := range p {
			params[k] = v
		}
	}
	body["params"] = params
	bodyByte, err := json.Marshal(body)
	fmt.Println(string(bodyByte))
	if err != nil {
		return nil, err
	}
	post, err := http.Post(url, "application/json", bytes.NewReader(bodyByte))
	if err != nil {
		return nil, err
	}
	readAll, err := io.ReadAll(post.Body)
	if err != nil {
		return nil, err
	}
	apiResp := ApiResp{}
	err = json.Unmarshal(readAll, &apiResp)
	if err != nil {
		return nil, err
	}
	return &apiResp, err
}

func NewApi(url string) *Api {
	return &Api{url: url}
}

type Api struct {
	url string
}

// GetQQList 取框架QQ号
func (a *Api) GetQQList() ([]string, error) {
	resp, err := ApiGetQQList().Exec(a.url, nil)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	return str.Split(resp.GetString("ret"), "\r\n"), nil
}

// GetOnlineQQList 取框架在线QQ号
func (a *Api) GetOnlineQQList() ([]string, error) {
	resp, err := ApiGetOnlineQQList().Exec(a.url, nil)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	return str.Split(resp.GetString("ret"), "\r\n"), nil
}

// SendMsg 发送消息
func (a *Api) SendMsg(
	qq string, // 响应QQ
	msgType MsgType, // 信息类型
	toQQG string, // 收信群_讨论组
	toQQ string, // 收信对象
	content string, // 信息内容
) error {
	p := map[string]any{
		"c1": qq,
		"c2": msgType,
		"c3": toQQG,
		"c4": toQQ,
		"c5": content,
	}
	exec, err := ApiSendMsg().Exec(a.url, p)
	if err != nil {
		return err
	}
	if !exec.Success {
		return errors.New(exec.Msg)
	}
	return nil
}

// SearchGroup 搜索群
func (a *Api) SearchGroup(
	qq string, // QQ
	keyword string, // 关键词
	page int, // 页码
) ([]*QQGSearch, error) {
	p := map[string]any{
		"c1": qq,
		"c2": keyword,
		"c3": page,
	}
	exec, err := ApiSearchGroup().Exec(a.url, p)
	if err != nil {
		return nil, err
	}
	getMap := exec.GetMap("ret")
	qqgList := getMap["group_list"].([]any)
	qss := make([]*QQGSearch, 0, len(qqgList))
	for _, q := range qqgList {
		qm := q.(map[string]any)
		// 类别标签
		tags := ""
		if qm["gcate"] != nil {
			for i, cate := range qm["gcate"].([]any) {
				if i > 0 {
					tags += ","
				}
				tags += cate.(string)
			}
		}
		// 群特征
		features := ""
		if qm["group_label"] != nil {
			for i, gl := range qm["group_label"].([]any) {
				if i > 0 {
					features += ","
				}
				features += gl.(map[string]any)["item"].(string)
			}
		}
		// 群标签
		labels := ""
		if qm["labels"] != nil {
			for i, gl := range qm["labels"].([]any) {
				if i > 0 {
					labels += ","
				}
				labels += gl.(map[string]any)["label"].(string)
			}
		}
		// 省份城市
		qaddr := qm["qaddr"].([]any)
		province, city := "", ""
		if len(qaddr) >= 1 {
			province = qaddr[0].(string)
		}
		if len(qaddr) >= 2 {
			city = qaddr[1].(string)
		}
		// 城市ID
		cityId := 0
		if qm["cityid"] != nil {
			cityId = int(qm["cityid"].(float64))
		}
		qs := QQGSearch{
			Code:         str.String(int(qm["code"].(float64))),
			Gid:          str.String(int(qm["gid"].(float64))),
			Class:        qm["class_text"].(string),
			ClassId:      int(qm["class"].(float64)),
			Tags:         tags,
			Features:     features,
			Labels:       labels,
			MaxMemberNum: int(qm["max_member_num"].(float64)),
			MemberNum:    int(qm["member_num"].(float64)),
			Memo:         qm["memo"].(string),
			Name:         qm["name"].(string),
			Image:        qm["url"].(string),
			CityId:       cityId,
			Province:     province,
			City:         city,
			Level:        int(qm["level"].(float64)),
			Activity:     int(qm["activity"].(float64)),
		}
		qss = append(qss, &qs)
	}
	return qss, nil
}
