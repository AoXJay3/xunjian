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

	// resp.Body返回的不是字符串，也不是[]byte，而是stream，使用ReadAll将响应体变成[]byte，用于后续打印/转Json
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	// fmt.Println(string(body)) // 查看响应原始数据

	var result DeviceResp
	// 增加错误排查，如果存在，则返回错误内容的前300字符
	if err := json.Unmarshal(body, &result); err != nil { // Unmarshal(b []byte, v any)，将字节数据的内容传入第二个参数指针
		//fmt.Println("返回内容前 300 字符")
		previewLen := min(300, len(body))
		preview := string(body[:previewLen])
		return nil, fmt.Errorf("设备 %s 查询AC数据失败 %w, 返回前300行:\n%s", mac, err, preview)

	}
	// 如果返回的数据中，字段长度为空，则判断为未查询到
	if len(result.Result.Rows) == 0 {
		return nil, fmt.Errorf("未查询到设备%s", mac)
	}

	row := result.Result.Rows[0] // eg.{xxxxS_e02efe80c2d1 1 cd-ydy 30-day recording 60056380 142719078}

	// 从 AC 中获取到套餐天数
	dvrDays, err := func() (int, error) {
		if row.ServiceName == "" {
			return 0, nil
		}
		if len(row.ServiceName) < 2 {
			return 0, fmt.Errorf("%s格式不合法", row.ServiceName)
		}
		serviceNameDay := strings.Split(row.ServiceName, "-")[0]
		serviceNameDayInt, err := strconv.Atoi(serviceNameDay)
		if err != nil {
			return 0, fmt.Errorf("ServiceName 套餐天数不是数字: %s", row.ServiceName)
		}

		return serviceNameDayInt, nil
	}()

	row.DvrDays = dvrDays // 将ServiceName中获取的天数给到响应体内

	return &DeviceRespRow{
		DeviceID:     row.DeviceID,
		OnlineStatus: row.OnlineStatus,
		Region:       row.Region,
		ServiceName:  row.ServiceName,
		UID:          row.UID,
		DID:          row.DID,
		DvrDays:      dvrDays,
	}, nil

}
