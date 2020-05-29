package api

import (
	"net/http"
	"crypto/tls"
	"strings"
	"net/url"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"errors"
	"testing"
)

var token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTA5NDA3OTksImp0aSI6IjEwMDAwMDEiLCJpc3MiOiJ5YW5nanVuIiwibmJmIjoxNTkwNjQ5Mjk5LCJzdWIiOiJBY2Nlc3NUb2tlbiIsInVzZXIiOnsiaWQiOjEwMDAwMDEsInBob25lIjoiMTgzNTY2Mjg4NDgiLCJwYXNzd29yZCI6IjEyMzQ1In0sImJ1c2luZXNzIjp7ImJ1c2luZXNzX2lkIjoxMDEsImJ1c2lfdHlwZSI6InNlcnZpY2UiLCJzdG9yZV9pZCI6MjAwMDF9fQ.MBS1-x0UG7ZApJk6gAu0aW7xKCm-y0-ADkn6CCjEjyw`
var comUrl = "http://127.0.0.1:8808"
var addMember = "/user/add-member"

func TestApi(t *testing.T) {
	url := comUrl + addMember

	//param := map[string]string{
	//	"Authorization": "12234",
	//}
	//
	//body := map[string]interface{}{
	//	"user_ids":      "",
	//	"delegate_code": 12234,
	//}

	bodyJson, err := json.Marshal(Mem)
	if err != nil {
		panic(err)
	}

	err = Fetchdata(http.MethodPost, url, nil, bodyJson, nil)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("成功")
	}
	//fmt.Println(resp)
}

func Fetchdata(method, u string, param map[string]string, body []byte, res interface{}, fileNames ...string) error {
	defer func() {
		if err := recover(); err != nil {
			//con.GetLogger().WithField("opera", "Fetchdata").Error("recover", err)
			panic(err)
		}
	}()

	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	//client := http.Client{}
	req, err := http.NewRequest(method, u, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", token)
	//req.Header.Add("Content-Type", multipart.NewWriter(&bytes.Buffer{}).FormDataContentType())
	values := url.Values{}
	for k, v := range param {
		values.Set(k, v)
	}
	//req.URL.RawQuery = values.Encode()
	req.PostForm = values

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("--------respBody-------")
	if res != nil {
		err = json.Unmarshal(respBody, res)
		if err != nil {
			fmt.Println("解析错误")
			return err
		}
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println(string(respBody))
	} else {
		return errors.New(fmt.Sprintf("服务器响应失败:%s", string(respBody)))
	}
	return nil
}
