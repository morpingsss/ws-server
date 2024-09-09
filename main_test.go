package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
)

const baseURL = "http://127.0.0.1:7800"

type RegisterRequest struct {
	SystemId string `json:"systemId"`
}

type BindToGroupRequest struct {
	ClientId  string `json:"clientId"`
	GroupName string `json:"groupName"`
	UserId    string `json:"userId"`
	Extend    string `json:"extend"`
}

type GetOnlineListRequest struct {
	GroupName string `json:"groupName"`
}

var conn *websocket.Conn

func TestAPIs(t *testing.T) {
	// 测试 /api/register
	registerReq := RegisterRequest{SystemId: "abcd"}
	registerResp, err := sendPostRequest("/api/register", registerReq)
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}
	fmt.Println("注册响应:", registerResp)

	// 测试 /ws WebSocket 连接
	clientId, err := testWebSocketConnection("abcd")
	if err != nil {
		t.Fatalf("WebSocket 连接失败: %v", err)
	}
	fmt.Println("WebSocket 连接成功，ClientId:", clientId)

	// 测试 /api/bind_to_group
	bindToGroupReq := BindToGroupRequest{
		ClientId:  clientId,
		GroupName: "testGroup",
		UserId:    "user123",
		Extend:    "someData",
	}
	bindToGroupResp, err := sendPostRequestWithHeader("/api/bind_to_group", bindToGroupReq, "abcd")
	if err != nil {
		t.Fatalf("绑定到组失败: %v", err)
	}
	fmt.Println("绑定到组响应:", bindToGroupResp)

	// 测试 /api/get_online_list
	getOnlineListReq := GetOnlineListRequest{GroupName: "testGroup"}
	getOnlineListResp, err := sendPostRequestWithHeader("/api/get_online_list", getOnlineListReq, "abcd")
	if err != nil {
		t.Fatalf("获取在线列表失败: %v", err)
	}
	fmt.Println("获取在线列表响应:", getOnlineListResp)

	// 关闭 WebSocket 连接
	if conn != nil {
		conn.Close()
	}
}

func sendPostRequest(endpoint string, data interface{}) (string, error) {
	url := baseURL + endpoint
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON 序列化失败: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("发送 POST 请求失败: %v", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

func sendPostRequestWithHeader(endpoint string, data interface{}, systemId string) (string, error) {
	url := baseURL + endpoint
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON 序列化失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SystemId", systemId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 POST 请求失败: %v", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

func testWebSocketConnection(systemId string) (string, error) {
	url := fmt.Sprintf("ws://127.0.0.1:7800/ws?systemId=%s", systemId)
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return "", fmt.Errorf("WebSocket 连接失败: %v", err)
	}

	// 接收连接成功消息
	_, response, err := conn.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("接收消息失败: %v", err)
	}
	fmt.Println("WebSocket 响应:", string(response))

	// 假设响应中包含 clientId
	var respData map[string]interface{}
	if err := json.Unmarshal(response, &respData); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	data, ok := respData["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("响应中没有 data 字段")
	}
	clientId, ok := data["clientId"].(string)
	if !ok {
		return "", fmt.Errorf("响应中没有 clientId")
	}

	return clientId, nil
}
