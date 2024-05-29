package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/punpeo/punpeo-lib/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

var (
	callerDepth = 4
)

// DatabaseGormWriter 在默认 logx.Writer 的基础上，实现将 error 日志记录到数据表中
type DatabaseGormWriter struct {
	logx.Writer
	conn        *gorm.DB
	onceMigrate sync.Once
	tableName   string
}

func NewDatabaseGormWriter(writer logx.Writer, conn *gorm.DB) *DatabaseGormWriter {
	return &DatabaseGormWriter{
		Writer:    writer,
		conn:      conn,
		tableName: "err_log",
	}
}

func (d *DatabaseGormWriter) SetTableName(name string) *DatabaseGormWriter {
	d.tableName = name
	return d
}

func (d *DatabaseGormWriter) Error(v interface{}, fields ...logx.LogField) {
	defer func() {
		if pac := recover(); pac != nil {
			log.Println("ZLog-jz:Pac: ", pac)
		}
	}()

	d.Writer.Error(v, fields...)
	// 创建数据表
	d.onceMigrate.Do(func() {
		d.conn.
			Table(d.tableName).AutoMigrate(&LogErr{})
	})

	content := getContent(v)
	var span, trace = "", ""
	for _, field := range fields {
		switch field.Key {
		case spanKey:
			span, _ = field.Value.(string)
		case traceKey:
			trace, _ = field.Value.(string)
		}
	}

	// record
	log := LogErr{
		Caller:     utils.GetCaller(callerDepth),
		Content:    sql.NullString{String: content, Valid: true},
		Level:      levelError,
		Span:       span,
		Trace:      trace,
		CreateTime: time.Now(),
	}
	d.conn.Table(d.tableName).Create(&log)
}

// getContent 将用户日志信息转化为mysql string 类型数据
func getContent(v interface{}) string {
	var content string
	switch v.(type) {
	case string:
		content = v.(string)
	default:
		contentByte, err := json.Marshal(v)
		if err != nil {
			content = fmt.Sprintf("Marshal err: %s", err.Error())
		} else {
			content = string(contentByte)
		}
	}
	return content
}
