package qcloud

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func initService() *CosService {
	// 服务初始化
	s := &CosService{
		SecretId:   os.Getenv("CosSecretId"),
		SecretKey:  os.Getenv("CosSecretKey"),
		BucketName: os.Getenv("CosBucketName"),
		Region:     os.Getenv("CosRegion"),
		Domain:     os.Getenv("CosDomain"),
		AppId:      os.Getenv("AppId"),
	}

	s = s.InitClient()

	return s
}

// 高级接口上传对象
func TestTransferUploadFile(t *testing.T) {
	// 服务初始化
	cosService := initService()

	// 上传
	key := "exampleobject.xlsx"     // 存储桶的位置
	fileUrl := "exampleobject.xlsx" // 上传的文件路径

	resp, err := cosService.TransferUploadFile(key, fileUrl)
	if err != nil {
		fmt.Printf("ListFiles() error, %v", err)
	}

	fmt.Println(fmt.Sprintf("响应结果：%+v", resp))
}

// 简单上传对象（上传文件的内容，可以为文件流或字节流）
func TestPutObject(t *testing.T) {
	// 服务初始化
	cosService := initService()
	f, err := os.Open("./test.jpg") // 文件流
	//f := strings.NewReader("")  // 字节流

	// 上传
	key := "xiaomeng.jpg" // 存储桶的位置
	resp, err := cosService.PutObject(key, f)
	if err != nil {
		fmt.Printf("PutObject error, %v", err)
	}

	fmt.Println(fmt.Sprintf("响应结果：%+v", resp))
	fmt.Println(fmt.Sprintf("访问链接：%s", cosService.Domain+key))
}

// 获取对象访问url
func TestGetObjectUrl(t *testing.T) {
	// 服务初始化
	cosService := initService()

	key := "exampleobject.png"
	//key := "/exampleobject.xlsx"
	resp := cosService.GetObjectUrl(key)

	fmt.Println(fmt.Sprintf("响应结果：%+v", resp))
}

// 获取对象访问预签名Url
func TestGetObjectPreSignedURL(t *testing.T) {
	// 服务初始化
	cosService := initService()

	//key := "exampleobject.png"
	key := "exampleobject.xlsx"
	resp, err := cosService.GetPreSignedURL(key, time.Minute*3)
	if err != nil {
		fmt.Printf("TestGetObjectPreSignedURL() error, %v", err)
	}

	fmt.Println(fmt.Sprintf("响应结果：%+v", resp))
}

// 获取对象访问预签名Url（通过自定义域名）
func TestGetPreSignedURLByDomain(t *testing.T) {
	// 服务初始化
	s := &CosService{
		SecretId:   os.Getenv("CosSecretId"),
		SecretKey:  os.Getenv("CosSecretKey"),
		BucketName: os.Getenv("CosBucketName"),
		Region:     os.Getenv("CosRegion"),
		Domain:     os.Getenv("CosDomain"),
	}

	key := "exampleobject.xlsx"
	resp, err := s.GetPreSignedURLByDomain(key, time.Minute*3)
	if err != nil {
		fmt.Printf("TestGetObjectPreSignedURL() error, %v", err)
	}

	fmt.Println(fmt.Sprintf("响应结果：%+v", resp))
}

// 下载对象
func TestDownloadObject(t *testing.T) {
	// 服务初始化
	cosService := initService()

	//key := "exampleobject.png"
	key := "/exampleobject.xlsx"
	localFile := "./exampleobject_download.xlsx"
	resp, err := cosService.DownloadObject(key, localFile)
	if err != nil {
		fmt.Printf("TestDownloadObject() error, %v", err)
	}

	fmt.Println(fmt.Sprintf("响应结果：%+v", resp))
}

// 获取临时密钥
func TestGetTemporaryKey(t *testing.T) {
	// 服务初始化
	cosService := initService()
	resp, err := cosService.GetTemporaryKey(20 * 60)
	if err != nil {
		fmt.Printf("TestGetTemporaryKey() error, %v", err)
	}

	fmt.Println(fmt.Sprintf("响应结果：%+v", resp))
}
