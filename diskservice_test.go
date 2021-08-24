package diskservice

import (
	"github.com/daqnext/meson.network-lts-disk-service/filedefine"
	"testing"
)

func TestInit(t *testing.T) {
	ds := New()
	err := ds.SetRootPath("./root", 2)
	if err != nil {
		t.Error(err)
		return
	}

	wg := ds.FullSpaceHolderInLaunch()

	wg.Wait()
}

func TestFreeSpace(t *testing.T) {
	ds := New()
	err := ds.SetRootPath("./root", 2)
	if err != nil {
		t.Error(err)
		return
	}

	wg := ds.FullSpaceHolderInLaunch()
	wg.Wait()

	err = ds.FreeSpace(3 * filedefine.UnitG)
	if err != nil {
		t.Error(err)
	}
}

func TestSaveFile(t *testing.T) {

}

func TestDeleteFile(t *testing.T) {

}
