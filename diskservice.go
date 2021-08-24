package diskservice

import (
	"errors"
	"github.com/daqnext/meson.network-lts-disk-service/filedefine"
	"github.com/daqnext/meson.network-lts-disk-service/filemgr"
	"github.com/daqnext/meson.network-lts-disk-service/spaceholdermgr"
	"github.com/shirou/gopsutil/v3/disk"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

const spaceHolderFolder = "spaceholder"
const fileFolder = "files"

const headSpace = 300 * filedefine.UnitM

type DiskService struct {
	rootPath              string
	fileMgr               *filemgr.FileMgr
	spaceHolderMgr        *spaceholdermgr.SpaceHolderMgr
	spaceHolderHandleLock sync.Mutex
	cdnSpaceProvide       uint64
}

func New() (ds *DiskService) {
	ds = &DiskService{}
	return ds
}

//SetRootPath give a path as root folder
func (ds *DiskService) SetRootPath(rootPath string, cdnSpaceProvideGB int) error {
	err := os.MkdirAll(rootPath, 0755)
	if err != nil {
		return err
	}
	ds.rootPath = rootPath
	ds.cdnSpaceProvide = uint64(cdnSpaceProvideGB * filedefine.UnitG)

	//SpaceHolderPath
	spaceHolderPath := path.Join(rootPath, spaceHolderFolder)
	spaceHolderMgr, err := spaceholdermgr.New(spaceHolderPath)
	if err != nil {
		return err
	}
	ds.spaceHolderMgr = spaceHolderMgr
	err = ds.spaceHolderMgr.InitSpaceHolder()
	if err != nil {
		return err
	}

	//FileFolderPath
	fileFolderPath := path.Join(rootPath, fileFolder)
	fileMgr, err := filemgr.New(fileFolderPath)
	if err != nil {
		return err
	}
	ds.fileMgr = fileMgr

	free := uint64(0)
	d, err := disk.Usage(rootPath)
	if err != nil {
		return err
	} else {
		free = d.Free
	}

	cdnSpaceUsed, err := ds.fileMgr.GetSavedFileTotalSize(true)
	if err != nil {
		return err
	}
	holdFileSize, err := ds.spaceHolderMgr.GetAllSpaceHolderFileSize(true)
	if err != nil {
		return err
	}

	total := cdnSpaceUsed + holdFileSize + free
	if total < ds.cdnSpaceProvide {
		return errors.New("disk no enough space")
	}

	return nil
}

// FullSpaceHolder Full the space with space hold files
func (ds *DiskService) FullSpaceHolderInLaunch() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		count := 0
		for true {
			holdSize, err := ds.spaceHolderMgr.GetAllSpaceHolderFileSize(false)
			cacheSize, err := ds.fileMgr.GetSavedFileTotalSize(false)
			if holdSize+cacheSize > ds.cdnSpaceProvide-headSpace {
				wg.Done()
				break
			}
			if holdSize > 4*filedefine.UnitG {
				time.Sleep(300 * time.Millisecond)
			}
			err = ds.spaceHolderMgr.GenSpaceHolderFile()
			if err != nil {
				log.Println(err)
				count++
				if count > 5 {
					log.Println("Can not create space holder files, please check the free space of disk")
					return
				}
				continue
			}
			count = 0
		}
	}()
	return wg
}

func (ds *DiskService) FullSpaceHolderAfterDeleteCache() {

}

// FreeSpace Free some space form space holder to save new file
func (ds *DiskService) FreeSpace(needSpace uint64) error {
	holdSize, _ := ds.spaceHolderMgr.GetAllSpaceHolderFileSize(false)
	if needSpace > holdSize {
		return errors.New("not enough disk space")
	}
	ds.spaceHolderHandleLock.Lock()
	defer ds.spaceHolderHandleLock.Unlock()
	for true {
		holdSize, _ := ds.spaceHolderMgr.GetAllSpaceHolderFileSize(false)
		cacheSize, _ := ds.fileMgr.GetSavedFileTotalSize(false)
		freeSpace := ds.cdnSpaceProvide - (holdSize + cacheSize)
		if freeSpace > needSpace+headSpace {
			return nil
		}
		if holdSize <= 0 {
			return errors.New("not enough disk space")
		}
		ds.spaceHolderMgr.DeleteSpaceHolderFile()
	}
	return nil
}

// SaveFile Download file from originUrl and save
func (ds *DiskService) SaveFile(fileSavePath string, originUrl string) {

}

// DeleteFile Delete file form disk
func (ds *DiskService) DeleteFile(filePath string) {

}
