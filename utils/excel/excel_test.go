package excel

import (
	"fmt"
	"testing"
)

type UserItem struct {
	Id   int
	name string
	age  int
}

func initServer() *Service {
	s := &Service{
		FileSavePath: "E:\\runtime\\files\\",
		FileName:     "export_test",
		SheetName:    "test",
		Headers:      nil,
		RowLastIndex: 0,
		IsDownload:   false,
	}

	return s
}

// 导出excel，流式写入（通过map数据）
func TestExportExcelByStreamWriteMapData(*testing.T) {
	s := initServer()
	headers := []interface{}{"用户ID", "名称", "年龄", "性别", "爱好"}
	headersMap := map[interface{}]interface{}{
		"age":  "年龄",
		"id":   "用户ID",
		"sex":  "性别",
		"like": "爱好",
		"name": "名称",
	}

	s.Headers = headers
	s.HeadersMap = headersMap
	s.InitStreamWrite()

	rows := []map[interface{}]interface{}{
		{"id": "1", "age": "18", "name": "张三", "sex": "男", "like": "篮球", "unknown": "aa"},
		{"id": "2", "age": "18", "name": "李四", "sex": "女", "like": "羽毛球", "unknown": "bb"},
		{"id": "3", "age": "25", "name": "王五", "sex": "未知", "like": "乒乓球", "unknown": "cc"},
	}
	s.StreamWriteByMapData(rows).StreamClose()
	s.SaveFile()
}

// 导出excel流式写入（通过数据列表）
func TestExportExcelByStreamWriteListData(*testing.T) {
	s := initServer()
	headers := []interface{}{"用户ID", "名称", "年龄", "性别", "爱好"}
	headersMap := map[interface{}]interface{}{
		"age":  "年龄",
		"id":   "用户ID",
		"sex":  "性别",
		"like": "爱好",
		"name": "名称",
	}

	s.Headers = headers
	s.HeadersMap = headersMap
	s.InitStreamWrite()

	rows := [][]interface{}{
		{"1", "test1", "11", "男", "篮球"},
		{"2", "test2", "12", "女", "足球"},
		{"3", "test3", "13", "未知", "羽毛球"},
	}
	s.StreamWriteByListData(rows).StreamClose()
	s.SaveFile()
}

// 导出excel流式写入（通过单行写入）
func TestExportExcelByStreamWriteRowData(*testing.T) {
	s := initServer()
	headers := []interface{}{"用户ID", "名称", "年龄", "性别", "爱好"}
	headersMap := map[interface{}]interface{}{
		"age":  "年龄",
		"id":   "用户ID",
		"sex":  "性别",
		"like": "爱好",
		"name": "名称",
	}

	s.Headers = headers
	s.HeadersMap = headersMap
	s.InitStreamWrite()

	for i := 1; i < 4; i++ {
		rows := []interface{}{i, fmt.Sprintf("test%d", i), 10 + i, fmt.Sprintf("sex%d", i), "篮球"}
		s.StreamWriteByRowData(rows)
	}

	s.StreamClose()
	s.SaveFile()
}
