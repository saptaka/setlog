package setlog

import (
	"log"
	"io"
	"os"
	"time"
	"setupconfig"
	"strconv"
	"io/ioutil"
)

var (
	Debug *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
)

type Log_variable struct {
	FILE_NAME      string
	DIR_NAME       string
	DIR_NAME_INFO  string
	DIR_NAME_ERROR string
	DEBUGMODE      bool
}


func initLog(debugProcess,infoProcess,warningProcess,errorProcess  io.Writer) {

	Debug=log.New(debugProcess, "DEBUG: ",log.Ldate|log.Ltime|log.Lshortfile)

	Info=log.New(infoProcess, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	Warning=log.New(warningProcess, "WARNING: ",log.Ldate|log.Ltime|log.Lshortfile)

	Error=log.New(errorProcess, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func logManagement(configfile string)  Log_variable{
	config:=setupconfig.ReadWriteConfig(configfile)
	var log_setting Log_variable
	var now string
	switch config["FORMAT_FILENAME"]{
	case "YYYY-MM-DD":
		now=time.Now().Format("2006-01-02")
	case "YYYY-DD-MM":
		now=time.Now().Format("2006-02-01")
	case "DD-MM-YYYY":
		now=time.Now().Format("02-01-2006")
	case "MM-DD-YYYY":
		now=time.Now().Format("01-02-2006")
	default:
		now=time.Now().Format("02-01-2006")
	}
	log_setting.FILE_NAME = now
	log_setting.DIR_NAME = config["DIR_NAME"]
	log_setting.DIR_NAME_INFO=config["DIR_LOG_INFO"]
	log_setting.DIR_NAME_ERROR=config["DIR_LOG_ERROR"]
	log_setting.DEBUGMODE,_= strconv.ParseBool(config["DEBUG"])

	return log_setting
}

func SetupLog(configFile string) {

	file_name:=logManagement(configFile).FILE_NAME
	dir_name:=logManagement(configFile).DIR_NAME
	dir_name_info:=logManagement(configFile).DIR_NAME_INFO
	dir_name_error:=logManagement(configFile).DIR_NAME_ERROR
	debugMode:=logManagement(configFile).DEBUGMODE

	if debugMode{
		if _, err := os.Stat("Debug"); os.IsNotExist(err) {
			os.Mkdir("Debug", os.ModePerm)
		}
		debug_event, err:=os.OpenFile("Debug"+"/"+file_name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

		if err != nil {
			log.Fatalln("Failed to open log file", debug_event, ":", err)
		}

		debug_write := io.MultiWriter(debug_event, os.Stdout)

		initLog(debug_write, debug_write, debug_write, debug_write)


	}else{
		if _, err := os.Stat(dir_name); os.IsNotExist(err) {
			os.Mkdir(dir_name, os.ModePerm)
		}

		if _, err := os.Stat(dir_name+"/"+dir_name_info); os.IsNotExist(err) {
			os.Mkdir(dir_name+"/"+dir_name_info, os.ModePerm)
		}

		file_event, err := os.OpenFile(dir_name+"/"+dir_name_info+"/"+file_name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			log.Fatalln("Failed to open log file", file_event, ":", err)
		}

		if _, err := os.Stat(dir_name+"/"+dir_name_error); os.IsNotExist(err) {
			os.Mkdir(dir_name+"/"+dir_name_error, os.ModePerm)
		}

		file_error, err := os.OpenFile(dir_name+"/"+dir_name_error+"/"+file_name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			log.Fatalln("Failed to open log file", file_error, ":", err)
		}

		initLog(ioutil.Discard, file_event, os.Stdout, file_error)
	}
}

