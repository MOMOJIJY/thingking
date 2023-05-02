package models

type TemplateMsg struct {
	ToUser     string      `json:"touser"`
	TemplateId string      `json:"template_id"`
	URL        string      `json:"url"`
	TopColor   string      `json:"topcolor"`
	Data       interface{} `json:"data"`
}

func NewTemplateMsg(toUser string) TemplateMsg {
	type data struct {
		Value string
		Color string
	}
	return TemplateMsg{
		ToUser:     toUser,
		TemplateId: "6X_ZMmuz4dWwPnhpN2DGiDIZALSX2e96qfbqe1JXvyY", // 先写死
		URL:        "http://www.baidu.com",                        // 先写死
		TopColor:   "#FF0000",                                     // 先写死
		Data: map[string]data{
			"toUser": {Value: "先生/小姐", Color: "#123137"}, // 先写死
		},
	}
}
