package util

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mu sync.Mutex
)

// Counter 传输大小
type Counter struct {
	TransferSize int64
}

// Upload 上传文件
func Upload(c *cos.Client, fileUrl StorageUrl, cosUrl StorageUrl, fo *FileOperations) {
	startT := time.Now().UnixNano() / 1000 / 1000
	localPath := fileUrl.ToString()

	fo.Monitor.init(fo.CpType)
	chProgressSignal = make(chan chProgressSignalType, 10)
	go progressBar(fo)

	chFiles := make(chan fileInfoType, ChannelSize)
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

	// 统计文件数量及大小数据
	go fileStatistic(localPath, fo)
	// 生成文件列表
	go generateFileList(localPath, chFiles, chListError, fo)

	for i := 0; i < fo.Operation.Routines; i++ {
		go uploadFiles(c, cosUrl, fo, chFiles, chError, chLog)
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

	closeProgress()
	fmt.Printf(fo.Monitor.progressBar(true, normalExit))

	endT := time.Now().UnixNano() / 1000 / 1000
	PrintTransferStats(startT, endT, fo)
}

func uploadFiles(c *cos.Client, cosUrl StorageUrl, fo *FileOperations, chFiles <-chan fileInfoType, chError chan<- error, chLog chan<- string) {
	for file := range chFiles {
		var skip, isDir bool
		var err error
		var size, transferSize int64
		var msg string
		var processMsg string
		var sleepTime time.Duration
		for retry := 0; retry <= fo.Operation.ErrRetryNum; retry++ {
			startT := time.Now().UnixNano() / 1000 / 1000
			skip, err, isDir, size, transferSize, msg = SingleUpload(c, fo, file, cosUrl)
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
				break // Upload succeeded, break the loop
			} else {

				if fo.Operation.ErrRetryInterval == 0 {
					// If the retry interval is not specified, retry after a random interval of 1~10 seconds.
					sleepTime = time.Duration(rand.Intn(10)+1) * time.Second
				} else {
					sleepTime = time.Duration(fo.Operation.ErrRetryInterval) * time.Second
				}

				time.Sleep(sleepTime)

				fo.Monitor.updateDealSize(-transferSize)
			}
		}

		fo.Monitor.updateMonitor(skip, err, isDir, size)
		chLog <- processMsg
		if err != nil {
			// 获取当前时间
			chError <- fmt.Errorf("[%s] %s failed: %w\n", time.Now().Format("2006-01-02 15:04:05"), msg, err)
			continue
		}
	}

	chError <- nil
}

// SingleUpload 单文件上传
func SingleUpload(c *cos.Client, fo *FileOperations, file fileInfoType, cosUrl StorageUrl) (skip bool, rErr error, isDir bool, size, transferSize int64, msg string) {
	skip = false
	rErr = nil
	isDir = false
	size = 0
	transferSize = 0

	localFilePath, cosPath := UploadPathFixed(file, cosUrl.(*CosUrl).Object)

	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		rErr = err
		return
	}

	var snapshotKey string

	msg = fmt.Sprintf("Upload %s to %s", localFilePath, getCosUrl(cosUrl.(*CosUrl).Bucket, cosPath))
	if fileInfo.IsDir() {
		isDir = true
		if fo.Operation.SkipDir {
			skip = true
		} else {
			// 在cos创建文件夹
			_, err = c.Object.Put(context.Background(), cosPath, strings.NewReader(""), nil)
			if err != nil {
				rErr = err
			}
		}
		return

	} else {
		size = fileInfo.Size()

		// 仅sync命令执行skip
		if fo.Command == CommandSync {
			absLocalFilePath, _ := filepath.Abs(localFilePath)
			snapshotKey = getUploadSnapshotKey(absLocalFilePath, cosUrl.(*CosUrl).Bucket, cosUrl.(*CosUrl).Object)
			skip, err = skipUpload(snapshotKey, c, fo, fileInfo.ModTime().Unix(), cosPath, localFilePath)
			if err != nil {
				rErr = err
				return
			}
		}

		if skip {
			return
		}

		opt := &cos.MultiUploadOptions{
			OptIni: &cos.InitiateMultipartUploadOptions{
				ACLHeaderOptions: &cos.ACLHeaderOptions{
					XCosACL:              "",
					XCosGrantRead:        "",
					XCosGrantWrite:       "",
					XCosGrantFullControl: "",
					XCosGrantReadACP:     "",
					XCosGrantWriteACP:    "",
				},
				ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
					CacheControl:             fo.Operation.Meta.CacheControl,
					ContentDisposition:       fo.Operation.Meta.ContentDisposition,
					ContentEncoding:          fo.Operation.Meta.ContentEncoding,
					ContentType:              fo.Operation.Meta.ContentType,
					ContentMD5:               fo.Operation.Meta.ContentMD5,
					ContentLength:            fo.Operation.Meta.ContentLength,
					ContentLanguage:          fo.Operation.Meta.ContentLanguage,
					Expect:                   "",
					Expires:                  fo.Operation.Meta.Expires,
					XCosContentSHA1:          "",
					XCosMetaXXX:              fo.Operation.Meta.XCosMetaXXX,
					XCosStorageClass:         fo.Operation.StorageClass,
					XCosServerSideEncryption: "",
					XCosSSECustomerAglo:      "",
					XCosSSECustomerKey:       "",
					XCosSSECustomerKeyMD5:    "",
					XOptionHeader:            nil,
					XCosTrafficLimit:         (int)(fo.Operation.RateLimiting * 1024 * 1024 * 8),
				},
			},
			PartSize:        fo.Operation.PartSize,
			ThreadPoolSize:  fo.Operation.ThreadNum,
			CheckPoint:      true,
			DisableChecksum: fo.Operation.DisableChecksum,
		}

		counter := &Counter{TransferSize: 0}
		// 未跳过则通过监听更新size(仅需要分块文件的通过sdk监听进度)
		if size > fo.Operation.PartSize*1024*1024 {
			opt.OptIni.Listener = &CosListener{fo, counter}
			size = 0
		}

		_, _, err = c.Object.Upload(context.Background(), cosPath, localFilePath, opt)

		if err != nil {
			transferSize = counter.TransferSize
			rErr = err
			return
		}
	}

	if snapshotKey != "" && fo.Operation.SnapshotPath != "" {
		// 上传成功后添加快照
		fo.SnapshotDb.Put([]byte(snapshotKey), []byte(strconv.FormatInt(fileInfo.ModTime().Unix(), 10)), nil)
	}

	return
}
