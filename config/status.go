package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Copyright 2013-2026 Temple Private Devlopment. All rights reserved.
// license that can be found in the LICENSE file.

type ClientData struct {
	Name    string `json:"name"`    // 客户端名称
	IP      string `json:"ip"`      // 客户端 IP 地址
	Ports   string `json:"ports"`   // 逗号分隔的端口列表，例如 "8080,8081"
	DBs     string `json:"dbs"`     // 数据库连接信息，例如 "mysql:3306"
	Webs    string `json:"webs"`    // Web 服务地址，例如 "http://example.com"
	Message string `json:"message"` //
}

// Message 定义消息结构
type Message struct {
	ClientName string `json:"name,omitempty"` // 用于注册消息
	Type       string `json:"type"`           // 消息类型
	Data       string `json:"data,omitempty"` // 状态数据，可选
	// Stop       func() // 停止函数
}

func ReportStatusJSON(clientData *ClientData) (func(), error) {
	jsonData, err := json.Marshal(clientData)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	msg := Message{
		ClientName: clientData.Name,
		Type:       "status", // 初始消息使用 "status" 类型
		Data:       string(jsonData),
	}

	stop, err := ReportStatus(&msg)
	if err != nil {
		return nil, err
	}

	// 监控 clientData.Message 的变化并发送更新
	go func() {
		lastMessage := clientData.Message
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if clientData.Message != "" && clientData.Message != lastMessage {
					updatedJSON, err := json.Marshal(clientData)
					if err != nil {
						log.Printf("Failed to marshal updated clientData: %v", err)
						continue
					}
					msg.Data = string(updatedJSON)
					lastMessage = clientData.Message
					// 注意：这里无法直接访问 ReportStatus 的 dataChan，但由于 msg.Data 被更新，ReportStatus 的内部逻辑会检测并发送
				}
			}
		}
	}()

	return stop, nil
}

func ReportStatus(msg *Message) (func(), error) {
	var err error

	if msg == nil {
		err = errors.New("消息不能为空")
		return nil, err
	}

	// 创建 clientData 并验证 msg.Data
	clientData := ClientData{
		Name: msg.ClientName,
	}
	if msg.Data != "" {
		if err := json.Unmarshal([]byte(msg.Data), &clientData); err != nil {
			log.Printf("无法解析 msg.Data 到 clientData: %v", err)
			return nil, fmt.Errorf("无效的 msg.Data: %v", err)
		}
	}

	wsURL := "wss://okfarm.hashken.com/okfarm.status"
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 测试用，生产环境移除
		},
	}

	dataChan := make(chan Message, 10)
	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	// 连接和重连逻辑
	go func(cd ClientData) { // 传递 clientData 副本
		defer wg.Done()
		for {
			select {
			case <-stopChan:
				return
			default:
				err := connectAndRun(wsURL, dialer, &cd, dataChan, stopChan)
				if err != nil {
					log.Printf("连接错误: %v", err)
					time.Sleep(5 * time.Second)
					continue
				}
			}
		}
	}(clientData) // 传递副本

	// 检查消息变化并发送
	go func(m Message) { // 使用 msg 的副本
		lastData := m.Data
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		if m.Type != "" && m.Data != "" {
			dataChan <- Message{Type: m.Type, Data: m.Data}
			lastData = m.Data
			m.Data = ""
		}

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				if m.Data != "" && m.Data != lastData {
					dataChan <- Message{Type: m.Type, Data: m.Data}
					lastData = m.Data
					m.Data = ""
				}
			}
		}
	}(*msg) // 传递 msg 副本

	stopFunc := func() {
		close(stopChan)
		wg.Wait()
	}

	return stopFunc, nil
}

func connectAndRun(wsURL string, dialer websocket.Dialer, clientData *ClientData, dataChan <-chan Message, stopChan <-chan struct{}) error {
	// 获取客户端的IP
	clientIP, err := getLocalIP()
	if err != nil {
		log.Printf("Failed to get local IP: %v", err)
		clientIP = "unknown"
	}
	clientData.IP = clientIP

	// 序列化为 JSON
	jsonData, err := json.Marshal(clientData)
	if err != nil {
		return err
	}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 注册客户端，包含IP和端口
	registerMsg := Message{
		Type:       "register",
		ClientName: clientData.Name,
		Data:       string(jsonData), // 使用Data字段发送IP:Port
	}
	if err := sendJSONMessage(conn, registerMsg); err != nil {
		return err
	}

	// 创建心跳定时器
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return nil
		case <-ticker.C:
			// 发送心跳
			heartbeat := Message{Type: "heartbeat"}
			if err := sendJSONMessage(conn, heartbeat); err != nil {
				return err
			}
		case msg := <-dataChan:
			// 发送数据
			if err := sendJSONMessage(conn, msg); err != nil {
				return err
			}
		}
	}
}

// sendJSONMessage 发送 JSON 格式消息
func sendJSONMessage(conn *websocket.Conn, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, jsonData)
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("no IP found")
}
