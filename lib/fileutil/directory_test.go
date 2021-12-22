package fileutil

import (
	"testing"
)

func Test_CreateDirIfNotExist(t *testing.T) {
	err := CreateDirIfNotExist("/tmp/harp_test")
	if err != nil {
		t.Fatal("Failed to create dir!")
	}

	err = CreateDirIfNotExist("/proc/harp_test")
	if err != nil {
		t.Log("Tried to create dir in /proc folder")
	}
}

func Test_DeleteDir(t *testing.T) {
	err := DeleteDir("/tmp/harp_test")
	if err != nil {
		t.Fatal("Failed to delete dir!")
	}

	err = DeleteDir("/proc/cpuinfo")
	if err != nil {
		t.Log("Tried to delete file in /proc folder")
	}
}
