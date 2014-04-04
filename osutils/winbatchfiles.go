package osutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

func RunTempWindowsBatchFile(commandsContent string, mustWait bool, tempDir string) error {
	var err error
	tempBatchFilesDir := path.Join(tempDir, "_TempBatchFiles")
	if !PathExists(tempBatchFilesDir) {
		err = os.MkdirAll(tempBatchFilesDir, 0600)
		if err != nil {
			return err
		}
	}
	now := time.Now()
	tmpBatchFile := path.Join(tempBatchFilesDir, fmt.Sprintf("%04d-%02d-%02d %02d_%02d_%02d.%09d.bat", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()))

	batchFileContent := "@echo off\r\n"
	batchFileContent += "start \"\" /wait "
	batchFileContent += commandsContent
	err = ioutil.WriteFile(tmpBatchFile, []byte(batchFileContent), 0600)
	if err != nil {
		return err
	}
	defer os.Remove(tmpBatchFile)

	cmd := exec.Command(tmpBatchFile)
	err = cmd.Start()
	if err != nil {
		return err
	}

	if mustWait {
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
