package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
)

var DEBUG = false

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Debug(s string) {
	if DEBUG {
		fmt.Printf("---------\n%s\n---------\n\n", s)
	}
}

func GetTerraformTokenFromConfig() string {
	homeDir, err := os.UserHomeDir()
	Check(err)

	var tfCredFile string
	if runtime.GOOS == "windows" {
		tfCredFile = fmt.Sprintf("%s\\AppData\\Roaming\\terraform.d\\credentials.tfrc.json", homeDir)
	} else {
		tfCredFile = fmt.Sprintf("%s/.terraform.d/credentials.tfrc.json", homeDir)
	}

	dat, err := os.Open(tfCredFile)
	Check(err)
	defer dat.Close()

	byteValue, _ := io.ReadAll(dat)
	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	Check(err)

	token := result["credentials"].(map[string]interface{})["app.terraform.io"].(map[string]interface{})["token"].(string)
	return token
}
