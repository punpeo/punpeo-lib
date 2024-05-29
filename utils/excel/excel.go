package excel

import (
	"fmt"
	"github.com/punpeo/punpeo-lib/utils"
	"github.com/syyongx/php2go"
	"github.com/xuri/excelize/v2"
	"net/http"
	"net/url"
)

type Service struct {
	FileSavePath string                      // 文件保存路径
	FileName     string                      // 文件名
	FileFullPath string                      // 文件全路径
	SheetName    string                      // 工作簿
	StreamWriter *excelize.StreamWriter      // 流式写入
	File         *excelize.File              // 操作文件
	Headers      []interface{}               // 表头（用于指定数据顺序），如["用户ID", "名称", "年龄"]
	HeadersKey   []interface{}               // 表头key切片，如["id", "name"]
	HeadersMap   map[interface{}]interface{} // 表头Map，如map["value"] = "表头"
	RowLastIndex int                         // 最后一行索引
	IsDownload   bool                        // 是否下载：false-否（保存文件）；true-下载文件
}

// 初始化表格文件
func (s *Service) InitFile() {
	f := excelize.NewFile()
	// 设置sheet名
	_, _ = f.NewSheet(s.SheetName)

	// 删除默认sheet
	_ = f.DeleteSheet("Sheet1")

	s.File = f
}

// 初始化单元格样式（默认）
func (s *Service) InitColStyle() {
	// 修改列宽
	_ = s.StreamWriter.SetColWidth(1, 15, 12)
}

// 初始化流式写入
func (s *Service) InitStreamWrite() *Service {
	s.InitFile()

	// 创建流式写入
	writer, _ := s.File.NewStreamWriter(s.SheetName)
	s.StreamWriter = writer

	// 初始化
	s.InitColStyle()
	s.RowLastIndex = 1

	// 设置表头
	if len(s.Headers) > 0 {
		_ = writer.SetRow("A1", s.Headers, excelize.RowOpts{StyleID: 1})
		s.RowLastIndex = 2 // 存在表头，则从第二行开始

		if len(s.HeadersMap) > 0 {
			headersMapFlip := php2go.ArrayFlip(s.HeadersMap)
			headersKey := make([]interface{}, 0)
			for _, v := range s.Headers {
				if php2go.ArrayKeyExists(v, headersMapFlip) {
					headersKey = append(headersKey, headersMapFlip[v])
				}
			}

			s.HeadersKey = headersKey
		}
	}

	return s
}

// 流式写入（map结构体数据，按表头顺序写）
func (s *Service) StreamWriteByMapData(mapData []map[interface{}]interface{}) *Service {
	for _, v := range mapData {
		row := make([]interface{}, 0)
		if len(s.HeadersKey) > 0 { // 若存在表头，则按表头数据位置写入行，否则无序写入行（单元格无序）
			for _, h := range s.HeadersKey {
				vv, ok := v[h]
				vv = utils.Ternary(ok, vv, "")
				row = append(row, vv)
			}

		} else {
			row = php2go.ArrayValues(v)
		}

		cell, _ := excelize.CoordinatesToCellName(1, s.RowLastIndex) // 索引转单元格坐标
		_ = s.StreamWriter.SetRow(cell, row)
		s.RowLastIndex++
	}

	return s
}

// 流式写入（行数据列表，需自定义顺序）
func (s *Service) StreamWriteByListData(listData [][]interface{}) *Service {
	for _, v := range listData {
		cell, _ := excelize.CoordinatesToCellName(1, s.RowLastIndex) // 索引转单元格坐标
		_ = s.StreamWriter.SetRow(cell, v)
		s.RowLastIndex++
	}

	return s
}

// 流式写入（单行数据，需自定义顺序）
func (s *Service) StreamWriteByRowData(rowData []interface{}) *Service {
	cell, _ := excelize.CoordinatesToCellName(1, s.RowLastIndex) // 索引转单元格坐标
	_ = s.StreamWriter.SetRow(cell, rowData)
	s.RowLastIndex++

	return s
}

// 结束流式写入
func (s *Service) StreamClose() error {
	err := s.StreamWriter.Flush()
	if err != nil {
		return err
	}

	return nil
}

// 保存文件
func (s *Service) SaveFile() (string, error) {
	fileFullPath := fmt.Sprintf("%s.xlsx", s.FileSavePath+s.FileName)
	err := s.File.SaveAs(fileFullPath)
	if err != nil {
		return "", err
	}

	s.FileFullPath = fileFullPath

	return fileFullPath, nil
}

// 下载文件
func (s *Service) DownloadFile(w http.ResponseWriter) {
	disposition := fmt.Sprintf("attachment; filename=%s.xlsx", url.QueryEscape(s.FileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	s.File.WriteTo(w)
	_ = s.File.Close()
}
