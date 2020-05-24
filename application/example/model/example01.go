package model

import(
	"store-management-be/database/mysql"
)

var masterName = "master"
var slaveName = "slave"

func QryMasterData() (map[string]string, error){
	sqlStr := `
		select * from table_test;
	`
	data, err := mysql.DbMap[masterName].FetchRowD(sqlStr)
	return data,err
}

func QrySlaveData() (map[string]string, error){
	sqlStr := `
		select * from user;
	`
	data, err := mysql.DbMap[slaveName].FetchRowD(sqlStr)
	return data,err
}