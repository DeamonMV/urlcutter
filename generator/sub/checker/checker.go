package checker

import (
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"strconv"
)

func Unmarshal(bytes []byte) map[string]string {
	log.Info("Raw JSON, which will be unmarshaled: ", string(bytes))
	var data map[string]string
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Warn("Error  %v", err)
	}
	log.Info("Unmarshaled JSON: ", data)
	return data
}

// exapmle of data ./stan-pub -c test-cluster -id stan-ex "dbproc" "{\"dbproc\":\"true\",\"codelen\":\"2\"}"
func Checker(mapa map[string]string) (int, string){

	codebool, codeint := checkCodeLen(mapa["codelen"])

	if mapa["dbproc"] == "true" && codebool == true {

		return codeint, ""
	} else {
		log.Warn("Code lenght in message not in the range. CodeLenght is ", mapa["codelen"])
		return 0, "Wrong code"
	}
}

func checkCodeLen (len string) (bool, int){
	//var intlen int
	codeint , err := strconv.Atoi(len)
	if err != nil{
		log.Warn("Can not convert String int Int, error is: ", err )
	}

	if codeint >= 0 && codeint <=9 {
		return true, codeint
	} else {
		return false, 0
	}
}
