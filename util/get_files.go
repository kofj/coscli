package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var once sync.Once

func fileStatistic(localPath string, fo *FileOperations) {
	f, err := os.Stat(localPath)
	if err != nil {
		fo.Monitor.setScanError(err)
		return
	}
	if f.IsDir() {
		if !strings.HasSuffix(localPath, string(os.PathSeparator)) {
			localPath += string(os.PathSeparator)
		}

		err := getFileListStatistic(localPath, fo)
		if err != nil {
			fo.Monitor.setScanError(err)
			return
		}
	} else {
		fo.Monitor.updateScanSizeNum(f.Size(), 1)
	}

	fo.Monitor.setScanEnd()
	freshProgress()
}

func getFileListStatistic(dpath string, fo *FileOperations) error {
	if fo.Operation.OnlyCurrentDir {
		return getCurrentDirFilesStatistic(dpath, fo)
	}

	name := dpath
	symlinkDiretorys := []string{dpath}
	walkFunc := func(fpath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		realFileSize := f.Size()
		dpath = filepath.Clean(dpath)
		fpath = filepath.Clean(fpath)
		fileName, err := filepath.Rel(dpath, fpath)
		if err != nil {
			return fmt.Errorf("list file error: %s, info: %s", fpath, err.Error())
		}

		if f.IsDir() {
			if fpath != dpath {
				if matchPatterns(filepath.Join(dpath, fileName), fo.Operation.Filters) {
					fo.Monitor.updateScanNum(1)
				}
			}
			return nil
		}

		if fo.Operation.DisableAllSymlink && (f.Mode()&os.ModeSymlink) != 0 {
			return nil
		}

		// 处理软链文件或文件夹
		if f.Mode()&os.ModeSymlink != 0 {

			realInfo, err := os.Stat(fpath)
			if err != nil {
				return err
			}

			if realInfo.IsDir() {
				realFileSize = 0
			} else {
				realFileSize = realInfo.Size()
			}

			if fo.Operation.EnableSymlinkDir && realInfo.IsDir() {
				// 软链文件夹，如果有"/"后缀，os.Lstat 将判断它是一个目录
				if !strings.HasSuffix(name, string(os.PathSeparator)) {
					name += string(os.PathSeparator)
				}
				linkDir := name + fileName + string(os.PathSeparator)
				symlinkDiretorys = append(symlinkDiretorys, linkDir)
				return nil
			}
		}
		if matchPatterns(filepath.Join(dpath, fileName), fo.Operation.Filters) {
			fo.Monitor.updateScanSizeNum(realFileSize, 1)
		}
		return nil
	}

	var err error
	for {
		symlinks := symlinkDiretorys
		symlinkDiretorys = []string{}
		for _, v := range symlinks {
			err = filepath.Walk(v, walkFunc)
			if err != nil {
				return err
			}
		}
		if len(symlinkDiretorys) == 0 {
			break
		}
	}
	return err
}

func getCurrentDirFilesStatistic(dpath string, fo *FileOperations) error {
	if !strings.HasSuffix(dpath, string(os.PathSeparator)) {
		dpath += string(os.PathSeparator)
	}

	fileList, err := ioutil.ReadDir(dpath)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileList {
		if !fileInfo.IsDir() {
			realInfo, errF := os.Stat(dpath + fileInfo.Name())
			if errF == nil && realInfo.IsDir() {
				// for symlink
				continue
			}
			if matchPatterns(filepath.Join(dpath, fileInfo.Name()), fo.Operation.Filters) {
				fo.Monitor.updateScanSizeNum(fileInfo.Size(), 1)
			}
		}
	}
	return nil
}

