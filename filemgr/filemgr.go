package filemgr

import (
	"github.com/daqnext/meson.network-lts-disk-service/utils"
	"os"
)

type FileMgr struct {
	FileFolderPath string
	CacheFileSize  uint64
}

func New(fileFolderPath string) (*FileMgr, error) {
	err := os.MkdirAll(fileFolderPath, 0755)
	if err != nil {
		return nil, err
	}
	fileMgr := &FileMgr{
		FileFolderPath: fileFolderPath,
	}
	return fileMgr, nil
}

func (f *FileMgr) SaveFile(filePath string) {
	//start download
}

func (f *FileMgr) DeleteFile(filePath string) error {
	//delete file in disk
	return os.Remove(filePath)

	//delete header file in disk
}

//GetSavedFileTotalSize get all saved file size (Byte)
func (f *FileMgr) GetSavedFileTotalSize(refresh bool) (uint64, error) {
	if !refresh {
		return f.CacheFileSize, nil
	}
	cacheSize, err := utils.GetDirSize(f.FileFolderPath)
	if err != nil {
		return 0, err
	}
	f.CacheFileSize = cacheSize
	return f.CacheFileSize, nil
}

//GetUnAccessedFiles get all unaccessed file in given interval second
func (f *FileMgr) GetUnAccessedFiles(intervalSecond int) {

}
