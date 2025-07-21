package util

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// CosCopy copies a file from srcClient to destClient using the provided URLs and FileOperations.
// srcClient and destClient are *cos.Client instances.
// srcUrl and destUrl are StorageUrl instances.
// fo is a *FileOperations instance.
func CosCopy(srcClient, destClient *cos.Client, srcUrl, destUrl StorageUrl, fo *FileOperations) error {
	startT := time.Now().UnixNano() / 1000 / 1000

	fo.Monitor.init(fo.CpType)
	chProgressSignal = make(chan chProgressSignalType, 10)
	go progressBar(fo)

	if srcUrl.(*CosUrl).Object != "" && !strings.HasSuffix(srcUrl.(*CosUrl).Object, CosSeparator) {
		// 单对象copy
		index := strings.LastIndex(srcUrl.(*CosUrl).Object, "/")
		prefix := ""
		relativeKey := srcUrl.(*CosUrl).Object
		if index > 0 {
			prefix = srcUrl.(*CosUrl).Object[:index+1]
			relativeKey = srcUrl.(*CosUrl).Object[index+1:]
		}
		// 获取文件信息
		resp, err := GetHead(srcClient, srcUrl.(*CosUrl).Object, fo.Operation.VersionId)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				// 源文件不在cos上
				return fmt.Errorf("Object not found : %v", err)
			}
			return fmt.Errorf("Head object err : %v", err)
		}

		// copy文件
		skip, err, isDir, size, msg := singleCopy(srcClient, destClient, fo, objectInfoType{prefix, relativeKey, resp.ContentLength, resp.Header.Get("Last-Modified")}, srcUrl, destUrl, fo.Operation.VersionId)

		fo.Monitor.updateMonitor(skip, err, isDir, size)
		if err != nil {
			return fmt.Errorf("%s failed: %v", msg, err)
		}

	} else {
		// 多对象copy
		batchCopyFiles(srcClient, destClient, srcUrl, destUrl, fo)
	}

	CloseErrorOutputFile(fo)
	closeProgress()
	fmt.Printf(fo.Monitor.progressBar(true, normalExit))

	endT := time.Now().UnixNano() / 1000 / 1000
	PrintTransferStats(startT, endT, fo)

	return nil
}

func batchCopyFiles(srcClient, destClient *cos.Client, srcUrl, destUrl StorageUrl, fo *FileOperations) {
	chObjects := make(chan objectInfoType, ChannelSize)
	chError := make(chan error, fo.Operation.Routines)
	chLog := make(chan string, fo.Operation.Routines)
	chListError := make(chan error, 1)

	// 启动进程日志处理协程
	var wgLogger sync.WaitGroup
	wgLogger.Add(1)
	go func() {
		defer wgLogger.Done() // 确保在退出时通知等待组
		for processMsg := range chLog {
			writeProcessLog(processMsg, fo)
		}
	}()

	if fo.BucketType == BucketTypeOfs {
		// 扫描ofs对象大小及数量
		go getOfsObjectList(srcClient, srcUrl, nil, nil, fo, true, false)
		// 获取ofs对象列表
		go getOfsObjectList(srcClient, srcUrl, chObjects, chListError, fo, false, true)
	} else {
		// 扫描cos对象大小及数量
		go getCosObjectList(srcClient, srcUrl, nil, nil, fo, true, false)
		// 获取cos对象列表
		go getCosObjectList(srcClient, srcUrl, chObjects, chListError, fo, false, true)
	}

	for i := 0; i < fo.Operation.Routines; i++ {
		go copyFiles(srcClient, destClient, srcUrl, destUrl, fo, chObjects, chError, chLog)
	}

	completed := 0
	for completed <= fo.Operation.Routines {
		select {
		case err := <-chListError:
			if err != nil {
				if fo.Operation.FailOutput {
					writeError(err.Error(), fo)
				}
			}
			completed++
		case err := <-chError:
			if err == nil {
				completed++
			} else {
				if fo.Operation.FailOutput {
					writeError(err.Error(), fo)
				}
			}
		}
	}

	close(chLog)
	wgLogger.Wait()
	logger.Info("logger record completed")
}

