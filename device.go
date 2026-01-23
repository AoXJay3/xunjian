// 该项目只针对登虹设备，在执行巡检前，需先确认设备为登虹直连设备
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// 解析 AC 套餐天数
func parseDaysFromServiceName(serviceName string) (int, error) {
	if serviceName == "" {
		return 0, fmt.Errorf("serviceName 为空")
	}
	svcDayStr := strings.Split(serviceName, "-")
	if len(serviceName) < 2 {
		return 0, fmt.Errorf("serviceName 格式不合法")
	}
	svcDayInt, err := strconv.Atoi(svcDayStr[0])
	if err != nil {
		return 0, fmt.Errorf("serviceName 天数不是数字: %s", serviceName)
	}
	return svcDayInt, nil

}

// // 从 AC 获取数据是针对单个设备的，从 core_data 获取数据是批量查的，所以对 AC 获取设备信息的函数中增加核对 core_data 中套餐是否一致
// func compareDvr(drr *DeviceRespRow, devicePackageList []DevicePackage) bool{
// 	// 确认 AC mac数据
// 	ACMac := row.DeviceID
// 	ACDvr := row.ServiceName

// 	// 获取 coreDvr 数据
// 	for _, DevicePackage := range devicePackageList{
// 		if ACMac != DevicePackage.ClientID{
// 			continue
// 		}

// 		coreDvr := DevicePackage.DvrDays // int

// 	}

// }

// 请求获取设备数据，传参包含client实例、mac地址，输出指针和错误
func queryDevice(client *http.Client, mac string) (*DeviceRespRow, error) {
	url := fmt.Sprintf(
		baseURL+
			"/admin/device/getDeviceListV2.do?email=&deviceId=%s&countFlag=1&deviceStatus=&page=1&rows=20&pageNumber=1&pageSize=20&pageIndex=0&orderby=asc",
		mac,
	)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Apifox/1.0.0")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// fmt.Println("resp.Body为：", resp.Body)
	// resp.Body返回的不是字符串，也不是[]byte，而是stream，使用ReadAll将响应体变成[]byte，用于后续打印/转Json
	body, _ := io.ReadAll(resp.Body)

	var result DeviceResp
	// 增加错误排查，如果存在，则返回错误内容的前300字符
	if err := json.Unmarshal(body, &result); err != nil { // Unmarshal(b []byte, v any)，将字节数据的内容传入第二个参数指针
		fmt.Println("返回内容前 300 字符")
		fmt.Println(string(body[:min(300, len(body))]))
		return nil, err
	}
	// 如果返回的数据中，字段长度为空，则判断为未查询到
	if len(result.Result.Rows) == 0 {
		return nil, fmt.Errorf("未查询到设备")
	}

	row := result.Result.Rows[0] // []struct{...}

	// 增加套餐判断
	// compareDvr

	return &DeviceRespRow{
		DeviceID:     row.DeviceID,
		OnlineStatus: row.OnlineStatus,
		Region:       row.Region,
		ServiceName:  row.ServiceName,
		UID:          row.UID,
	}, nil

}
