package common

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ExtractFile(content []byte, fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return ioutil.WriteFile(fileName, content, 0644)
	}
	fmt.Printf("File %s exists/couldn't be checked. Skipping to write\n", fileName)
	return nil
}
