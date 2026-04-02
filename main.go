package main

import (
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bjrjk/fakebuild/pkg/config"
	"github.com/bjrjk/fakebuild/pkg/generator"
	"github.com/bjrjk/fakebuild/pkg/output"
)

// Job represents a compilation job
type Job struct {
	FilePath  string
	IsC       bool
	IsRust    bool
	HasWarning bool
	Warning   *generator.Warning
	HasLink   bool
	LinkTarget string
	LinkIsExe  bool
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	cfg := config.Parse()
	if cfg == nil {
		return
	}

	out := output.New(cfg)

	// Set up signal handling for graceful exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Job channel for worker pool
	jobChan := make(chan Job, cfg.Parallel)

	// WaitGroup for workers
	var wg sync.WaitGroup

	// Start workers - they do the delay (simulate compilation) and output warnings/linking
	worker := func() {
		defer wg.Done()
		for job := range jobChan {
			// Simulate compilation time
			out.RandomDelay()

			// After compilation, output warnings and linking (order doesn't matter)
			if job.HasWarning {
				out.PrintWarning(job.Warning.File, job.Warning.Line, job.Warning.Message, job.Warning.Option)
			}
			if job.HasLink {
				out.PrintLinking(job.LinkTarget, job.LinkIsExe)
			}
		}
	}

	// Start workers
	wg.Add(cfg.Parallel)
	for i := 0; i < cfg.Parallel; i++ {
		go worker()
	}

	// Main goroutine generates ALL "Building" lines in order - GUARANTEES increasing progress
	// This is the key fix: we print Building... immediately in order before sending to worker
	if cfg.Endless {
		current := 0
		for {
			select {
			case <-sigChan:
				goto exit
			default:
			}

			progress := (current * 100) / 1000
			if progress > 99 {
				progress = 99
			}

			job := generateJob(cfg)
			// Print Building IMMEDIATELY in correct order
			out.PrintCompiling(progress, job.FilePath+".o", job.IsC, job.IsRust)
			// Send to worker for delay and warnings
			jobChan <- job
			current++
		}
	} else {
		// Finite mode: print Building lines in order
		for current := 0; current < cfg.TotalFiles; current++ {
			select {
			case <-sigChan:
				goto exit
			default:
			}

			progress := (current * 100) / cfg.TotalFiles
			if progress > 100 {
				progress = 100
			}

			job := generateJob(cfg)
			// Print Building IMMEDIATELY in correct order
			out.PrintCompiling(progress, job.FilePath+".o", job.IsC, job.IsRust)
			// Send to worker for delay and warnings
			jobChan <- job
		}
		close(jobChan)

		// Wait for all workers to finish
		go func() {
			wg.Wait()
			sigChan <- os.Interrupt // trigger exit
		}()

		// Wait for completion
		<-sigChan
		// After all workers done, print finished
		out.PrintFinished("fakebuild")
		goto exit
	}

exit:
	os.Exit(0)
}

func generateJob(cfg *config.Config) Job {
	filePath := generator.RandomFilePath()
	isC := generator.IsC(filePath)
	isRust := generator.IsRust(filePath)

	job := Job{
		FilePath: filePath,
		IsC:      isC,
		IsRust:   isRust,
	}

	if rand.Float64() < cfg.WarningFreq {
		job.HasWarning = true
		job.Warning = generator.GenerateWarning(filePath)
	}

	if rand.Float32() < 0.05 {
		job.HasLink = true
		job.LinkTarget = generator.RandomTargetName()
		job.LinkIsExe = rand.Float32() < 0.7
	}

	return job
}
