package qiniu

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"sync"
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	token := GetToken("video")
	fmt.Println("TOKEN : ", token)
	//assert.Equal(t, "a", token, "七牛token获取失败")
}

// 迁移存储
func TestListFiles(t *testing.T) {
	bucket := "video"
	delimiter, marker := "", ""

	loc, _ := time.LoadLocation("Asia/Shanghai")
	startDayUx, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-02-01 00:00:00", loc)
	endDayUx, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-03-01 00:00:00", loc)
	startDaySt := startDayUx.Unix()
	endDaySt := endDayUx.Unix()

	for curDaySt := startDaySt; curDaySt <= endDaySt; curDaySt = time.Unix(curDaySt, 0).AddDate(0, 1, 0).Unix() {
		prefix := "dsh/front/" + time.Unix(curDaySt, 0).Format("200601")
		fmt.Println(prefix)
		wg := sync.WaitGroup{}
		limit, page := 1000, 1
		for {
			entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, delimiter, marker, limit)
			if err != nil {
				fmt.Printf("ListFiles() error, %s", err)
				break
			}
			if len(entries) == 0 {
				fmt.Printf("empty entries")
				break
			}
			marker = nextMarker
			lastKey := ""
			var chtypeKeys = map[string]int{}
			for index, entry := range entries {
				chtypeKeys[entry.Key] = 1
				if index == len(entries)-1 {
					lastKey = entry.Key
				}
			}
			fmt.Println("len |", len(entries))
			chtypeOps := make([]string, 0, len(chtypeKeys))
			for key, fileType := range chtypeKeys {
				chtypeOps = append(chtypeOps, storage.URIChangeType(bucket, key, fileType))
			}
			wg.Add(1)
			go func(chtypeOps []string, limit, page int, lastKey string) {
				bucketManager.Batch(chtypeOps)
				fmt.Printf("七牛存储桶迁移 limit %d| page: %d | last key : %s ", limit, page, lastKey)
				fmt.Println()
				wg.Done()
			}(chtypeOps, limit, page, lastKey)

			page++
			if len(entries) != limit {
				fmt.Printf("ListFiles() failed, unexpected items count, expected: %d, actual: %d", limit, len(entries))
				break
			}
			if !hasNext {
				fmt.Printf("ListFiles() failed, unexpected hasNext")
				break
			}
		}
		wg.Wait()
	}
}
