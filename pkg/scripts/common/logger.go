package common

import (
	"log"
	"os"
)

var ScriptsLogger = log.New(os.Stdout, "", log.LstdFlags)

func ScriptsDebug(format string, v ...any) {
	ScriptsLogger.Printf("[DEBUG] "+format, v...)
}

func ScriptsWarn(format string, v ...any) {
	ScriptsLogger.Printf("[WARN] "+format, v...)
}
