package qiniu

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/client"
	"github.com/qiniu/go-sdk/v7/storage"
	"net/http"
	"os"
	"time"
)

var (
	AK            = os.Getenv("qiniuAccessKey")
	SK            = os.Getenv("qiniuAsecretKey")
	clt           client.Client
	mac           *qbox.Mac
	bucketManager *storage.BucketManager
)

func init() {
	clt = client.Client{
		Client: &http.Client{
			Timeout: time.Minute * 10,
		},
	}
	mac = auth.New(AK, SK)
	cfg := storage.Config{}
	cfg.UseCdnDomains = true
	bucketManager = storage.NewBucketManagerEx(mac, &cfg, &clt)
}

func GetToken(bucket string) string {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)
	return upToken
}
