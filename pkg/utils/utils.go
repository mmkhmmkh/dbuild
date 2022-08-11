package utils

import (
	"fmt"
	"time"
)

const (
	DbuildRepo        = "mmkhmmkh/dbuild"
	ControllerContext = "controller"
	WorkerContext     = "worker"
	DbuildBranch      = "master"
	DbuildPrefix      = "dbuild-"
	DbuildDir         = "/dbuild"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}
