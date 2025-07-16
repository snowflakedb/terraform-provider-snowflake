package common

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)
