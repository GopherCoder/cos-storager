package config

import (
	"os"
)

// 七牛云 (Access/Secret Key）
const (
	QNAK     = "tFn2q_EINsikTSpZSHDYERs7fJgvt7muaRLd0iOpX"
	QNSK     = "3Q-rrrnCqLqRkQnEDLvYZOVWzAoXwgXln7-axopDW"
	QNBUCKET = "qiniuyun"
	QNLINK   = "http://pkuescs8i.bkt.clouddn.com"
)

// 腾讯云 ((Access/Secret Key)
const (
	TXAK     = "AKIDYau6DWoAJf2HQRiDzcnoiRp3DICJFCuEX"
	TXSK     = "kMqKUWP8jHhykaqALj0ABRQHzYQAkE8kW"
	TXBUCKET = "tenxun-1253878710"
	TXREGION = "ap-shanghai"
	TXLINK   = "http://tenxun-1253878710.cos.ap-shanghai.myqcloud.com"
)

// 又拍云

const (
	UPYUNBUCKET   = "upyun-wuxiaoshenX"
	UPYUNOPERATOR = "wuxiaoshenX"
	UPYUNPASSOWRD = "upyun123456789W"
	UPYUNLink     = "http://upyun-wuxiaoshen.test.upcdn.net"
)

// SMMS

const (
	SMMSBUCKET = "smms"
	SMMSLINK   = "https://sm.ms/"
)

// 阿里云
const (
	ALIYUNTYPE            = "aliyun"
	ALIYUNBUCKET          = "aliyun-wuxiaoshen"
	ALIYUNENDPOINT        = "http://oss-cn-shanghai.aliyuncs.com"
	ALIYUNACCESSKEYID     = "LTAIWVJGDdDYTIRpX"
	ALIYUNACCESSKEYSECRET = "aVX9Sd3JdcHGPsReGs67Dr1zMORq0IW"
)

// 数据库相关设置

const (
	DATABASENAME     = "cos_storage"
	DATABASEPORT     = "5433"
	DATABASEPASSWORD = ""
	DATABASESSLMODE  = "disable"
	DATABASEUSER     = "postgres"
	DATABASEHOST     = "127.0.0.1"
)

var EnvGlobalParams struct {
	QN struct {
		Bucket string
		AK     string
		SK     string
		Link   string
	}
	TX struct {
		Bucket string
		AK     string
		SK     string
		RG     string
		Link   string
	}
	UP struct {
		Bucket string
		OP     string
		SK     string
		Link   string
	}
	AL struct {
		Type     string
		Bucket   string
		AK       string
		SK       string
		EndPoint string
	}
}

func EnvInit() {
	EnvGlobalParams.QN = struct {
		Bucket string
		AK     string
		SK     string
		Link   string
	}{Bucket: os.Getenv("QNBUCKET"), AK: os.Getenv("QNAK"), SK: os.Getenv("QNSK"), Link: os.Getenv("QNLINK")}
	EnvGlobalParams.TX = struct {
		Bucket string
		AK     string
		SK     string
		RG     string
		Link   string
	}{Bucket: os.Getenv("TXBUCKET"), AK: os.Getenv("TXAK"), SK: os.Getenv("TXSK"), RG: os.Getenv("TXREGION"), Link: os.Getenv("TXLINK")}
	EnvGlobalParams.UP = struct {
		Bucket string
		OP     string
		SK     string
		Link   string
	}{Bucket: os.Getenv("UPYUNBUCKET"), OP: os.Getenv("UPYUNOPERATOR"), SK: os.Getenv("UPYUNPASSOWRD"), Link: os.Getenv("UPYUNLink")}
	EnvGlobalParams.AL = struct {
		Type     string
		Bucket   string
		AK       string
		SK       string
		EndPoint string
	}{Type: os.Getenv("ALIYUNTYPE"), Bucket: os.Getenv("ALIYUNBUCKET"), AK: os.Getenv("ALIYUNACCESSKEYID"), SK: os.Getenv("ALIYUNACCESSKEYSECRET"), EndPoint: os.Getenv("ALIYUNENDPOINT")}
}

// 从环境变量中取值
