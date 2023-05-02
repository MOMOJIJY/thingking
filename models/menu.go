package models

type Button struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
	Key  string `json:"key,omitempty"`
}

type Menu struct {
	Button []interface{} `json:"button"`
}

type MultiButton struct {
	Name      string   `json:"name"`
	SubButton []Button `json:"sub_button"`
}

func NewMenu() Menu {
	return Menu{
		Button: []interface{}{
			Button{Type: "click", Name: "转发", Key: "V1001_FOLLOW"},
			MultiButton{
				Name: "菜单",
				SubButton: []Button{
					{Type: "click", Name: "投币", Key: "V2001_POINT"},
					{Type: "click", Name: "点赞", Key: "V2002_LIKE"},
					{Type: "click", Name: "收藏", Key: "V2003_COLLECT"},
				},
			},
		},
	}
}
