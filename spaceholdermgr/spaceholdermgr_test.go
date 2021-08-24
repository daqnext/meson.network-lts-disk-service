package spaceholdermgr

import (
	"github.com/daqnext/meson.network-lts-disk-service/filedefine"
	"log"
	"testing"
)

func TestGenSpaceHolderFile(t *testing.T) {
	sh, err := New("./spaceholder")
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	err = sh.InitSpaceHolder()
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	count := 0
	for true {
		if sh.HoldFileSize >= 5*filedefine.UnitG {
			break
		}
		err := sh.GenSpaceHolderFile()
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

	log.Println(sh.HoldFileSize)
	log.Println(sh.SpaceHoldFiles)
}

func TestDeleteSpaceHolderFile(t *testing.T) {
	sh, err := New("./spaceholder")
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}

	err = sh.InitSpaceHolder()
	if err != nil {
		log.Println(err)
		t.Error(err)
		return
	}
	log.Println(sh.HoldFileSize)
	log.Println(sh.SpaceHoldFiles)

	for true {
		if sh.HoldFileSize < 1*filedefine.UnitG {
			break
		}
		err := sh.DeleteSpaceHolderFile()
		if err != nil {
			log.Println(err)
			t.Error(err)
		}
		if sh.HoldFileSize <= 0 {
			break
		}
	}

	log.Println(sh.HoldFileSize)
	log.Println(sh.SpaceHoldFiles)
}
