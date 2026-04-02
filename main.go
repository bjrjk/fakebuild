package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bjrjk/fakebuild/pkg/config"
	"github.com/bjrjk/fakebuild/pkg/generator"
	"github.com/bjrjk/fakebuild/pkg/output"
)

// Job represents a compilation job to be processed by a worker
type Job struct {
	Current   int
	Progress  int
	FilePath  string
	IsC       bool
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

	// Output channel - all output goes through here to guarantee order
	// This ensures "Building" lines are printed in increasing progress order
	outputChan := make(chan string, 64)

	// WaitGroup for waiting on workers
	var wg sync.WaitGroup

	// Done channel signals when all work is complete
	done := make(chan struct{}, 1)

	// Single-threaded output printer - guarantees order of output
	go func() {
		for msg := range outputChan {
			fmt.Print(msg)
		}
	}()

	// Capture output into the channel instead of printing directly
	capturePrint := func(format string, args ...interface{}) {
		outputChan <- fmt.Sprintf(format, args...)
	}

	// Worker function - processes jobs
	worker := func() {
		defer wg.Done()
		for job := range jobChan {
			// Simulate compilation by sleeping after "Building" has already been printed
			out.RandomDelay()

			// After compilation completes, send warnings and linking to output
			if job.HasWarning {
				// We need to capture the warning output to send through the channel
				// But for simplicity, we'll generate it here and send
				if rand.Float32() < 0.3 {
					capturePrint("In file included from %s:%d:\n", job.Warning.File, rand.Intn(50)+1)
				}
				capturePrint("%s%s:%d: %swarning:%s %s [%s%s]\n",
					output.ColorBold, job.Warning.File, job.Warning.Line,
					output.ColorYellow, output.ColorReset,
					job.Warning.Message, output.ColorBold, job.Warning.Option)

				// Re-generate code line and caret (random is fine)
				indent := strings.Repeat(" ", len(fmt.Sprintf("%s:%d: ", job.Warning.File, job.Warning.Line)))
				words := []string{"int", "size_t", "char", "bool", "static", "const", "auto"}
				word := words[rand.Intn(len(words))]
				varName := job.Warning.Message
				if strings.Contains(job.Warning.Message, "'") {
					start := strings.Index(job.Warning.Message, "'")
					end := strings.LastIndex(job.Warning.Message, "'")
					if start >= 0 && end > start {
						varName = job.Warning.Message[start+1:end]
					}
				}
				var codeLine string
				line := job.Warning.Line
				if strings.Contains(job.Warning.Message, "unused variable") || strings.Contains(job.Warning.Message, "unused parameter") {
					codeLine = fmt.Sprintf("   %d | %s %s = %d;\n", line, word, varName, rand.Intn(10000))
				} else if strings.Contains(job.Warning.Message, "implicit conversion") {
					codeLine = fmt.Sprintf("   %d | %s = (int)%s;\n", line, varName, varName)
				} else {
					codeLine = fmt.Sprintf("   %d |   %s;\n", line, varName)
				}
				capturePrint("%s", indent+codeLine)
				caretPos := len(indent) + len(fmt.Sprintf("   %d | ", line)) + len(word) + 1
				capturePrint("%s%s^%s\n", indent+strings.Repeat(" ", caretPos), output.ColorBold, output.ColorReset)
			}

			if job.HasLink {
				pct := rand.Intn(99) + 1
				what := "executable"
				if !job.LinkIsExe {
					what = "shared library"
				}
				capturePrint("%s[%3d%%]%s Linking %s %s%s%s\n",
					output.ColorGray, pct, output.ColorReset,
					what,
					output.ColorBold, job.LinkTarget, output.ColorReset)
			}
		}
	}

	// Start workers
	wg.Add(cfg.Parallel)
	for i := 0; i < cfg.Parallel; i++ {
		go worker()
	}

	// Generate jobs in order - this guarantees "Building" lines are output in increasing order
	go func() {
		if cfg.Endless {
			// Endless mode
			current := 0
			for {
				progress := (current * 100) / 1000
				if progress > 99 {
					progress = 99
				}
				generateAndSendJob(current, progress, cfg, out, capturePrint, jobChan)
				current++
			}
		} else {
			// Finite mode
			for current := 0; current < cfg.TotalFiles; current++ {
				progress := (current * 100) / cfg.TotalFiles
				if progress > 100 {
					progress = 100
				}
				generateAndSendJob(current, progress, cfg, out, capturePrint, jobChan)
			}
			close(jobChan)
		}
	}()

	// Wait for completion
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
		capturePrint("%s[100%%]%s Built target %s%sfakebuild%s\n",
			output.ColorGray, output.ColorReset,
			output.ColorGreen+output.ColorBold, output.ColorReset, "")
	}

	// Close output channel after we're done with all output
	close(outputChan)

	// Exit
	os.Exit(0)
}

// generateAndSendJob generates a new job, prints the "Building" line immediately (in order), then sends to worker
func generateAndSendJob(current int, progress int, cfg *config.Config, out *output.Output, capturePrint func(format string, args ...interface{}), jobChan chan<- Job) {
	filePath := generator.RandomFilePath()
	isC := generator.IsC(filePath)

	// Print "Building..." IMMEDIATELY, in order - this is what guarantees increasing progress
	lang := "CXX"
	if isC {
		lang = "C"
	}
	pct := fmt.Sprintf("%3d%%", progress)
	if cfg.NoColor {
		capturePrint("[%s] Building %s object %s%s\n", pct, lang, filePath, ".o")
	} else {
		capturePrint("%s[%s]%s Building %s object %s%s%s%s\n",
			output.ColorGray, pct, output.ColorReset,
			lang,
			output.ColorBold, filePath, ".o", output.ColorReset)
	}

	// Prepare job for worker
	job := Job{
		Current:  current,
		Progress: progress,
		FilePath: filePath,
		IsC:      isC,
	}

	// Chance of warning after compilation
	if rand.Float64() < cfg.WarningFreq {
		job.HasWarning = true
		job.Warning = generator.GenerateWarning(filePath)
	}

	// Chance of linking after compilation
	if rand.Float32() < 0.05 {
		job.HasLink = true
		job.LinkTarget = generator.RandomTargetName()
		job.LinkIsExe = rand.Float32() < 0.7
	}

	// Send job to worker for processing (sleep)
	jobChan <- job
}
