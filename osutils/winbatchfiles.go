package osutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

var runningCount int = 0

func HasRunningBatchFiles() bool {
	return runningCount > 0
}

func RunTempBatchFileAndGetOutput(commandsContent string, tempDir string) ([]byte, error) {
	var err error
	tempBatchFilesDir := path.Join(tempDir, "_TempBatchFiles")
	if !PathExists(tempBatchFilesDir) {
		err = os.MkdirAll(tempBatchFilesDir, 0600)
		if err != nil {
			return nil, err
		}
	}
	now := time.Now()

	fileNameOnly := fmt.Sprintf("%04d-%02d-%02d %02d_%02d_%02d.%09d.bat", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond())

	tmpBaseBatchFilePath := path.Join(tempBatchFilesDir, "1 - "+fileNameOnly)
	tmpBatchFileContentPath := path.Join(tempBatchFilesDir, "2 - "+fileNameOnly)

	baseBatchContent := "@echo off\r\n"
	//baseBatchContent += "start \"\" /wait \"" + tmpBatchFileContentPath + "\""
	baseBatchContent += "call \"" + tmpBatchFileContentPath + "\""
	baseBatchContent += "\r\n"
	err = ioutil.WriteFile(tmpBaseBatchFilePath, []byte(baseBatchContent), 0600)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpBaseBatchFilePath)

	batchFileContent := "@echo off\r\n"
	batchFileContent += commandsContent
	batchFileContent += "\r\n"
	err = ioutil.WriteFile(tmpBatchFileContentPath, []byte(batchFileContent), 0600)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpBatchFileContentPath)

	cmd := exec.Command(tmpBaseBatchFilePath)
	return cmd.CombinedOutput()
}
