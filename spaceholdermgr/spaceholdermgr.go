package spaceholdermgr

import (
	"fmt"
	"github.com/daqnext/meson.network-lts-disk-service/filedefine"
	"github.com/daqnext/meson.network-lts-disk-service/utils"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
)

type SpaceHolderMgr struct {
	SpaceHolderPath string
	SpaceHoldFiles  []string
	HoldFileSize    uint64
	handleFileLock  sync.Mutex
}

var array = make([]byte, filedefine.UnitM)

const eachHoldFileSize = uint64(100 * filedefine.UnitM)

func New(spaceHolderPath string) (*SpaceHolderMgr, error) {
	err := os.MkdirAll(spaceHolderPath, 0755)
	if err != nil {
		return nil, err
	}
	spaceHolderMgr := &SpaceHolderMgr{
		SpaceHolderPath: spaceHolderPath,
	}

	return spaceHolderMgr, nil
}

func (s *SpaceHolderMgr) InitSpaceHolder() error {
	//disk space holder
	if !utils.Exists(s.SpaceHolderPath) {
		err := os.Mkdir(s.SpaceHolderPath, 0755)
		if err != nil {
			return err
		}
	}

	holdFiles, err := ioutil.ReadDir(s.SpaceHolderPath)
	if err != nil {
		return err
	}

	s.SpaceHoldFiles = []string{}
	s.HoldFileSize = 0
	for _, file := range holdFiles {
		s.SpaceHoldFiles = append(s.SpaceHoldFiles, file.Name())
		s.HoldFileSize += uint64(file.Size())
	}
	return nil
}

func (s *SpaceHolderMgr) GenSpaceHolderFile() error {
	s.handleFileLock.Lock()
	defer s.handleFileLock.Unlock()

	holdFileCount := len(s.SpaceHoldFiles)
	name := fmt.Sprintf("%010d%d", holdFileCount+1, rand.Intn(99999999))
	name = utils.GetStringHash(name)
	name = name + ".bin"
	fileName := s.SpaceHolderPath + "/" + name
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	os.Chmod(fileName, 0755)
	defer f.Close()
	fileSize := uint64(0)
	for i := 0; i < 100; i++ {
		size, err := f.Write(array)
		if err != nil {
			return err
		}
		fileSize += uint64(size)
	}

	s.SpaceHoldFiles = append(s.SpaceHoldFiles, name)
	s.HoldFileSize += fileSize
	return nil
}

func (s *SpaceHolderMgr) DeleteSpaceHolderFile() error {
	s.handleFileLock.Lock()
	defer s.handleFileLock.Unlock()

	if len(s.SpaceHoldFiles) == 0 {
		files, err := ioutil.ReadDir(s.SpaceHolderPath)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return nil
		}
		s.HoldFileSize = 0
		for _, v := range files {
			if v.IsDir() {
				continue
			}
			s.SpaceHoldFiles = append(s.SpaceHoldFiles, v.Name())
			s.HoldFileSize += uint64(v.Size())
		}
	}
	name := s.SpaceHoldFiles[len(s.SpaceHoldFiles)-1]
	fileName := s.SpaceHolderPath + "/" + name

	fileStat, err := os.Stat(fileName)
	if err != nil {
		s.SpaceHoldFiles = s.SpaceHoldFiles[:len(s.SpaceHoldFiles)-1]
		return err
	}

	s.HoldFileSize -= uint64(fileStat.Size())
	s.SpaceHoldFiles = s.SpaceHoldFiles[:len(s.SpaceHoldFiles)-1]
	os.Remove(fileName)

	return nil
}

//GetAllSpaceHolderFileSize get the space holder file size (byte)
func (s *SpaceHolderMgr) GetAllSpaceHolderFileSize(refresh bool) (uint64, error) {
	if !refresh && s.HoldFileSize > 0 {
		return s.HoldFileSize, nil
	}
	err := s.InitSpaceHolder()
	if err != nil {
		return 0, err
	}
	return s.HoldFileSize, nil
}
