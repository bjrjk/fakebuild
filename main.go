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

	// Progress tracking with mutex for thread-safe access
	var totalCompleted int
	var mu sync.Mutex

	// WaitGroup for waiting on workers
	var wg sync.WaitGroup

	// Done channel signals when all work is complete
	done := make(chan struct{}, 1)

	// Worker function for parallel jobs
	worker := func() {
		defer wg.Done()
		for {
			// Check if we should stop (with lock)
			mu.Lock()
			if !cfg.Endless && totalCompleted >= cfg.TotalFiles {
				mu.Unlock()
				return
			}
			current := totalCompleted
			totalCompleted++
			mu.Unlock()

			// Generate file
			filePath := generator.RandomFilePath()
			isC := generator.IsC(filePath)

			// Calculate progress
			var progress int
			if cfg.Endless {
				progress = (current * 100) / 1000 // Fake progress up to 99%
				if progress > 99 {
					progress = 99
				}
			} else {
				if cfg.TotalFiles > 0 {
					progress = (current * 100) / cfg.TotalFiles
					if progress > 100 {
						progress = 100
					}
				} else {
					progress = 0
				}
			}

			// Random delay to simulate compilation
			out.RandomDelay()

			out.PrintCompiling(progress, filePath+".o", isC)

			// Chance of linking after some files
			if rand.Float32() < 0.05 {
				out.RandomDelay()
				target := generator.RandomTargetName()
				out.PrintLinking(target, rand.Float32() < 0.7)
			}

			// Chance of warning
			if rand.Float64() < cfg.WarningFreq {
				warning := generator.GenerateWarning(filePath)
				out.PrintWarning(warning.File, warning.Line, warning.Message, warning.Option)
			}
		}
	}

	// Start workers
	wg.Add(cfg.Parallel)
	for i := 0; i < cfg.Parallel; i++ {
		go worker()
	}

	// If not endless, wait for all workers to finish then signal done
	if !cfg.Endless {
		go func() {
			wg.Wait()
			done <- struct{}{}
		}()
	}

	// Wait for signal or completion
	select {
	case <-sigChan:
	case <-done:
	}

	if !cfg.Endless {
		// Finished normally
		out.PrintFinished("fakebuild")
	}

	// Exit
	os.Exit(0)
}