func copyFiles(srcClient, destClient *cos.Client, srcUrl, destUrl StorageUrl, fo *FileOperations, chObjects <-chan objectInfoType, chError chan<- error, chLog chan<- string) {
	for object := range chObjects {
		var skip, isDir bool
		var err error
		var size int64
		var msg string
		var processMsg string
		var sleepTime time.Duration
		for retry := 0; retry <= fo.Operation.ErrRetryNum; retry++ {
			startT := time.Now().UnixNano() / 1000 / 1000
			skip, err, isDir, size, msg = singleCopy(srcClient, destClient, fo, object, srcUrl, destUrl)
			endT := time.Now().UnixNano() / 1000 / 1000
			costTime := int(endT - startT)
			if retry == 0 {
				if err == nil {
					processMsg += fmt.Sprintf("[%s] %s successed,cost %dms\n", time.Now().Format("2006-01-02 15:04:05"), msg, costTime)
				} else {
					processMsg += fmt.Sprintf("[%s] %s failed: %v,cost %dms\n", time.Now().Format("2006-01-02 15:04:05"), msg, err, costTime)
				}
			} else {
				if err == nil {
					processMsg += fmt.Sprintf("[%s] retry[%d] with sleep[%v] %s successed,cost %dms\n", time.Now().Format("2006-01-02 15:04:05"), retry, sleepTime.Seconds(), msg, costTime)
				} else {
					processMsg += fmt.Sprintf("[%s] retry[%d] with sleep[%v] %s failed: %v,cost %dms\n", time.Now().Format("2006-01-02 15:04:05"), retry, sleepTime.Seconds(), msg, err, costTime)
				}
			}
			if err == nil {
				break // Copy succeeded, break the loop
			} else {
				if fo.Operation.ErrRetryInterval == 0 {
					// If the retry interval is not specified, retry after a random interval of 1~10 seconds.
					sleepTime = time.Duration(rand.Intn(10)+1) * time.Second
				} else {
					sleepTime = time.Duration(fo.Operation.ErrRetryInterval) * time.Second
				}

				time.Sleep(sleepTime)
			}
		}

		fo.Monitor.updateMonitor(skip, err, isDir, size)
		chLog <- processMsg
		if err != nil {
			chError <- fmt.Errorf("[%s] %s failed: %w\n", time.Now().Format("2006-01-02 15:04:05"), msg, err)
			continue
		}
	}

	chError <- nil
}

func singleCopy(srcClient, destClient *cos.Client, fo *FileOperations, objectInfo objectInfoType, srcUrl, destUrl StorageUrl, VersionId ...string) (skip bool, rErr error, isDir bool, size int64, msg string) {
	skip = false
	rErr = nil
	isDir = false
	size = objectInfo.size
	object := objectInfo.prefix + objectInfo.relativeKey

	destPath := copyPathFixed(objectInfo.relativeKey, destUrl.(*CosUrl).Object)
	msg = fmt.Sprintf("Copy %s to %s", getCosUrl(srcUrl.(*CosUrl).Bucket, object), getCosUrl(destUrl.(*CosUrl).Bucket, destPath))

	var err error
	// 是文件夹则直接创建并退出
	if size == 0 && strings.HasSuffix(object, "/") {
		isDir = true
	}

	// 仅sync命令执行skip
	if fo.Command == CommandSync && !isDir {
		skip, err = skipCopy(srcClient, destClient, object, destPath, fo)
		if err != nil {
			rErr = err
			return
		}
	}

	if skip {
		return
	}

	threadNum := fo.Operation.ThreadNum
	if threadNum == 0 {
		// 若未设置文件分块并发数,需要根据文件大小和分块大小计算默认分块并发数
		threadNum, err = getThreadNumByPartSize(size, fo.Operation.PartSize)
		if err != nil {
			rErr = err
			return
		}
	}

	url, err := GenURL(fo.Config, fo.Param, srcUrl.(*CosUrl).Bucket)

	srcURL := fmt.Sprintf("%s/%s", url.BucketURL.Host, object)

	opt := &cos.MultiCopyOptions{
		OptCopy: &cos.ObjectCopyOptions{
			&cos.ObjectCopyHeaderOptions{
				CacheControl:       fo.Operation.Meta.CacheControl,
				ContentDisposition: fo.Operation.Meta.ContentDisposition,
				ContentEncoding:    fo.Operation.Meta.ContentEncoding,
				ContentType:        fo.Operation.Meta.ContentType,
				Expires:            fo.Operation.Meta.Expires,
				XCosStorageClass:   fo.Operation.StorageClass,
				XCosMetaXXX:        fo.Operation.Meta.XCosMetaXXX,
			},
			nil,
		},
		PartSize:       fo.Operation.PartSize,
		ThreadPoolSize: threadNum,
	}
	if fo.Operation.Meta.CacheControl != "" || fo.Operation.Meta.ContentDisposition != "" || fo.Operation.Meta.ContentEncoding != "" ||
		fo.Operation.Meta.ContentType != "" || fo.Operation.Meta.Expires != "" || fo.Operation.Meta.MetaChange {
	}
	{
		opt.OptCopy.ObjectCopyHeaderOptions.XCosMetadataDirective = "Replaced"
	}

	if fo.BucketType == BucketTypeOfs {
		_, _, err = destClient.Object.MultiCopy(context.Background(), destPath, srcURL, opt)
	} else {
		_, _, err = destClient.Object.MultiCopy(context.Background(), destPath, srcURL, opt, VersionId...)
	}

	if err != nil {
		rErr = err
		return
	}

	if fo.Operation.Move {
		if err == nil {
			_, err = srcClient.Object.Delete(context.Background(), object, nil)
			rErr = err
			return
		}
	}

	return
}
