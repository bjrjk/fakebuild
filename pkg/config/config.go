package config

import (
	"flag"
	"fmt"
	"runtime"
)

// Config holds all configuration options for fakebuild
type Config struct {
	Speed       float64  // Speed multiplier (1.0 = normal)
	Parallel    int      // Number of parallel compilation jobs
	Endless     bool     // Run forever
	TotalFiles  int      // Total files to compile (0 = endless)
	WarningFreq float64  // Probability of a warning (0.0 - 1.0)
	MinDelay    float64  // Minimum compilation delay in seconds (default: 0)
	MaxDelay    float64  // Maximum compilation delay in seconds (default: 10)
	NoColor     bool     // Disable ANSI colors
}

// Parse parses command-line arguments and returns a Config
func Parse() *Config {
	cfg := &Config{
		Speed:       1.0,
		Parallel:    runtime.NumCPU(),
		Endless:     true,
		TotalFiles:  0,
		WarningFreq: 0.15,
		MinDelay:    0.0,
		MaxDelay:    10.0,
		NoColor:     false,
	}

	flag.Float64Var(&cfg.Speed, "speed", cfg.Speed, "Speed multiplier (default: 1.0)")
	flag.Float64Var(&cfg.Speed, "s", cfg.Speed, "Speed multiplier (short)")

	flag.IntVar(&cfg.Parallel, "parallel", cfg.Parallel, "Number of parallel jobs (default: number of CPUs)")
	flag.IntVar(&cfg.Parallel, "p", cfg.Parallel, "Number of parallel jobs (short)")

	flag.BoolVar(&cfg.Endless, "endless", cfg.Endless, "Run forever (default: true)")
	flag.BoolVar(&cfg.Endless, "e", cfg.Endless, "Run forever (short)")

	flag.IntVar(&cfg.TotalFiles, "total", cfg.TotalFiles, "Total files to compile (0 = endless)")
	flag.IntVar(&cfg.TotalFiles, "t", cfg.TotalFiles, "Total files to compile (short)")

	flag.Float64Var(&cfg.WarningFreq, "warnings", cfg.WarningFreq, "Warning frequency (0.0 - 1.0)")
	flag.Float64Var(&cfg.WarningFreq, "w", cfg.WarningFreq, "Warning frequency (short)")

	flag.Float64Var(&cfg.MinDelay, "min-delay", cfg.MinDelay, "Minimum compilation delay in seconds (default: 0)")
	flag.Float64Var(&cfg.MinDelay, "m", cfg.MinDelay, "Minimum compilation delay in seconds (short)")

	flag.Float64Var(&cfg.MaxDelay, "max-delay", cfg.MaxDelay, "Maximum compilation delay in seconds (default: 10)")
	flag.Float64Var(&cfg.MaxDelay, "M", cfg.MaxDelay, "Maximum compilation delay in seconds (short)")

	flag.BoolVar(&cfg.NoColor, "no-color", cfg.NoColor, "Disable ANSI colored output")

	help := flag.Bool("help", false, "Show this help message")
	flag.Bool("h", false, "Show this help message (short)")

	flag.Parse()

	if *help {
		printHelp()
		return nil
	}

	// If TotalFiles > 0, disable endless mode
	if cfg.TotalFiles > 0 {
		cfg.Endless = false
	}

	// Validate ranges
	if cfg.Speed <= 0 {
		cfg.Speed = 1.0
	}
	if cfg.Parallel <= 0 {
		cfg.Parallel = 1
	}
	if cfg.WarningFreq < 0 {
		cfg.WarningFreq = 0
	}
	if cfg.WarningFreq > 1 {
		cfg.WarningFreq = 1
	}
	if cfg.MinDelay < 0 {
		cfg.MinDelay = 0
	}
	if cfg.MaxDelay < cfg.MinDelay {
		cfg.MaxDelay = cfg.MinDelay
	}

	return cfg
}

func printHelp() {
	fmt.Println("fakebuild - pretend to compile large CMake projects for fake productivity")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  fakebuild [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -s, --speed FLOAT       Speed multiplier (default: 1.0)")
	fmt.Println("  -p, --parallel INT       Number of parallel jobs (default: number of CPUs)")
	fmt.Println("  -e, --endless            Run forever (default: true)")
	fmt.Println("  -t, --total INT          Total files to compile (0 = endless, default: 0)")
	fmt.Println("  -w, --warnings FLOAT     Warning frequency 0.0 - 1.0 (default: 0.15)")
	fmt.Println("  -m, --min-delay FLOAT    Minimum compilation delay in seconds (default: 0)")
	fmt.Println("  -M, --max-delay FLOAT    Maximum compilation delay in seconds (default: 10)")
	fmt.Println("  --no-color               Disable ANSI colored output")
	fmt.Println("  -h, --help               Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  fakebuild")
	fmt.Println("  fakebuild --parallel 16 --min-delay 1 --max-delay 15")
	fmt.Println("  fakebuild --speed 2 --warnings 0.3")
	fmt.Println("  fakebuild --total 1000")
}
