package member

import (
	"net/http"
	"store-management-be/baselib/network"
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"store-management-be/application/user_management/model"
	"regexp"
	"time"
)

func AddMember(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	var m model.Member
	err := readData(r.Body, &m)
	if err != nil {
		network.NewFailure(parameterFail).AppendErr(err.Error()).Response(w)
	}

	e := checkMem(&m)
	if e != 0 {
		network.NewFailure(e).Response(w)
	}

	id, err := model.AddMember(m)
	if err != nil {
		network.NewHttpFailure(http.StatusInternalServerError, dbOperaFail).Response(w)
		return
	}
	network.NewSuccess(map[string]int64{"user_id": id}).Response(w)

}

func checkMem(m *model.Member) (memberFailure) {
	if m.Name == "" {
		return NameNilFail
	}
	if len(m.Phone) < 11 {
		return PhoneFail
	}

	if m.MemberLevel <= 0 {
		return MemLevelFail
	}
	if m.Amount < 0 {
		m.Amount = 0
	}
	m.AddTime = int(time.Now().Unix())
	m.Status = 1
	if m.EffectTime < m.AddTime {
		return EffectTimeFail
	}
	return 0
}

func isIdCard(idCard string) bool {
	res, _ := regexp.Match("^[1-9]\\d{7}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}$|^[1-9]\\d{5}[1-9]\\d{3}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}([0-9]|X)$", []byte(idCard))
	return res
}

func MemberInfo(requester *network.Requester, w http.ResponseWriter, r *http.Request) {

	fmt.Println("uid=", requester.Token.Identifier())
	fmt.Println("user=", requester.Token.UserCert())
	fmt.Println("business=", requester.Token.BusinessCert())

	m := make(map[string]interface{})
	m["id"] = 100001
	m["age"] = 23
	m["name"] = "yangjun"
	m["phone"] = "18356628848"
	m["sex"] = "男"

	network.NewSuccess(m).Response(w)
}

func readData(r io.Reader, val interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &val)
	if err != nil {
		return err
	}
	return nil
}
