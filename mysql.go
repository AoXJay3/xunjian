// 查询设备的core_data.device_package：本地-->龙洲湾ops-->db从库
package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 测试mysql连接用
func testJoinMysql() error {
	dsn := "go_user:Axj123456!@tcp(10.12.29.23:13306)/core_data?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("打开 mysql 失败: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("Ping mysql 失败: %w", err)
	}
	fmt.Println("目的 mysql 连接成功")
	return nil
}

// 创建连接,输出一个sql.DB连接池的
func openDB() (*sql.DB, error) {
	dsn := "go_user:Axj123456!@tcp(10.12.29.23:13306)/core_data?charset=utf8mb4&parseTime=true&loc=Local"
	// sql.Open创建数据库句柄，并没有真实连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// 使用db.Ping()才能测试连通性
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}

// 通过openDB中输出的*sql.DB来查询数据，使用db.Query方法批量查询，查询需输入请求和设备mac
func queryCoreDevice(db *sql.DB, MacList []string) ([]DevicePackage, error) {
	// 处理 Maclist，将数据改为 xxxxS_ 前缀: coreMacList-->[xxxxS_mac1 xxxxS_mac2]
	var coreMacList []string
	for _, mac := range MacList {
		coreMac := "xxxxS_" + mac
		coreMacList = append(coreMacList, coreMac)
	}
	if len(coreMacList) == 0 {
		return []DevicePackage{}, nil
	}

	// 构造占位符和参数
	placeholders := make([]string, len(coreMacList)) // 长度为mac数量的string切片
	args := make([]interface{}, len(coreMacList))    // 长度为mac数量的interface切片
	for i, mac := range coreMacList {
		placeholders[i] = "?"
		args[i] = mac
	}
	// 将占位符给到 sql 语句中
	sqlStr := fmt.Sprintf(`select 
		did,uid,client_id,device_id,service_id,start_time,end_time,dvr_days,clip_hours,state_code,create_time,modify_time
 		from core_data.device_package where device_id in (%s)`,
		strings.Join(placeholders, ",")) // "?,?,..."
	// 使用db.Query对数据库进行查询，传参：sqlStr-->sql查询语句，args的类型为...any，代表可变数量的任意类型参数，返回查询结果Rows和一个error
	rows, err := db.Query(sqlStr, args...) // 传参sql语句，后跟mac1,mac2,...,批量查询
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 定义一个切片用于存储 db 中查询的信息
	var devicePackageList []DevicePackage
	for rows.Next() {
		var dp DevicePackage
		err := rows.Scan(
			&dp.Did,
			&dp.Uid,
			&dp.ClientID,
			&dp.DeviceID,
			&dp.ServiceID,
			&dp.StartTime,
			&dp.EndTime,
			&dp.DvrDays,
			&dp.ClipHours,
			&dp.StateCode,
			&dp.CreateTime,
			&dp.ModifyTime,
		)
		if err != nil {
			return nil, err
		}
		devicePackageList = append(devicePackageList, dp)
	}

	//遍历完成后检查错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devicePackageList, nil

}

// 查询core_data.device_package中的数据，数据内容如下：
// +-------+-------+--------------+--------------------+------------+---------------------+---------------------+----------+------------+------------+---------------------+---------------------+
// | did   | uid   | client_id    | device_id          | service_id | start_time          | end_time            | dvr_days | clip_hours | state_code | create_time         | modify_time         |
// +-------+-------+--------------+--------------------+------------+---------------------+---------------------+----------+------------+------------+---------------------+---------------------+
// | 87482 | 11519 | 910fb531-4d0 | xxxxS_2059a0b8cc10 |         17 | 2019-06-16 11:36:01 | 2019-07-16 00:00:08 |        3 |          3 |          1 | 2017-11-02 11:04:00 | 2020-04-26 14:27:11 |
// | 90033 |  9310 | 910fb531-4d0 | xxxxS_2059a074cd3f |         22 | 2014-10-16 14:53:33 | 2015-10-16 00:00:00 |        7 |          7 |          1 | 2017-11-02 11:04:00 | 2018-11-02 10:00:00 |
// +-------+-------+--------------+--------------------+------------+---------------------+---------------------+----------+------------+------------+---------------------+---------------------+
