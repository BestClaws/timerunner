package main

import (
	"time"
	"timerunner/pkg/jsondisk"
)

type RunInfo struct {
	Period      Duration `json:"period"`
	Offset      Duration `json:"offset"`
	PackageName string   `json:"packageName"`
	ExecPath    string   `json:"execPath"`
}

func (r RunInfo) GetOffset() time.Duration {
	return r.Offset.Duration
}

func (r RunInfo) GetPeriod() time.Duration {
	return r.Period.Duration
}

// loads run info from disk at startup.
func loadRunInfo() []RunInfo {

	runInfo, err := jsondisk.Load[[]RunInfo]("data/descriptions.json")
	if err != nil {
		return nil
	}
	return *runInfo

}
