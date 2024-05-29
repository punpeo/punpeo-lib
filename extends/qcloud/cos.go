package qcloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"github.com/zeromicro/go-zero/core/logx"
)

type CosService struct {
	// CI client
	Client     *cos.Client
	Credential *cos.Credential
	suite.Suite

	SecretId   string // 用户的 SecretId，
	SecretKey  string // 用户的 SecretKey，
	BucketName string // 存储桶名称，由 bucketname-appid 组成
	Region     string // ap-guangzhou
	Domain     string // 自定义域名
	AppId      string // 应用ID
}

type PutObjectResponse struct {
	ETag           string // 上传文件的 MD5 值
	XCosExpiration string // 设置生命周期后，返回文件过期规则
}

func (s *CosService) InitClient() *CosService {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", s.BucketName, s.Region)) // 原始
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  s.SecretId,
			SecretKey: s.SecretKey,
		},
	})

	s.Client = c
	s.Credential = s.Client.GetCredential()

	return s
}

/*
*

	    高级接口上传对象

		key := "exampleobject"  // 存储桶的位置
		fileUrl := "../test"  // 文件路径
*/
func (s *CosService) TransferUploadFile(key, fileUrl string) (*cos.CompleteMultipartUploadResult, error) {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(fmt.Sprintf("高级接口上传对象transferUploadFile，err：%+v", err))
		}
	}()

	client := s.Client
	//.cssg-snippet-body-start:[transfer-upload-file]

	resp, _, err := client.Object.Upload(
		context.Background(), key, fileUrl, nil,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil

	//.cssg-snippet-body-end
}

/*
*
获取预签名URL

key 对象键（Key）是对象在存储桶中的唯一标识
expired 有效时间
*/
func (s *CosService) GetPreSignedURL(key string, expired time.Duration) (*url.URL, error) {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(fmt.Sprintf("获取预签名getPreSignedURL，err：%+v", err))
		}
	}()

	// PresignedURLOptions 提供用户添加请求参数和请求头部
	opt := &cos.PresignedURLOptions{
		// http 请求参数，传入的请求参数需与实际请求相同，能够防止用户篡改此 HTTP 请求的参数
		Query: &url.Values{},
		// http 请求头部，传入的请求头部需包含在实际请求中，能够防止用户篡改签入此处的 HTTP 请求头部
		Header: &http.Header{},
	}
	// 添加请求参数, 返回的预签名 url 将包含该参数
	opt.Query.Add("x-cos-security-token", s.Credential.SessionToken)

	// 获取预签名 URL
	client := s.Client
	presignedURL, err := client.Object.GetPresignedURL(context.Background(), http.MethodGet, key, s.SecretId, s.SecretKey, expired, opt)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}

/*
*
简单上传对象（上传文件的内容，可以为文件流或字节流）
若不是 bytes.Buffer/bytes.Reader/strings.Reader 时，必须指定 opt.ObjectPutHeaderOptions.ContentLength
*/
func (s *CosService) PutObject(key string, f io.Reader) (*PutObjectResponse, error) {
	client := s.Client
	resp, err := client.Object.Put(context.Background(), key, f, nil)
	if err != nil {
		return nil, err
	}

	etag := resp.Header.Get("ETag")
	exp := resp.Header.Get("x-cos-expiration")

	//.cssg-snippet-body-end
	return &PutObjectResponse{
		ETag:           etag,
		XCosExpiration: exp,
	}, nil
}

/*
*
获取预签名URL（自定义域名）

key 对象键（Key）是对象在存储桶中的唯一标识
expired 有效时间
*/
func (s *CosService) GetPreSignedURLByDomain(key string, expired time.Duration) (*url.URL, error) {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(fmt.Sprintf("获取预签名URL（自定义域名）getPreSignedURLByDomain，err：%+v", err))
		}
	}()

	// 修改成用户的自定义域名
	u, _ := url.Parse(s.Domain)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{})
	ctx := context.Background()

	// 获取预签名
	presignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodGet, key, s.SecretId, s.SecretKey, expired, nil)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}

/*
*
对象访问 URL
获取对象访问 URL 用于匿名下载或分发。

key 对象键（Key）是对象在存储桶中的唯一标识
*/
func (s *CosService) GetObjectUrl(key string) *url.URL {
	client := s.Client

	ourl := client.Object.GetObjectURL(key)

	return ourl
}

/*
*
下载对象（高级接口）
*/
func (s *CosService) DownloadObject(key, localFile string) (*cos.Response, error) {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(fmt.Sprintf("下载对象（高级接口）downloadObject，err：%+v", err))
		}
	}()

	client := s.Client
	opt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
	}

	resp, err := client.Object.Download(
		context.Background(), key, localFile, opt,
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 获取临时密钥
func (s *CosService) GetTemporaryKey(durationSeconds int64) (*sts.CredentialResult, error) {
	c := sts.NewClient(
		s.SecretId, s.SecretKey, nil,
		// sts.Host("sts.internal.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	opt := &sts.CredentialOptions{
		DurationSeconds: durationSeconds, // 有效期（单位：秒）
		Region:          "ap-guangzhou",
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						// 这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						// 存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						"qcs::cos:" + s.Region + ":uid/" + s.AppId + ":" + s.BucketName + "/*",
					},
					// 开始构建生效条件 condition
					// 关于 condition 的详细设置规则和COS支持的condition类型可以参考https://cloud.tencent.com/document/product/436/71306
					Condition: map[string]map[string]interface{}{
						//"ip_equal": map[string]interface{}{
						//	"qcs:ip": []string{
						//		"10.217.182.3/24",
						//		"111.21.33.72/24",
						//	},
						//},
					},
				},
			},
		},
	}

	// 请求临时密钥
	res, err := c.GetCredential(opt)
	if err != nil {
		return nil, err
	}

	return res, nil

}
