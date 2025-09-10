package util

import "github.com/tencentyun/cos-go-sdk-v5"

func getThreadNumByPartSize(totalSize, partSize int64) (int, error) {
	var threadNum int
	_, partNum, err := cos.SplitSizeIntoChunksToDownload(totalSize, partSize*1024*1024)
	if err != nil {
		return threadNum, err
	}

	if partNum < 2 {
		threadNum = 1
	} else if partNum < 4 {
		threadNum = 2
	} else if partNum <= 20 {
		threadNum = 4
	} else if partNum <= 300 {
		threadNum = 8
	} else if partNum <= 500 {
		threadNum = 10
	} else {
		threadNum = 12
	}
	return threadNum, err
}
