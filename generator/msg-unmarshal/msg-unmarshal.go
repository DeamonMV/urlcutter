package msg_unmarshal

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

// Deprecated -- exapmle of data ./stan-pub -c test-cluster -id stan-ex "dbproc" "{\"dbproc\":\"true\",\"codelen\":\"2\"}"
// exapmle of data ./stan-pub -c test-cluster -id stan-ex "dbproc" "{\"name\":\"dbproc\",\"title\":\"code\",\"msg\":\"2\"}"

func GeneratorDataCheck(mapa map[string]string) (int, string){

	codebool, codeint := checkCodeLen(mapa["msg"])

	if mapa["name"] == "dbproc" && codebool == true {

		return codeint, ""
	} else {
		log.Warn("Code lenght in message not in the range. CodeLenght is ", mapa["msg"])
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

// Ф-ция на доставания кода от герератора
// по всей идее отсюда должнен уходить просто стринг, которых содержит Код
// подумать о преверек на длинну кода - т.е. в сообщении передавать длинну кода и сранивать ее с фактической
func DbprocGetCode(mapa map[string]string) string {

	if mapa["name"] == "generator" {

		return mapa["msg"]
	} else   {
		log.Warn("Code in message not in the range. Code is ", mapa["msg"])
		return ""
	}
}

// подумать о преверек на длинну кода - т.е. в сообщении передавать длинну кода и сранивать ее с фактической
func FaceDataCheck(mapa map[string]string) string {

	if mapa["name"] == "Dbproc" && mapa["msg"] == "existed" {

		return mapa["code"]
	} else {
		log.Warn("Code in message not in the range. Code is ", mapa["msg"])
		return ""
	}
}

// подумать о преверек на длинну кода - т.е. в сообщении передавать длинну кода и сранивать ее с фактической
func DeliveryDataCheck(mapa map[string]string) string {

	if mapa["name"] == "Dbproc" && mapa["msg"] == "new" {
		return ""
	} else if mapa["name"] == "Dbproc" && mapa["msg"] == "existed" {
		return mapa["code"]
	} else {
		log.Warn("Code in message not in the range. Code is ", mapa["msg"])
		return "ERROR"
	}
}

// подумать о преверек на длинну кода - т.е. в сообщении передавать длинну кода и сранивать ее с фактической
func DbprocDataCheck(mapa map[string]string) string {

	if mapa["name"] == "Dbproc" {

		return mapa["code"]
	} else {
		log.Warn("Code in message not in the range. Code is ", mapa["msg"])
		return ""
	}
}