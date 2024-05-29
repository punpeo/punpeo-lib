package emailz

import (
	"crypto/tls"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var efrom = "<xxx@qq.com>"
var eto = []string{"xxx@jianzhikeji.com"}
var ecc = []string{"xxx@jianzhikeji.com"}  //设置抄送如果抄送多人逗号隔开
var ebcc = []string{"xxx@jianzhikeji.com"} //设置秘密抄送
var eSubject = "测试主题"
var eText = "今天是周五 by test mail"
var eHTML = `<h1><a href="http://www.topgoer.com/">go语言中文网站</a></h1>`
var eAttachFile = "./test.txt"

func getMailConfig() MailConfig {
	conf := MailConfig{
		address:  "smtp.qq.com:587",
		host:     "smtp.qq.com",
		identity: "授权码_xxx",
		username: "xxx@qq.com",
		password: "xxxx",
		poolNum:  4,
		timeout:  10 * time.Second,
		tls:      &tls.Config{InsecureSkipVerify: true},
	}
	return conf
}

func Test_email_send(t *testing.T) {

	data := MailContent{
		emailFrom:    efrom,
		emailTo:      eto,
		emailSubject: eSubject,
		//emailCc:
		emailText:       eText,
		emailHtml:       eHTML,
		emailAttachFile: eAttachFile,
	}
	conf := getMailConfig()
	t.Run("Test_email_send", func(t *testing.T) {
		pool, err := NewEmailPool(conf)
		assert.Equal(t, nil, err)
		err = pool.SendEmail(data)
		assert.Equal(t, nil, err)

	})
}

// 腾讯企业邮箱测试
func Test_email_send_exmail(t *testing.T) {
	data := MailContent{
		emailFrom:    "<jianzhizhiku@jianzhikeji.com>",
		emailTo:      eto,
		emailSubject: eSubject,
		//emailCc:
		emailText: eText,
		//emailHtml:       eHTML,
		emailAttachFile: eAttachFile,
	}
	conf := MailConfig{
		address: "smtp.exmail.qq.com:587",
		host:    "smtp.exmail.qq.com",
		//identity: "授权码_7ab5",
		username: "jianzhizhiku@jianzhikeji.com",
		password: "8fDWkuphz2N4ZA6t",
		poolNum:  4,
		timeout:  10 * time.Second,
		tls:      &tls.Config{InsecureSkipVerify: true},
	}

	t.Run("Test_email_send_exmail", func(t *testing.T) {
		pool, err := NewEmailPool(conf)
		assert.Equal(t, nil, err)
		err = pool.SendEmail(data)
		assert.Equal(t, nil, err)

	})
}

func Test_email_send_nofrom(t *testing.T) {
	data := MailContent{
		//emailFrom:    efrom,
		emailTo:      eto,
		emailSubject: eSubject,
		//emailCc:
		emailText:       eText,
		emailHtml:       eHTML,
		emailAttachFile: eAttachFile,
	}
	conf := getMailConfig()
	t.Run("Test_email_send", func(t *testing.T) {
		pool, err := NewEmailPool(conf)
		assert.Equal(t, nil, err)
		err = pool.SendEmail(data)
		assert.Equal(t, errors.New("email from cannot empty"), err)
	})
}

func Test_email_send_noto(t *testing.T) {
	data := MailContent{
		emailFrom: efrom,
		//emailTo:      eto,
		emailSubject: eSubject,
		//emailCc:
		emailText:       eText,
		emailHtml:       eHTML,
		emailAttachFile: eAttachFile,
	}
	conf := getMailConfig()
	t.Run("Test_email_send", func(t *testing.T) {
		pool, err := NewEmailPool(conf)
		assert.Equal(t, nil, err)
		err = pool.SendEmail(data)
		assert.Equal(t, errors.New("email to cannot empty"), err)
	})
}

func Test_email_send_cc(t *testing.T) {

	data := MailContent{
		emailFrom:       efrom,
		emailTo:         eto,
		emailSubject:    eSubject,
		emailCc:         ecc,
		emailBcc:        ebcc,
		emailText:       eText,
		emailHtml:       eHTML,
		emailAttachFile: eAttachFile,
	}
	conf := getMailConfig()
	t.Run("Test_email_send_cc", func(t *testing.T) {
		pool, err := NewEmailPool(conf)
		assert.Equal(t, nil, err)
		err = pool.SendEmail(data)
		assert.Equal(t, nil, err)

	})
}
