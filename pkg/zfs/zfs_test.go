package zfs

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func sleep(delay int) {
	time.Sleep(time.Duration(delay) * time.Second)
}

func pow2(x int) int64 {
	return int64(math.Pow(2, float64(x)))
}

//https://github.com/benbjohnson/testing
// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// nok fails the test if an err is nil.
func nok(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func zpoolTest(t *testing.T, fn func()) {
	tempfiles := make([]string, 3)
	for i := range tempfiles {
		f, _ := ioutil.TempFile("/tmp/", "zfs-")
		defer f.Close()
		err := f.Truncate(pow2(30))
		ok(t, err)
		tempfiles[i] = f.Name()
		defer os.Remove(f.Name())
	}

	_, err := CreateZpool("test", nil, tempfiles...)
	ok(t, err)
	defer pool.Destroy()
	ok(t, err)
	fn()

}

func TestDatasets(t *testing.T) {
	zpoolTest(t, func() {
		_, err := Datasets("")
		ok(t, err)

		ds, err := GetDataset("test")
		ok(t, err)
		equals(t, DatasetFilesystem, ds.Type)
		equals(t, "", ds.Origin)
		if runtime.GOOS != "solaris" {
			assert(t, ds.Logicalused != 0, "Logicalused is not greater than 0")
		}
	})
}

func TestDatasetGetProperty(t *testing.T) {
	zpoolTest(t, func() {
		ds, err := GetDataset("test")
		ok(t, err)

		prop, err := ds.GetProperty("foobarbaz")
		nok(t, err)
		equals(t, "", prop)

		prop, err = ds.GetProperty("compression")
		ok(t, err)
		equals(t, "off", prop)
	})
}

func TestSnapshots(t *testing.T) {

	zpoolTest(t, func() {
		snapshots, err := Snapshots("")
		ok(t, err)

		for _, snapshot := range snapshots {
			equals(t, DatasetSnapshot, snapshot.Type)
		}
	})
}

func TestFilesystems(t *testing.T) {
	zpoolTest(t, func() {
		f, err := CreateFilesystem("test/filesystem-test", nil)
		ok(t, err)

		filesystems, err := Filesystems("")
		ok(t, err)

		for _, filesystem := range filesystems {
			equals(t, DatasetFilesystem, filesystem.Type)
		}

		ok(t, f.Destroy(DestroyDefault))
	})
}

func TestCreateFilesystemWithProperties(t *testing.T) {
	zpoolTest(t, func() {
		props := map[string]string{
			"compression": "lz4",
		}

		f, err := CreateFilesystem("test/filesystem-test", props)
		ok(t, err)

		equals(t, "lz4", f.Compression)

		filesystems, err := Filesystems("")
		ok(t, err)

		for _, filesystem := range filesystems {
			equals(t, DatasetFilesystem, filesystem.Type)
		}

		ok(t, f.Destroy(DestroyDefault))
	})
}

func TestVolumes(t *testing.T) {
	zpoolTest(t, func() {
		v, err := CreateVolume("test/volume-test", uint64(pow2(23)), nil)
		ok(t, err)

		// volumes are sometimes "busy" if you try to manipulate them right away
		sleep(1)

		equals(t, DatasetVolume, v.Type)
		volumes, err := Volumes("")
		ok(t, err)

		for _, volume := range volumes {
			equals(t, DatasetVolume, volume.Type)
		}

		ok(t, v.Destroy(DestroyDefault))
	})
}

func TestSnapshot(t *testing.T) {
	zpoolTest(t, func() {
		f, err := CreateFilesystem("test/snapshot-test", nil)
		ok(t, err)

		filesystems, err := Filesystems("")
		ok(t, err)

		for _, filesystem := range filesystems {
			equals(t, DatasetFilesystem, filesystem.Type)
		}

		s, err := f.Snapshot("test", false)
		ok(t, err)

		equals(t, DatasetSnapshot, s.Type)

		equals(t, "test/snapshot-test@test", s.Name)

		ok(t, s.Destroy(DestroyDefault))

		ok(t, f.Destroy(DestroyDefault))
	})
}

func TestClone(t *testing.T) {
	zpoolTest(t, func() {
		f, err := CreateFilesystem("test/snapshot-test", nil)
		ok(t, err)

		filesystems, err := Filesystems("")
		ok(t, err)

		for _, filesystem := range filesystems {
			equals(t, DatasetFilesystem, filesystem.Type)
		}

		s, err := f.Snapshot("test", false)
		ok(t, err)

		equals(t, DatasetSnapshot, s.Type)
		equals(t, "test/snapshot-test@test", s.Name)

		c, err := s.Clone("test/clone-test", nil)
		ok(t, err)

		equals(t, DatasetFilesystem, c.Type)

		ok(t, c.Destroy(DestroyDefault))

		ok(t, s.Destroy(DestroyDefault))

		ok(t, f.Destroy(DestroyDefault))
	})
}

func TestSendSnapshot(t *testing.T) {
	zpoolTest(t, func() {
		f, err := CreateFilesystem("test/snapshot-test", nil)
		ok(t, err)

		filesystems, err := Filesystems("")
		ok(t, err)

		for _, filesystem := range filesystems {
			equals(t, DatasetFilesystem, filesystem.Type)
		}

		s, err := f.Snapshot("test", false)
		ok(t, err)

		file, _ := ioutil.TempFile("/tmp/", "zfs-")
		defer file.Close()
		err = file.Truncate(pow2(30))
		ok(t, err)
		defer os.Remove(file.Name())

		err = s.SendSnapshot(file)
		ok(t, err)

		ok(t, s.Destroy(DestroyDefault))

		ok(t, f.Destroy(DestroyDefault))
	})
}

func TestChildren(t *testing.T) {
	zpoolTest(t, func() {
		f, err := CreateFilesystem("test/snapshot-test", nil)
		ok(t, err)

		s, err := f.Snapshot("test", false)
		ok(t, err)

		equals(t, DatasetSnapshot, s.Type)
		equals(t, "test/snapshot-test@test", s.Name)

		children, err := f.Children(0)
		ok(t, err)

		equals(t, 1, len(children))
		equals(t, "test/snapshot-test@test", children[0].Name)

		ok(t, s.Destroy(DestroyDefault))
		ok(t, f.Destroy(DestroyDefault))
	})
}

func TestListZpool(t *testing.T) {
	zpoolTest(t, func() {
		pools, err := ListZpools()
		ok(t, err)
		for _, pool := range pools {
			if pool.Name == "test" {
				equals(t, "test", pool.Name)
				return
			}
		}
		t.Fatal("Failed to find test pool")
	})
}

func TestRollback(t *testing.T) {
	zpoolTest(t, func() {
		f, err := CreateFilesystem("test/snapshot-test", nil)
		ok(t, err)

		filesystems, err := Filesystems("")
		ok(t, err)

		for _, filesystem := range filesystems {
			equals(t, DatasetFilesystem, filesystem.Type)
		}

		s1, err := f.Snapshot("test", false)
		ok(t, err)

		_, err = f.Snapshot("test2", false)
		ok(t, err)

		s3, err := f.Snapshot("test3", false)
		ok(t, err)

		err = s3.Rollback(false)
		ok(t, err)

		err = s1.Rollback(false)
		assert(t, err != nil, "should error when rolling back beyond most recent without destroyMoreRecent = true")

		err = s1.Rollback(true)
		ok(t, err)

		ok(t, s1.Destroy(DestroyDefault))

		ok(t, f.Destroy(DestroyDefault))
	})
}

func TestDiff(t *testing.T) {
	zpoolTest(t, func() {
		fs, err := CreateFilesystem("test/origin", nil)
		ok(t, err)

		linkedFile, err := os.Create(filepath.Join(fs.Mountpoint, "linked"))
		ok(t, err)

		movedFile, err := os.Create(filepath.Join(fs.Mountpoint, "file"))
		ok(t, err)

		snapshot, err := fs.Snapshot("snapshot", false)
		ok(t, err)

		unicodeFile, err := os.Create(filepath.Join(fs.Mountpoint, "i ❤ unicode"))
		ok(t, err)

		err = os.Rename(movedFile.Name(), movedFile.Name()+"-new")
		ok(t, err)

		err = os.Link(linkedFile.Name(), linkedFile.Name()+"_hard")
		ok(t, err)

		inodeChanges, err := fs.Diff(snapshot.Name)
		ok(t, err)
		equals(t, 4, len(inodeChanges))

		unicodePath := "/test/origin/i\x040\x1c2\x135\x144\x040unicode"
		wants := map[string]*InodeChange{
			"/test/origin/linked": &InodeChange{
				Type:                 File,
				Change:               Modified,
				ReferenceCountChange: 1,
			},
			"/test/origin/file": &InodeChange{
				Type:    File,
				Change:  Renamed,
				NewPath: "/test/origin/file-new",
			},
			"/test/origin/i ❤ unicode": &InodeChange{
				Path:   "❤❤ unicode ❤❤",
				Type:   File,
				Change: Created,
			},
			unicodePath: &InodeChange{
				Path:   "❤❤ unicode ❤❤",
				Type:   File,
				Change: Created,
			},
			"/test/origin/": &InodeChange{
				Type:   Directory,
				Change: Modified,
			},
		}
		for _, change := range inodeChanges {
			want := wants[change.Path]
			want.Path = change.Path
			delete(wants, change.Path)

			equals(t, want, change)
		}

		equals(t, 1, len(wants))
		for _, want := range wants {
			equals(t, "❤❤ unicode ❤❤", want.Path)
		}

		ok(t, movedFile.Close())
		ok(t, unicodeFile.Close())
		ok(t, linkedFile.Close())
		ok(t, snapshot.Destroy(DestroyForceUmount))
		ok(t, fs.Destroy(DestroyForceUmount))
	})
}
