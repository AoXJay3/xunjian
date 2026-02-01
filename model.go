package main

import "time"

// 定义需要的数据格式 JSON格式
type DeviceResp struct {
	Result struct {
		Rows []struct {
			// 下面的字段都必须打上标签，如果不大写首字母，字段无法被json.Unmarshel识别
			DeviceID     string `json:"deviceid"`
			OnlineStatus string `json:"onlineStatus"`
			Region       string `json:"region"`
			ServiceName  string `json:"servicename"`
			UID          string `json:"uid"`
			DID          string `json:"did"`
			DvrDays      int
		} `json:"rows"` // 标签
	} `json:"result"`
}

type DeviceRespRow struct {
	DeviceID     string
	OnlineStatus string
	Region       string
	ServiceName  string
	UID          string
	DID          string
	DvrDays      int
	//DID string
}

// +-------+-------+--------------+--------------------+------------+---------------------+---------------------+----------+------------+------------+---------------------+---------------------+
// | did   | uid   | client_id    | device_id          | service_id | start_time          | end_time            | dvr_days | clip_hours | state_code | create_time         | modify_time         |
// +-------+-------+--------------+--------------------+------------+---------------------+---------------------+----------+------------+------------+---------------------+---------------------+
// | 87482 | 11519 | 910fb531-4d0 | xxxxS_2059a0b8cc10 |         17 | 2019-06-16 11:36:01 | 2019-07-16 00:00:08 |        3 |          3 |          1 | 2017-11-02 11:04:00 | 2020-04-26 14:27:11 |
// | 90033 |  9310 | 910fb531-4d0 | xxxxS_2059a074cd3f |         22 | 2014-10-16 14:53:33 | 2015-10-16 00:00:00 |        7 |          7 |          1 | 2017-11-02 11:04:00 | 2018-11-02 10:00:00 |
// +-------+-------+--------------+--------------------+------------+---------------------+---------------------+----------+------------+------------+---------------------+---------------------+

// device_package表数据结构；Go习惯驼峰命名结构体字段 + 下划线 JSON
type DevicePackage struct {
	Did        int64     `json:"did"`
	Uid        int64     `json:"uid"`
	ClientID   string    `json:"client_id"`
	DeviceID   string    `json:"device_id"`
	ServiceID  int       `json:"service_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	DvrDays    int       `json:"dvr_days"`
	ClipHours  int       `json:"clip_hours"`
	StateCode  int       `json:"state_code"`
	CreateTime time.Time `json:"create_time"`
	ModifyTime time.Time `json:"modify_time"`
}

type CoreDeviceRow struct {
	Did       int
	Uid       int
	Device_id string
	End_time  string
	Dvr_days  int
}

// 套餐天数检查
type PackageCheck struct {
	DeviceID string
	ACDays   int
	DBDays   int
}
