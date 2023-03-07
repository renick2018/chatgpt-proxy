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
	Status          bool
	Host            string
	ConvMap         map[string]*Conversation // [nickname]
	Asking          int
	count           int
	successCount    int
	askLock         sync.Mutex
	countLock       sync.Mutex
	Code            int
	OffTimestamp    time.Time
	AskingTimestamp time.Time
	Email           string
	Password        string
	ApiKey          string
	IsPlus          bool
	IsAPi           bool
}

func (s *Server) Workload() float32 {
	var activeConv = 0
	for _, v := range s.ConvMap {
		if time.Now().UnixMilli()-v.LastAskTime.UnixMilli() < 120000 {
			activeConv++
		}
	}

	return float32(s.Asking) + 0.5*float32(activeConv) + 1 - 1.0/float32((s.count+s.successCount)/2+1)
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

func (s *Server) IsVIP() bool {
	return s.IsAPi || s.IsPlus
}

func (s *Server) Ask(convId, message string) (*string, string) {
	s.updateCount(true)
	if !s.IsAPi {
		s.askLock.Lock()
	}
	defer func() {
		s.updateCount(false)
		if !s.IsAPi {
			s.askLock.Unlock()
		}
	}()

	if !s.Status {
		return nil, convId
	}
	s.AskingTimestamp = time.Now()
	logger.Info(fmt.Sprintf("%s %s try ask %d letter: %s", s.Host, convId, len([]rune(message)), message))

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

	s.OffTimestamp = time.Now()
	if !s.IsAPi {
		s.Status = false
	}
	// return nil, false
	return nil, ""
}

func (s *Server) post(convId, message string) *Response {

	url := s.Host + "/ask" // POST 请求的目标 URL
	data := make(map[string]interface{})
	data["message"] = message
	data["node"] = s.Email
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

	client := &http.Client{Timeout: 6 * time.Minute}
	resp, err := client.Do(req) // 发送请求
	if err != nil {
		logger.Warning(fmt.Sprintf("Error sending HTTP request: %+v", err))
		return nil
	}

	defer resp.Body.Close()

	logger.Info(fmt.Sprintf("HTTP Response Status: %+v", resp.Status))

	// 读取响应体
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	logger.Info(fmt.Sprintf("HTTP Response Body:: %+v", buf.String()))

	var rsp = make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &rsp)

	var res = rsp["response"]
	if res == nil {
		res = rsp["data"]
	}

	if res == nil {
		var msg = rsp["message"]
		var text = "阿巴阿巴"
		var code = -1
		switch msg.(type) {
		case map[string]interface{}:
			code = int(msg.(map[string]interface{})["statusCode"].(float64))
		case string:
			code = int(rsp["code"].(float64))
		}

		switch code {
		case 413:
			text = fmt.Sprintf("太长了，简短的问哈，脑阔已经打包了")
		case 500:
			text = fmt.Sprintf("这个问题太难了，换一个吧，脑瓜子嗡嗡嗡的")
		default:
			s.Code = code
			return nil
		}

		return &Response{
			Message:        text,
			ConversationID: convId,
		}
	}

	var response = res.(map[string]interface{})

	if len(response) == 0 {
		s.Code = -1
		return nil
	}

	return &Response{
		MessageID:      response["messageId"].(string),
		Message:        response["response"].(string),
		ConversationID: response["conversationId"].(string),
	}
}

func (s *Server) Info() map[string]interface{} {
	return map[string]interface{}{
		"host":               s.Host,
		"email":              s.Email,
		"status":             s.Status,
		"conversation_count": len(s.ConvMap),
		"asking":             s.Asking,
		"count":              s.count,
		"success_count":      s.successCount,
		"workload":           s.Workload(),
		"code":               s.Code,
		"off_timestamp":      s.OffTimestamp.Format("2006-01-02 15:04:05"),
		"asking_timestamp":   s.AskingTimestamp.Format("2006-01-02 15:04:05"),
		"plus":               s.IsPlus,
		"api":                s.IsAPi,
	}
}