func generateFileList(localPath string, chFiles chan<- fileInfoType, chListError chan<- error, fo *FileOperations) {
	defer close(chFiles)
	f, err := os.Stat(localPath)
	if err != nil {
		chListError <- err
		return
	}
	if f.IsDir() {
		if !strings.HasSuffix(localPath, string(os.PathSeparator)) {
			localPath += string(os.PathSeparator)
		}

		err := getFileList(localPath, chFiles, fo)
		if err != nil {
			chListError <- err
			return
		}
	} else {
		dir, fname := filepath.Split(localPath)
		chFiles <- fileInfoType{filePath: fname, dir: dir, size: f.Size(), isDir: f.IsDir()}
	}
	chListError <- nil
}

func getFileList(dpath string, chFiles chan<- fileInfoType, fo *FileOperations) error {
	if fo.Operation.OnlyCurrentDir {
		return getCurrentDirFileList(dpath, chFiles, fo)
	}

	name := dpath
	symlinkDiretorys := []string{dpath}
	walkFunc := func(fpath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		realFileSize := f.Size()
		dpath = filepath.Clean(dpath)
		fpath = filepath.Clean(fpath)
		fileName, err := filepath.Rel(dpath, fpath)
		if err != nil {
			return fmt.Errorf("list file error: %s, info: %s", fpath, err.Error())
		}

		if f.IsDir() {
			if fpath != dpath {
				if matchPatterns(filepath.Join(dpath, fileName), fo.Operation.Filters) {
					if strings.HasSuffix(fileName, "\\") || strings.HasSuffix(fileName, "/") {
						chFiles <- fileInfoType{filePath: fileName, dir: name, size: 0, lastModified: f.ModTime().Unix(), isDir: f.IsDir()}
					} else {
						chFiles <- fileInfoType{filePath: fileName + string(os.PathSeparator), dir: name, size: 0, lastModified: f.ModTime().Unix(), isDir: f.IsDir()}
					}
				}
			}
			return nil
		}

		if fo.Operation.DisableAllSymlink && (f.Mode()&os.ModeSymlink) != 0 {
			return nil
		}

		if fo.Operation.EnableSymlinkDir && (f.Mode()&os.ModeSymlink) != 0 {
			realInfo, err := os.Stat(fpath)
			if err != nil {
				return err
			}

			if realInfo.IsDir() {
				if !strings.HasSuffix(name, string(os.PathSeparator)) {
					name += string(os.PathSeparator)
				}
				linkDir := name + fileName + string(os.PathSeparator)
				symlinkDiretorys = append(symlinkDiretorys, linkDir)
				return nil
			}
		}

		if matchPatterns(filepath.Join(dpath, fileName), fo.Operation.Filters) {
			chFiles <- fileInfoType{filePath: fileName, dir: name, size: realFileSize, lastModified: f.ModTime().Unix(), isDir: f.IsDir()}
		}
		return nil
	}

	var err error
	for {
		symlinks := symlinkDiretorys
		symlinkDiretorys = []string{}
		for _, v := range symlinks {
			err = filepath.Walk(v, walkFunc)
			if err != nil {
				return err
			}
		}
		if len(symlinkDiretorys) == 0 {
			break
		}
	}
	return err
}

func getCurrentDirFileList(dpath string, chFiles chan<- fileInfoType, fo *FileOperations) error {
	if !strings.HasSuffix(dpath, string(os.PathSeparator)) {
		dpath += string(os.PathSeparator)
	}

	fileList, err := ioutil.ReadDir(dpath)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileList {
		if !fileInfo.IsDir() {
			realInfo, errF := os.Stat(dpath + fileInfo.Name())
			if errF == nil && realInfo.IsDir() {
				// for symlink
				continue
			}

			if matchPatterns(filepath.Join(dpath, fileInfo.Name()), fo.Operation.Filters) {
				chFiles <- fileInfoType{filePath: fileInfo.Name(), dir: dpath, size: fileInfo.Size(), lastModified: fileInfo.ModTime().Unix(), isDir: fileInfo.IsDir()}
			}
		}
	}
	return nil
}

