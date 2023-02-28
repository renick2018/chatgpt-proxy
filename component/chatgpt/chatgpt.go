package chatgpt

import (
	"chatgpt-proxy/component/email"
	"chatgpt-proxy/config"
	"chatgpt-proxy/lib/logger"
	"fmt"
	"strings"
	"sync"
	"time"
)

var locker sync.Mutex
var ServerMap = make(map[string]*Server)
var alertTimestamp int64 = 0

func LoadServers() {
	for _, host := range config.Global.ChatServerAddrs {
		ServerMap[host] = &Server{
			Host:    host,
			Status:  true,
			ConvMap: make(map[string]*Conversation),
		}
		logger.Info("load chatgpt server: ", host)
	}
}

func LoadServersV2() {
	for _, item := range config.Global.GPTServers {
		ServerMap[item.Email] = &Server{
			Host:     item.Host,
			Email:    item.Email,
			Password: item.Password,
			Status:   true,
			ConvMap:  make(map[string]*Conversation),
		}
		var data = make(map[string]interface{})
		var nodes = make([]map[string]interface{}, 1)
		nodes[0] = make(map[string]interface{})
		nodes[0]["email"] = item.Email
		nodes[0]["password"] = item.Password
		data["nodes"] = nodes
		rsp, err := callServer(fmt.Sprintf("%s/add_nodes", item.Host), data)
		logger.Info(fmt.Sprintf("load chatgpt server err: %+v\nrsp: %v", err, rsp))
	}
}

func AddServer(host, email, password string) bool {
	if _, ex := ServerMap[host]; ex {
		return false
	}
	var server = &Server{
		Host:    host,
		Status:  true,
		ConvMap: make(map[string]*Conversation),
	}
	ServerMap[email] = server
	if len(email) != 0 {
		var data = make(map[string]interface{})
		var nodes = make([]map[string]interface{}, 1)
		nodes[0] = make(map[string]interface{})
		nodes[0]["email"] = email
		nodes[0]["password"] = password
		data["nodes"] = nodes

		callServer(fmt.Sprintf("%s/restart_nodes", host), data)
	}
	logger.Info("add chatgpt server: ", host)
	return true
}

// Ask return rsp, nickname
func Ask(nickname, message string) (*string, string) {
	// get a freest ai server
	var server = fetchSever(nickname)

	// if no available return nil
	if server == nil {
		logger.Info(fmt.Sprintf("%s no available server", nickname))
		return nil, nickname
	}
	logger.Info(fmt.Sprintf("%s fetch server %v", nickname, server.Host))

	// ask
	var rsp, conv = server.Ask(nickname, message)

	// if ok return rsp
	if rsp != nil {
		return rsp, conv
	}

	// warning
	serverOffline(server)

	// Ask()
	return Ask(nickname, message)
}

func fetchSever(nickname string) (s *Server) {
	locker.Lock()
	defer func() {
		if s != nil {
			s.NewConv(nickname)
		}
		locker.Unlock()
	}()

	var freest *Server
	var cached *Server
	for _, v := range ServerMap {
		if !v.Status {
			continue
		}
		if v.ConvMap[nickname] != nil && (cached == nil || v.Workload() < cached.Workload()){
			cached = v
		}
		if freest == nil || freest.Workload() > v.Workload() {
			freest = v
		}
	}

	if cached != nil && (cached.Asking < 5 || time.Now().UnixMilli()-cached.ConvMap[nickname].LastAskTime.UnixMilli() < 300000){
		return cached
	}

	return freest
}

func serverOffline(server *Server) {
	var text = fmt.Sprintf("chatgpt server %s is offline, check it! %+v", server.Host, time.UnixMilli(alertTimestamp).Format("2006-01-02 15:04:05"))
	logger.Warning(text)
	if time.Now().UnixMilli() - alertTimestamp < 60000 {
		return
	}
	alertTimestamp = time.Now().UnixMilli()
	var alive = 0
	for _, v := range ServerMap{
		if v.Status {
			alive++
		}
	}
	// 企业微信/飞书/邮箱
	for _, to := range config.Global.Emails{
		email.Send(to,
			fmt.Sprintf("ChatGPT %s offline", strings.SplitAfter(server.Host, ":")[2]),
			fmt.Sprintf("ChatGPT server <strong> %s </strong> is offline <br> while <strong> %d </strong> asking <br> online: <strong> %d </strong> offline: <strong> %d </strong> <br> please check it as soon as possible!", server.Host, server.Asking, alive, len(ServerMap) - alive))
	}
}

func RestartServer(host string) int {
	var count = 0
	for _, s := range ServerMap {
		if (strings.Contains(s.Host, host) || strings.Contains(s.Email, host)) && s.Status == false {
			s.Status = true
			count++
			logger.Info(fmt.Sprintf("%s restart server: %s", host, s.Host))
			if len(s.Email) != 0 {
				var data = make(map[string]interface{})
				var nodes = make([]map[string]interface{}, 1)
				nodes[0] = make(map[string]interface{})
				nodes[0]["email"] = s.Email
				data["nodes"] = nodes

				callServer(fmt.Sprintf("%s/restart_nodes", s.Host), data)
			}
		}
	}
	return count
}
