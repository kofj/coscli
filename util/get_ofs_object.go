package util

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
	"strings"
)

func getOfsObjectList(c *cos.Client, cosUrl StorageUrl, chObjects chan<- objectInfoType, chError chan<- error, fo *FileOperations, scanSizeNum bool, withFinishSignal bool) {
	if chObjects != nil {
		defer close(chObjects)
	}

	prefix := cosUrl.(*CosUrl).Object
	marker := ""
	limit := 0

	delimiter := ""
	if fo.Operation.OnlyCurrentDir {
		delimiter = "/"
	}

	err := getOfsObjectListRecursion(c, cosUrl, chObjects, chError, fo, scanSizeNum, prefix, marker, limit, delimiter)

	if err != nil && scanSizeNum {
		fo.Monitor.setScanError(err)
	}

	if scanSizeNum {
		fo.Monitor.setScanEnd()
		freshProgress()
	}

	if withFinishSignal {
		chError <- err
	}
}

func getOfsObjectListRecursion(c *cos.Client, cosUrl StorageUrl, chObjects chan<- objectInfoType, chError chan<- error, fo *FileOperations, scanSizeNum bool, prefix string, marker string, limit int, delimiter string) error {
	isTruncated := true
	for isTruncated {
		// 实例化请求参数
		opt := &cos.BucketGetOptions{
			Prefix:       prefix,
			Delimiter:    delimiter,
			EncodingType: "url",
			Marker:       marker,
			MaxKeys:      limit,
		}
		res, err := tryGetObjects(c, opt)
		if err != nil {
			return err
		}

		for _, object := range res.Contents {
			object.Key, _ = url.QueryUnescape(object.Key)
			if cosObjectMatchPatterns(object.Key, fo.Operation.Filters) {
				if scanSizeNum {
					fo.Monitor.updateScanSizeNum(object.Size, 1)
				} else {
					objPrefix := ""
					objKey := object.Key
					index := strings.LastIndex(cosUrl.(*CosUrl).Object, "/")
					if index > 0 {
						objPrefix = object.Key[:index+1]
						objKey = object.Key[index+1:]
					}
					chObjects <- objectInfoType{prefix: objPrefix, relativeKey: objKey, size: object.Size, lastModified: object.LastModified}
				}
			}
		}

		if len(res.CommonPrefixes) > 0 {
			for _, commonPrefix := range res.CommonPrefixes {
				commonPrefix, _ = url.QueryUnescape(commonPrefix)

				if cosObjectMatchPatterns(commonPrefix, fo.Operation.Filters) {
					if scanSizeNum {
						fo.Monitor.updateScanSizeNum(0, 1)
					} else {
						objPrefix := ""
						objKey := commonPrefix
						index := strings.LastIndex(cosUrl.(*CosUrl).Object, "/")
						if index > 0 {
							objPrefix = commonPrefix[:index+1]
							objKey = commonPrefix[index+1:]
						}
						chObjects <- objectInfoType{prefix: objPrefix, relativeKey: objKey, size: int64(0), lastModified: ""}
					}
				}

				if delimiter == "" {
					err = getOfsObjectListRecursion(c, cosUrl, chObjects, chError, fo, scanSizeNum, commonPrefix, marker, limit, delimiter)
					if err != nil {
						return err
					}
				}
			}
		}

		isTruncated = res.IsTruncated
		marker, _ = url.QueryUnescape(res.NextMarker)
	}
	return nil
}
func getOfsObjectListForLs(c *cos.Client, prefix string, marker string, limit int, recursive bool) (err error, objects []cos.Object, commonPrefixes []string, isTruncated bool, nextMarker string) {

	delimiter := ""
	if !recursive {
		delimiter = "/"
	}

	// 实例化请求参数
	opt := &cos.BucketGetOptions{
		Prefix:       prefix,
		Delimiter:    delimiter,
		EncodingType: "url",
		Marker:       marker,
		MaxKeys:      limit,
	}
	res, err := tryGetObjects(c, opt)
	if err != nil {
		return
	}

	objects = res.Contents
	commonPrefixes = res.CommonPrefixes
	isTruncated = res.IsTruncated
	nextMarker, _ = url.QueryUnescape(res.NextMarker)
	return
}

// GetOfsKeys reads the keys from the COS URL and returns an error if any occurs.
// c: *cos.Client, cosUrl: StorageUrl, keys: map[string]string, fo: *FileOperations
func GetOfsKeys(c *cos.Client, cosUrl StorageUrl, keys map[string]commonInfoType, fo *FileOperations, objType string) error {

	chFiles := make(chan objectInfoType, ChannelSize)
	chFinish := make(chan error, 2)
	go ReadCosKeys(keys, cosUrl, chFiles, chFinish, fo, objType)
	go getOfsObjectList(c, cosUrl, chFiles, chFinish, fo, false, false)
	select {
	case err := <-chFinish:
		if err != nil {
			return err
		}
	}

	return nil
}