func getLocalFileKeys(fileUrl StorageUrl, keys map[string]commonInfoType, fo *FileOperations, objType string) error {
	strPath := fileUrl.ToString()
	if !strings.HasSuffix(strPath, string(os.PathSeparator)) {
		strPath += string(os.PathSeparator)
	}

	chFiles := make(chan fileInfoType, ChannelSize)
	chFinish := make(chan error, 2)
	go ReadLocalFileKeys(chFiles, chFinish, keys, fo, objType)
	go GetFileList(strPath, chFiles, chFinish, fo)
	select {
	case err := <-chFinish:
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadLocalFileKeys 读取本地文件keys
func ReadLocalFileKeys(chFiles <-chan fileInfoType, chFinish chan<- error, keys map[string]commonInfoType, fo *FileOperations, objType string) {

	results := make(chan commonInfoType, 1000)
	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU()*2; i++ { // 根据CPU核心数动态调整
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fileInfo := range chFiles {
				select {
				case results <- commonInfoType{key: fileInfo.filePath, dir: fileInfo.dir, size: fileInfo.size, lastModifiedUnix: fileInfo.lastModified, isDir: fileInfo.isDir}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// 3. 启动结果收集器
	go func() {
		defer close(done)
		totalCount := 0
		lastReport := time.Now()
		batchSize := 1000 // 批量处理大小
		batch := make([]commonInfoType, 0, batchSize)

		for res := range results {
			totalCount++
			batch = append(batch, res)

			if len(batch) >= batchSize || time.Since(lastReport) > 100*time.Millisecond {
				for _, item := range batch {
					keys[item.key] = item
				}
				batch = batch[:0] // 重置批次

				if objType == TypeSrc {
					fo.SyncDeleteObjectInfo.srcCount = totalCount
				} else {
					fo.SyncDeleteObjectInfo.destCount = totalCount
				}
				lastReport = time.Now()

				// 检查数量限制
				if len(keys) > MaxSyncNumbers {
					cancel() // 取消所有工作
					chFinish <- fmt.Errorf("over max sync numbers %d", MaxSyncNumbers)
					return
				}
			}
		}

		// 处理剩余批次
		for _, item := range batch {
			keys[item.key] = item
		}

		if objType == TypeSrc {
			fo.SyncDeleteObjectInfo.srcCount = totalCount
		} else {
			fo.SyncDeleteObjectInfo.destCount = totalCount
		}
		chFinish <- nil
	}()

	// 4. 等待所有工作器完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 5. 等待结果收集完成
	<-done
}

// GetFileList 获取文件列表
func GetFileList(strPath string, chFiles chan<- fileInfoType, chFinish chan<- error, fo *FileOperations) {
	defer close(chFiles)
	err := getFileList(strPath, chFiles, fo)
	if err != nil {
		chFinish <- err
	}
}

func generateFileListByKeys(srcKeys, uploadKeys map[string]commonInfoType, chFiles chan<- fileInfoType, chListError chan<- error, fo *FileOperations) {
	defer close(chFiles)

	// 使用 WaitGroup 等待两个任务完成
	var wg sync.WaitGroup
	wg.Add(2) // 等待两个任务

	// 任务1：扫描统计（在后台执行）
	go func() {
		defer wg.Done()
		for _, v := range srcKeys {
			fo.Monitor.updateScanSizeNum(v.size, 1)
		}
		fo.Monitor.setScanEnd()
		freshProgress()
	}()

	// 任务2：发送文件信息（在后台执行）
	go func() {
		defer wg.Done()
		for k, v := range srcKeys {
			if _, exists := uploadKeys[k]; exists {
				chFiles <- fileInfoType{v.key, v.dir, v.size, v.lastModifiedUnix, v.isDir, false}
			} else {
				chFiles <- fileInfoType{v.key, v.dir, v.size, v.lastModifiedUnix, v.isDir, true}
			}
		}
		// 发送完成信号
		chListError <- nil
	}()

	// 等待两个任务完成
	wg.Wait()

}
