package main

import (
	"os"
	"strings"
	"sync"
)

var writeLock sync.Mutex

func writeToFile(filename, content string) {
	writeLock.Lock()
	defer writeLock.Unlock()

	writeLog("Writing to", filename, content)
	if strings.TrimSpace(filename) == "" || strings.TrimSpace(content) == "" {
		errorLog("Writing to", filename, "failed because content or filename is empty(probably nothing to write)")
		return
	}

	// ensure path is ready
	folderPath := strings.Split(filename, "/")
	folderPath = folderPath[:len(folderPath)-1]

	folderPathS := strings.Join(folderPath, "/")

	err := os.MkdirAll(folderPathS, 0777)
	handleErrorAndPanic(err)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	handleErrorAndPanic(err)

	_, err = file.WriteString(content + "\n")
	handleErrorAndPanic(err)

	err = file.Close()
	handleErrorAndPanic(err)
}

//func writeSlice(sl []string, fileName string) {
//	for _, value := range sl {
//		writeLog("Writing", value, "to", fileName)
//		writeToFile(fileName, value)
//	}
//}
