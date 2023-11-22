package main

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"time"
	"timerunner/pkg/jsondisk"
)

type Scheduler struct {
	Timers   map[string]*time.Timer
	RunInfos []RunInfo
}

func NewScheduler(runInfos []RunInfo) Scheduler {
	return Scheduler{
		RunInfos: runInfos,
		Timers:   make(map[string]*time.Timer),
	}
}

func (s *Scheduler) ScheduleAllRuns() {
	fmt.Println(time.Now().Second())
	for _, run := range s.RunInfos {
		s.scheduleRun(run)
	}

}

func (s *Scheduler) scheduleRun(run RunInfo) {
	now := time.Now()

	// nextRunIn, determines the upcoming run for all the executables.
	// only used once after boot, before periodic runs take over.
	var nextRunIn time.Duration
	firstRun := TwelveAmTime(now).Add(run.Offset.Duration)
	if firstRun.Before(now) {
		nextRunIn = run.GetPeriod() - (time.Since(firstRun) % run.GetPeriod())
		if nextRunIn == run.Period.Duration {
			nextRunIn = 0 // rare special case, if next run exactly a run period away, then it must also naturally run now.
		}
	} else {
		nextRunIn = time.Until(firstRun)
	}

	//time.AfterFunc(nextRunIn, func() {
	//	log.Println("running :", run.PackageName)
	//	ticker := time.NewTicker(run.GetPeriod())
	//	for {
	//		select {
	//		case <-ticker.C:
	//			log.Println("rest of : ", run.PackageName)
	//		default:
	//
	//		}
	//	}
	//})

	timer, ok := s.Timers[run.PackageName]
	if ok {
		log.Println("cancelling scheduled run for: ", run.PackageName)
		timer.Stop() // remove already present timer for that packageName
	}

	s.Timers[run.PackageName] = time.AfterFunc(nextRunIn, func() {
		binName := "./" + run.ExecPath
		cmd := exec.Command(binName)
		cmd.Dir = path.Join(executablesPath, run.PackageName)
		log.Println("running package:", run.PackageName, " working dir: ", cmd.Dir, " bin invocation: ", binName)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("failed to run: %s", err)
		}
		log.Printf("output: %s", output)

		s.Timers[run.PackageName].Reset(run.GetPeriod())
	})

}

func (s *Scheduler) runImmediate(packageName string) {
	runInfos, err := jsondisk.Load[[]RunInfo]("data/descriptions.json")
	if err != nil {
		return
	}

	runInitiated := false
	for _, run := range *runInfos {
		if run.PackageName == packageName {
			go func() {
				binName := "./" + run.ExecPath

				cmd := exec.Command(binName)
				cmd.Dir = path.Join(executablesPath, run.PackageName)
				log.Println("running package:", run.PackageName, " working dir: ", cmd.Dir, " bin invocation: ", binName)
				output, err := cmd.Output()
				if err != nil {
					log.Printf("failed to run: %s", err)
				}
				log.Printf("output: %s", output)
			}()
			runInitiated = true
		}
	}

	if runInitiated {
		log.Println("run initiated")
	} else {
		log.Println("that package name was not found")
	}

}
