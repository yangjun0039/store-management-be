package model

import (
	"store-management-be/database/mysql"
)

var masterName = "master" //
var slaveName = "slave"

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

func AddMember(m Member) (int64, error) {
	tx, err := mysql.DbMap[masterName].Begin()
	if err != nil {
		return 0, err
	}
	sqlMemInfo := `
		insert into member_info (name, phone, sex, email, addr, position, id_card_num, photo, introducer_id, note, add_time)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	sqlMemData := `
		insert into member_data (user_id, member_level, password, amount, point, status, effect_time)
		values (?, ?, ?, ?, ?, ?, ?);
	`
	id, err := mysql.DbMap[masterName].InsertTx(tx, sqlMemInfo, m.Name, m.Phone, m.Sex, m.Email, m.Addr, m.Position, m.IdCardNum, m.Photo, m.IntroducerId, m.Note, m.AddTime)
	if err != nil {
		return id, err
	}
	_,err = mysql.DbMap[masterName].InsertTx(tx, sqlMemData, id, m.MemberLevel, m.Password, m.Amount, m.Point, m.Status, m.EffectTime)
	if err != nil {
		return id, err
	}
	return id,nil
}
