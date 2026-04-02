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
// All percentage-based output is printed by main goroutine to guarantee order
type Job struct {
	FilePath   string
	IsC        bool
	IsRust     bool
	IsAssembly bool
	HasWarning bool
	Warning    *generator.Warning
	// Link info - only need it for after-compilation output, percentage already printed
	HasLink    bool
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

	// Start workers - they do the delay (simulate compilation) and output warnings
	// All percentage-based output (Building/Linking) is already printed by main goroutine to guarantee order
	worker := func() {
		defer wg.Done()
		for job := range jobChan {
			// Simulate compilation time (delay after percentage output)
			out.RandomDelay()

			// After compilation completes, output warnings (warnings go after compilation anyway, order doesn't matter)
			if job.HasWarning {
				out.PrintWarning(job.Warning.File, job.Warning.Line, job.Warning.Message, job.Warning.Option)
			}
			// Linking percentage is already printed by main goroutine in correct order
			// No need to output again here
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

			job, linkProgress := generateJob(cfg, current, 0)
			// Print Building IMMEDIATELY in correct order (guarantees increasing percentage)
			out.PrintCompiling(progress, job.FilePath+".o", job.IsC, job.IsRust, job.IsAssembly)
			// If we need to do a link after this compilation, print it in order too
			if job.HasLink {
				out.PrintLinkingWithProgress(linkProgress, job.LinkTarget, job.LinkIsExe)
			}
			// Send to worker for delay and warnings
			jobChan <- job
			if job.HasLink {
				current += 2
			} else {
				current++
			}
		}
	} else {
		// Finite mode: print everything in order
		current := 0
		// Reserve the last step for final linking (99%), so we need one less compilation
		for current < cfg.TotalFiles - 1 {
			select {
			case <-sigChan:
				goto exit
			default:
			}

			progress := (current * 100) / cfg.TotalFiles
			if progress > 100 {
				progress = 100
			}

			job, linkProgress := generateJob(cfg, current, cfg.TotalFiles)
			// Print Building IMMEDIATELY in correct order (guarantees increasing percentage)
			out.PrintCompiling(progress, job.FilePath+".o", job.IsC, job.IsRust, job.IsAssembly)
			// If we need to do a link after this compilation, print it in order too
			if job.HasLink {
				out.PrintLinkingWithProgress(linkProgress, job.LinkTarget, job.LinkIsExe)
			}
			// Send to worker for delay and warnings
			jobChan <- job
			if job.HasLink {
				current += 2
			} else {
				current++
			}
		}

		// Print final linking step at 99% before 100% finished
		out.PrintLinkingWithProgress(99, "fakebuild", true)
		finalJob := Job{HasLink: true, LinkTarget: "fakebuild", LinkIsExe: true}
		jobChan <- finalJob
		close(jobChan)

		// Wait for all workers to finish including final linking delay
		wg.Wait()

		// After all workers done, print finished
		out.PrintFinished("fakebuild")
		goto exit
	}

exit:
	os.Exit(0)
}

// generateJob returns a new job with random content
// Link info is stored, but percentage is printed by main for order
func generateJob(cfg *config.Config, current int, total int) (Job, int) {
	filePath := generator.RandomFilePath()
	isC := generator.IsC(filePath)
	isRust := generator.IsRust(filePath)
	isAssembly := generator.IsAssembly(filePath)

	job := Job{
		FilePath:    filePath,
		IsC:         isC,
		IsRust:      isRust,
		IsAssembly:  isAssembly,
	}

	if rand.Float64() < cfg.WarningFreq {
		job.HasWarning = true
		job.Warning = generator.GenerateWarning(filePath)
	}

	var linkProgress int
	if rand.Float32() < 0.05 {
		job.HasLink = true
		job.LinkTarget = generator.RandomTargetName()
		job.LinkIsExe = rand.Float32() < 0.7
		current++
		if cfg.Endless {
			linkProgress = (current * 100) / 1000
			if linkProgress > 99 {
				linkProgress = 99
			}
		} else {
			linkProgress = (current * 100) / total
			if linkProgress > 100 {
				linkProgress = 100
			}
		}
	}

	return job, linkProgress
}
