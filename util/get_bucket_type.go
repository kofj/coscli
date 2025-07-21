package util

import (
	"context"
	logger "github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"strings"
)

// GetBucketType 获取桶类型
func GetBucketType(c *cos.Client, param *Param, config *Config, bucketName string) (string, error) {
	bucketType := BucketTypeCos
	if param.BucketType != "" {
		bucketType = strings.ToUpper(strings.TrimSpace(param.BucketType))
		if bucketType != "COS" && bucketType != "OFS" {
			logger.Fatalln("bucket type can only be either COS or OFS ")
		}
	} else {
		if config.Base.DisableAutoFetchBucketType == "true" {
			bucket, _, err := FindBucket(config, bucketName)
			if err != nil {
				return "", err
			}
			if bucket.Ofs {
				bucketType = BucketTypeOfs
			}
		} else {
			// 判断桶是否是ofs桶
			s, err := c.Bucket.Head(context.Background())
			if err != nil {
				return "", err
			}
			// 根据s.Header判断是否是融合桶或者普通桶
			if s.Header.Get("X-Cos-Bucket-Arch") == BucketTypeOfs {
				bucketType = BucketTypeOfs
			}
		}
	}
	logger.Info("桶类型", bucketType)
	return bucketType, nil
}
