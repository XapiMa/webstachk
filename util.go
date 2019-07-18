package webstachk

import (
	"fmt"
	"log"
	"os"
	"time"
)

func logFatal(err error) {
	log.Fatalf("Error: webstachk %s %s", time.Now().Format("2006/01/02 15:04:05"), err)
}

func logPrint(err error) {
	log.Printf("Error: webstachk %s %s", time.Now().Format("2006/01/02 15:04:05"), err)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func joinCode(codes []int, s string) string {
	str := ""
	for j, code := range codes {
		str += fmt.Sprintf("%d", code)
		if j != len(codes)-1 {
			str += s
		}
	}
	return str
}
func appendFile(outputPath, outputString string) error {
	if outputPath == "" {
		fmt.Printf("%s\n", outputString)
	} else {
		file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Write(([]byte)(outputString + "\n"))
		file.Close()
	}
	return nil
}
