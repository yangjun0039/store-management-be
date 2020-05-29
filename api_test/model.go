package api

import (
	"time"
)

type Member struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Sex          int    `json:"sex"`
	Email        string `json:"email"`
	Addr         string `json:"addr"`
	Position     string `json:"position"`
	IdCardNum    string `json:"id_card_num"`
	Photo        string `json:"photo"`
	IntroducerId int    `json:"introducer_id"`
	Note         string `json:"note"`
	AddTime      int    `json:"add_time"`

	MemberLevel int    `json:"member_level"`
	Password    string `json:"password"`
	Amount      int    `json:"amount"`
	Point       int    `json:"point"`
	Status      int    `json:"status"`
	EffectTime  int    `json:"effect_time"`
}

var Mem = Member{
	Name:     "yangjun",
	Phone:    "18356628848",
	Sex:      1,
	Email:    "1365651552@qq.com",
	Addr:     "安徽省芜湖市",
	Position: "搬砖的",

	MemberLevel: 1,
	Amount:      50000,
	EffectTime:  int(time.Now().Unix() + 3600*24*365),
}
