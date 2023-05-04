package models

import "encoding/xml"

type Type string

const (
	MsgTypeText  Type = "text"
	MsgTypeImage Type = "image"
	MsgTypeVideo Type = "video"

	XmlContentType = "application/xml; charset=utf-8"
)

// RequestRawMessage 解析用户发来的消息
type RequestRawMessage struct {
	XMLName xml.Name `xml:"xml"`
	// 基础字段
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"` // user's openID
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      Type   `xml:"MsgType"`
	MsgID        int64  `xml:"MsgId"`
	MsgDataId    string `xml:"MsgDataId"`
	// 文本字段
	Content string `xml:"Content"`
	// 图片、视频、语言共有字段
	MediaID string `xml:"MediaId"`
	// 图片字段
	PicURL string `xml:"PicUrl"`
	// 语言字段
	Format string `xml:"Format"`
	// 视频字段
	ThumbMediaID string `xml:"ThumbMediaId"`
	// 事件字段
	Event    string `xml:"Event"`
	EventKey string `xml:"EventKey"`
}

// ResponseBaseMessage 提供响应信息的基础字段
type ResponseBaseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      Type     `xml:"MsgType"`
}

// ResponseTextMessage 文本的响应信息
type ResponseTextMessage struct {
	ResponseBaseMessage
	Content string `xml:"Content"`
}

// ResponseImageMessage 图片的响应信息
type ResponseImageMessage struct {
	ResponseBaseMessage
	MediaID string `xml:"Image>MediaId"`
}

// ResponseChangeBG 替换背景外部接口的响应
type ResponseChangeBG struct {
	Errors []ChangeBGError
}

type ChangeBGError struct {
	Title string
	Code  string
}
