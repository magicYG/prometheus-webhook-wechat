package models

type WechatNotification struct {
	MsgType  string                      `json:"msgtype"`
	Text     *WechatNotificationText     `json:"text,omitempty"`
	Markdown *WechatNotificationMarkdown `json:"markdown,omitempty"`
}

type WechatNotificationText struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type WechatNotificationMarkdown struct {
	Content string `json:"content"`
}

type WechatNotificationResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
