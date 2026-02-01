package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

// 主要流程：建立持久会话 + 获取JSESSIONID + 密码登录（会话认证） + 调取数据

// 声明变量
const (
	baseURL = "https://cjoint.reservehemu.cn"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请传入设备MAC，支持逗号分隔")
		os.Exit(1)
	}

	macList := strings.Split(os.Args[1], ",") // macList-->[mac1 mac2 mac3]

	// 建立持久会话
	// jar类似python中的session，cookiejar是go标准库中的cookie管理器
	// 功能：自动保存服务器返回的Set-Cookie，下次请求同域名自动带上Cookie
	jar, _ := cookiejar.New(nil)

	// 可复用的HTTP客户端实例，所有请求和响应的Cookie，都交由这个jar来管
	client := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	// 访问goLogin.do获取JSESSIONID,复用client实例
	if err := visitLoginPage(client); err != nil {
		panic(err)
	}

	// 登录
	if err := login(client); err != nil {
		panic(err)
	}

	// 创建CSV文件名称
	now := time.Now()
	nowStr := now.Format("20060102")
	outputCSV := "xunjian" + nowStr + ".csv"

	// 创建CSV文件
	file, err := os.Create(outputCSV)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 先写入一行作为参考
	writer.Write([]string{"device_mac", "online_status", "region", "service_name", "uid"})

	// 查询设备，复用连接，但针对每一个设备单独查询
	var respAC []DeviceRespRow // 接收 AC 返回
	for _, mac := range macList {
		mac := strings.TrimSpace(mac)
		if mac == "" {
			continue
		}
		// 使用queryDevice函数获取设备信息
		resp, err := queryDevice(client, mac)
		if err != nil {
			fmt.Println(err)
			continue
		}

		respAC = append(respAC, *resp)

		// // 如果要输出到控制台
		// fmt.Printf("%s,%s,%s,%s,%s,%s\n",
		// 	resp.DeviceID,
		// 	resp.OnlineStatus,
		// 	resp.ServiceName,
		// 	resp.Region,
		// 	resp.UID,
		// 	resp.DID,
		// )

		// 写入文件
		writer.Write([]string{
			resp.DeviceID,
			resp.OnlineStatus,
			resp.ServiceName,
			resp.Region,
			resp.UID,
			resp.DID,
		})

		time.Sleep(50 * time.Microsecond)

	}
	// fmt.Println("第一个respAC为：", respAC[0])

	fmt.Println("完成，结果已写入:", outputCSV)

	// 前面通过查询AC接口获取到设备相关信息，下面进行查库核实套餐信息
	// 连接数据库
	db, err := openDB()
	if err != nil {
		fmt.Println("打开 mysql 失败,", err)
		return
	}
	defer db.Close()

	result, err := queryCoreDevice(db, macList) // 返回 []DevicePackage{}
	if err != nil {
		fmt.Println("获取core_data数据失败，失败原因:", err)
	}

	// // 控制台展示 DB 数据
	// for _, macInfo := range result {
	// 	fmt.Printf("deviceId:%s,uid:%d,did:%d,dvrDays:%d\n", macInfo.DeviceID, macInfo.Uid, macInfo.Did, macInfo.DvrDays)
	// }

	ACDvrMap := make(map[string]int)
	DBDvrMap := make(map[string]int)

	for _, dp := range respAC {
		ACDvrMap[dp.DeviceID] = dp.DvrDays
	}
	for _, dp := range result {
		DBDvrMap[dp.DeviceID] = dp.DvrDays
	}

	var dvrNotSameDevice []string
	for device_id, _ := range ACDvrMap {
		if ACDvrMap[strings.ToLower(device_id)] != DBDvrMap[strings.ToLower(device_id)] {
			dvrNotSameDevice = append(dvrNotSameDevice, device_id)
			continue
		}
		continue
	}

	fmt.Printf("套餐不一致的设备有%d个:%s\n", len(dvrNotSameDevice), dvrNotSameDevice)

}
