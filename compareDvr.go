package main

import (
	"fmt"
	"strconv"
	"strings"
)

func findACDvr(respAC []DeviceRespRow, mac string) (int, error) {
	for _, macInfo := range respAC { // macInfo --> {xxxxS_a0ff225ae598 1 fh3-ydy 30-day recording 61890679}
		if mac != macInfo.DeviceID {
			continue
		}
		macSvcName := macInfo.ServiceName // 30-day
		if macSvcName != "" {
			getDaystr := strings.Split(macSvcName, "-")[0] // "30"
			day, err := strconv.Atoi(getDaystr)
			if err != nil {
				return 0, fmt.Errorf("ServiceName=%s 转换失败:%w", macInfo.ServiceName, err)
			}
			fmt.Println("输出天数:", day)
			return day, nil
		}
		return 0, nil

	}
	// 未找到 mac
	return 0, nil
}

// func compareDvr(macList []string) error {

// }
