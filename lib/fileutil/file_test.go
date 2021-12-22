package fileutil

import "testing"

func Test_WriteFile(t *testing.T) {
	err := WriteFile("/tmp/harp_test_file", "test")
	if err != nil {
		t.Fatal("Failed to write file!")
	}

	err = WriteFile("/proc/harp_test_file", "test")
	if err != nil {
		t.Log("Tried to write non-exist file in /proc foler")
	}

	err = WriteFile("/proc/cpuinfo", "test")
	if err != nil {
		t.Log("Tried to write file in /proc foler")
	}
}

func Test_WriteFileAppend(t *testing.T) {
	err := WriteFileAppend("/tmp/harp_test_file", "_test_append")
	if err != nil {
		t.Fatal("Failed to write file with append mode!")
	}

	err = WriteFileAppend("/proc/harp_test_file", "_test_append")
	if err != nil {
		t.Log("Tried to write non-exist file with append mode in /proc foler")
	}

	err = WriteFileAppend("/proc/cpuinfo", "_test_append")
	if err != nil {
		t.Log("Tried to write file with append mode in /proc foler")
	}
}

func Test_DeleteFile(t *testing.T) {
	err := DeleteFile("/tmp/harp_test_file")
	if err != nil {
		t.Error("Failed to delete file!")
	}

	err = DeleteFile("/proc/cpuinfo")
	if err != nil {
		t.Log("Tried to delete file in /proc folder")
	}
}
