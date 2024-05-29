package emailz

import (
	"crypto/tls"
	"errors"
	"github.com/jordan-wright/email"
	"github.com/zeromicro/go-zero/core/logx"
	"net/smtp"
	"time"
)

type MailContent struct {
	emailFrom       string   // 设置发送方的邮箱
	emailTo         []string // 设置接收方的邮箱
	emailCc         []string // 设置抄送如果抄送多人逗号隔开
	emailBcc        []string // 设置秘密抄送
	emailSubject    string   // 设置主题
	emailText       string   // text邮件内容,text和html只能二选一
	emailHtml       string   // html邮件,text和html只能二选一
	emailAttachFile string   // 附件
}

type MailConfig struct {
	address  string        // 邮箱服务地址,smtp.qq.com:587
	host     string        // 邮箱服务host,smtp.qq.com
	identity string        // 授权码
	username string        // 邮箱账号名称
	password string        // 邮箱发邮件的授权密码
	poolNum  int           // 协程数量
	timeout  time.Duration // 超时时间
	tls      *tls.Config
}

type emailPool struct {
	pool *email.Pool
	conf MailConfig
}

// 获取邮件服务协程池
func NewEmailPool(config MailConfig) (*emailPool, error) {
	ePool, err := email.NewPool(
		config.address,
		config.poolNum,
		smtp.PlainAuth(config.identity, config.username, config.password, config.host),
		config.tls,
	)
	if err != nil {
		logx.Errorf("faild to create email pool:%v", err)
		return nil, err
	}
	return &emailPool{
		pool: ePool,
		conf: config,
	}, nil
}

// 发送单个邮件
func (ep *emailPool) SendEmail(data MailContent) error {
	if data.emailFrom == "" {
		return errors.New("email from cannot empty")
	}
	if len(data.emailTo) == 0 {
		return errors.New("email to cannot empty")
	}
	e := email.NewEmail()
	e.From = data.emailFrom
	e.To = data.emailTo
	e.Cc = data.emailCc
	e.Bcc = data.emailBcc
	e.Subject = data.emailSubject
	e.Text = []byte(data.emailText)
	e.HTML = []byte(data.emailHtml)
	e.AttachFile(data.emailAttachFile)
	err := ep.pool.Send(e, ep.conf.timeout)
	if err != nil {
		logx.Errorf("email:%v sent error:v%\n", e, err)
		return err
	}
	return nil
}
