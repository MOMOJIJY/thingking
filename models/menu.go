package models

const (
	EventKeyThreeLian  = "V1001_THREE"
	EventKeyChangeBG   = "V2001_IMG_BG"
	EventKeyOtherTools = "V2002_OTHER"
)

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
			Button{Type: "click", Name: "一键三连", Key: EventKeyThreeLian},
			MultiButton{
				Name: "菜单",
				SubButton: []Button{
					{Type: "click", Name: "证件照生成", Key: EventKeyChangeBG},
					{Type: "click", Name: "更多精彩", Key: EventKeyOtherTools},
				},
			},
		},
	}
}
