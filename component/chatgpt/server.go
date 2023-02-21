package chatgpt

import (
	"bytes"
	"chatgpt-proxy/lib/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	Status       bool
	Host         string
	ConvMap      map[string]*Conversation // [nickname]
	Asking       int
	count        int
	successCount int
	askLock      sync.Mutex
	countLock    sync.Mutex
}

func (s *Server) Workload() float32 {
	var activeConv = 0
	for _, v := range s.ConvMap{
		if time.Now().UnixMilli() - v.LastAskTime.UnixMilli() < 120000 {
			activeConv++
		}
	}

	return float32(s.Asking) + 0.5 * float32(activeConv) + 1 - 1.0/float32((s.count + s.successCount)/2 + 1)
}

func (s *Server) updateCount(plus bool) {
	s.countLock.Lock()
	defer s.countLock.Unlock()
	if plus {
		s.Asking++
		s.count++
	} else {
		s.Asking--
	}
}

func (s *Server) NewConv(nickname string) {
	if nickname == "" {
		return
	}
	if _, ex := s.ConvMap[nickname]; !ex {
		s.ConvMap[nickname] = &Conversation{
			Nickname: nickname,
			Server:   s,
		}
	}
	s.ConvMap[nickname].LastAskTime = time.Now()
}

func (s *Server) Ask(convId, message string) (*string, string) {
	s.updateCount(true)
	s.askLock.Lock()
	defer func() {
		s.updateCount(false)
		s.askLock.Unlock()
	}()

	if !s.Status {
		return nil, convId
	}
	logger.Info(fmt.Sprintf("%s %s try ask %s", s.Host, convId, message))

	// post for rsp
	var rsp = s.post(convId, message)

	logger.Info(fmt.Sprintf("%s %s try ask %s\nresponse: %+v", s.Host, convId, message, rsp))

	// if ok return rsp
	if rsp != nil {
		s.successCount++
		var convAddr = convId
		if convAddr == "" {
			convAddr = rsp.ConversationID
			s.NewConv(convAddr)
		}
		s.ConvMap[convAddr].ID = rsp.ConversationID
		s.ConvMap[convAddr].LastMessageID = rsp.MessageID
		s.ConvMap[convAddr].LastAskTime = time.Now()
		return &rsp.Message, convAddr
	}

	s.Status = false

	// return nil, false
	return nil, ""
}

func (s *Server) post(convId, message string) *Response {

	url := s.Host + "/ask" // POST 请求的目标 URL
	data := make(map[string]interface{})
	data["message"] = message
	if s.ConvMap[convId] != nil {
		data["messageId"] = s.ConvMap[convId].LastMessageID
		data["conversationId"] = s.ConvMap[convId].ID
	} else {
		data["messageId"] = ""
		data["conversationId"] = ""
	}
	bs, _ := json.Marshal(data) // POST 请求的数据

	logger.Info(fmt.Sprintf("%s request: %s", s.Host, string(bs)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		logger.Warning(fmt.Sprintf("Error creating HTTP request: %+v", err))
		return nil
	}

	req.Header.Set("Content-Type", "application/json") // 设置请求头

	client := &http.Client{}
	resp, err := client.Do(req) // 发送请求
	if err != nil {
		logger.Warning(fmt.Sprintf("Error sending HTTP request: %+v", err))
		return nil
	}

	defer resp.Body.Close()

	fmt.Println()
	logger.Info(fmt.Sprintf("HTTP Response Status: %+v", resp.Status))

	// 读取响应体
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	logger.Info(fmt.Sprintf("HTTP Response Body:: %+v", buf.String()))

	var rsp = make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &rsp)

	if rsp["response"] == nil {
		return nil
	}

	var response = rsp["response"].(map[string]interface{})

	return &Response{
		MessageID:      response["id"].(string),
		Message:        response["response"].(string),
		ConversationID: response["conversationId"].(string),
	}
}

func (s *Server) Info() map[string]interface{} {
	return map[string]interface{}{
		"host":               s.Host,
		"status":             s.Status,
		"conversation_count": len(s.ConvMap),
		"asking":             s.Asking,
		"count":              s.count,
		"success_count":      s.successCount,
		"workload":           s.Workload(),
	}
}
