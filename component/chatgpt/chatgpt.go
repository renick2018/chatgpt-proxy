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

func AddServer(host string) bool {
	if _, ex := ServerMap[host]; ex {
		return false
	}
	ServerMap[host] = &Server{
		Host:    host,
		Status:  true,
		ConvMap: make(map[string]*Conversation),
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
		if v.ConvMap[nickname] != nil {
			cached = v
		}
		if freest == nil || freest.Workload() > v.Workload() {
			freest = v
		}
	}

	if cached == nil || !cached.Status || (cached.Asking > 5 && time.Now().UnixMilli()-cached.ConvMap[nickname].LastAskTime.UnixMilli() > 600000) {
		return freest
	}
	return cached
}

func serverOffline(server *Server) {
	var text = fmt.Sprintf("chatgpt server %s is offline, check it!", server.Host)
	logger.Warning(text)
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
			fmt.Sprintf("ChatGPT server <strong> %s </strong> is offline <br> while <strong> %d </strong> asking <br> online: <strong> %d </strong> offline: <strong> %d </strong> <br> please check it as soon as possible!", server.Host, alive, len(ServerMap) - alive, server.Asking))
	}
}

func RestartServer(host string) int {
	var count = 0
	for _, s := range ServerMap {
		if strings.Contains(s.Host, host) && s.Status == false {
			s.Status = true
			count++
			logger.Info(fmt.Sprintf("%s restart server: %s", host, s.Host))
		}
	}
	return count
}
