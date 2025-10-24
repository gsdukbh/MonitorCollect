package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

type TelegrafJson struct {
	//telegraf 默认字段
	Tags      map[string]string      `json:"tags"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp int64                  `json:"timestamp"`
	Name      string                 `json:"name"`
}

// parseJson 函数用于解析 JSON 格式的数据
// Telegraf 发送的是数组格式: [{...}, {...}]
func parseJson(body []byte) {
	var metrics []TelegrafJson
	if err := json.Unmarshal(body, &metrics); err != nil {
		log.Printf("解析 JSON 出错: %v", err)
		return
	}

	fmt.Println("--- 收到 JSON 格式数据 ---")
	for _, metric := range metrics {
		fmt.Printf("Measurement: %s\n", metric.Name)
		// 打印 Tags
		for key, val := range metric.Tags {
			fmt.Printf("  Tag: %s = %s\n", key, val)
		}
		switch metric.Name {
		case "cpu":
			var cpu CPUFields
			if err := cpu.FromFieldsMap(metric.Fields); err != nil {
				log.Printf("解析 CPU 字段出错: %v", err)
				continue
			}
			fmt.Printf("CPU %s: 使用率=%.2f%%, 空闲=%.2f%%\n",
				metric.Tags["cpu"], cpu.UsageActive, cpu.UsageIdle)
			//todo直接在cpu.go 文件中处理逻辑。 // 其他类型数据同理，可以根据 metric.Name 进行不同的处理
		}
	}
}

// parseLineProtocol 函数用于解析 InfluxDB Line Protocol 格式的数据
func parseLineProtocol(body []byte) {
	// 使用官方的 line-protocol 解析器
	decoder := lineprotocol.NewDecoder(bytes.NewReader(body))

	fmt.Println("--- 收到 InfluxDB Line Protocol 格式数据 ---")
	for decoder.Next() {
		measurement, err := decoder.Measurement()
		if err != nil {
			log.Printf("解析 Measurement 出错: %v", err)
			continue
		}

		fmt.Printf("Measurement: %s\n", string(measurement))

		// 打印 Tags
		for {
			key, val, err := decoder.NextTag()
			if err != nil {
				break
			}
			fmt.Printf("  Tag: %s = %s\n", string(key), string(val))
		}

		// 打印 Fields
		for {
			key, val, err := decoder.NextField()
			if err != nil {
				break
			}
			fmt.Printf("  Field: %s = %v (%T)\n", string(key), val, val)
		}

		// 获取时间戳
		ts, err := decoder.Time(lineprotocol.Nanosecond, time.Time{})
		if err == nil {
			fmt.Printf("  Timestamp: %s\n", ts)
		}
		fmt.Println("-------------------------------------------")
	}

	if err := decoder.Err(); err != nil {
		log.Printf("解析 Line Protocol 出错: %v", err)
	}
}

// handleJsonMetrics 专门处理 JSON 格式的 Telegraf 数据
func handleJsonMetrics(w http.ResponseWriter, r *http.Request) {
	// 1. 确保是 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "只接受 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	// 2. 处理 Gzip 压缩
	var reader io.Reader = r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "无法解压 Gzip 数据", http.StatusBadRequest)
			return
		}
		defer gzReader.Close()
		reader = gzReader
	}

	// 3. 读取请求体
	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, "无法读取请求体", http.StatusBadRequest)
		return
	}

	// 4. 解析 JSON 格式
	parseJson(body)

	// 5. 返回成功响应
	w.WriteHeader(http.StatusNoContent)
}

// handleLineProtocolMetrics 专门处理 Line Protocol 格式的 Telegraf 数据
func handleLineProtocolMetrics(w http.ResponseWriter, r *http.Request) {
	// 1. 确保是 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "只接受 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	// 2. 处理 Gzip 压缩
	var reader io.Reader = r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "无法解压 Gzip 数据", http.StatusBadRequest)
			return
		}
		defer gzReader.Close()
		reader = gzReader
	}

	// 3. 读取请求体
	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, "无法读取请求体", http.StatusBadRequest)
		return
	}

	// 4. 解析 Line Protocol 格式
	parseLineProtocol(body)

	// 5. 返回成功响应
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// 注册两个不同的端点
	http.HandleFunc("/metrics/json", handleJsonMetrics)
	http.HandleFunc("/metrics/lineprotocol", handleLineProtocolMetrics)

	port := "8080"
	log.Printf("服务器启动，监听在端口 %s, 等待 Telegraf 数据...", port)
	log.Printf("JSON 格式请配置 url 为: http://localhost:%s/metrics/json", port)
	log.Printf("Line Protocol 格式请配置 url 为: http://localhost:%s/metrics/lineprotocol", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("启动服务器失败: %s\n", err)
	}
}
