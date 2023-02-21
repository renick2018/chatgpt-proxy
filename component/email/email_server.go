package email

import (
	"chatgpt-proxy/config"
	"chatgpt-proxy/lib/logger"
	"fmt"
	"github.com/go-gomail/gomail"
	"os"
)

/*
Send
files: [] string, attach files
 */
func Send(to, subject, body string, files ...string) bool {
	var email = config.Global.EmailServer
	var message = gomail.NewMessage()
	message.SetAddressHeader("From", email.Sender, "Device Watcher")
	// 收件人可以有多个，故用此方式
	message.SetHeader("To", to)
	// 主题
	message.SetHeader("Subject", subject)
	// 正文
	message.SetBody("text/html", body)

	for item := range files{
		if _, err := os.Stat(files[item]); err == nil || os.IsExist(err) {
			message.Attach(files[item])
		}
	}

	d := gomail.NewDialer(email.Host, email.Port, email.Sender, email.Password)
	// 发送
	if err := d.DialAndSend(message); err != nil {
		logger.Warning(fmt.Sprintf("send email err: %s", err.Error()))
		return false
	}
	return true
}
