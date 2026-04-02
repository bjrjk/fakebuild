package output

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bjrjk/fakebuild/pkg/config"
)

// Output handles terminal output with conditional coloring
type Output struct {
	config *config.Config
}

// New creates a new Output
func New(cfg *config.Config) *Output {
	return &Output{config: cfg}
}

// color returns a color code if colors are enabled, empty otherwise
func (o *Output) color(color string) string {
	if o.config.NoColor {
		return ""
	}
	return color
}

// PrintCompiling prints the "Building ..." line
func (o *Output) PrintCompiling(progress int, filePath string, isC bool) {
	lang := "CXX"
	if isC {
		lang = "C"
	}
	pct := fmt.Sprintf("%3d%%", progress)
	fmt.Printf("%s[%s]%s Building %s object %s%s%s\n",
		o.color(ColorGray), pct, o.color(ColorReset),
		lang,
		o.color(ColorBold), filePath, o.color(ColorReset))
}

// PrintLinking prints the "Linking ..." line
func (o *Output) PrintLinking(targetName string, isExecutable bool) {
	what := "executable"
	if !isExecutable {
		what = "shared library"
	}
	pct := rand.Intn(99) + 1
	pctStr := fmt.Sprintf("%3d%%", pct)
	fmt.Printf("%s[%s]%s Linking %s %s%s%s\n",
		o.color(ColorGray), pctStr, o.color(ColorReset),
		what,
		o.color(ColorBold), targetName, o.color(ColorReset))
}

// PrintInstalling prints the "Installing ..." line
func (o *Output) PrintInstalling(component string) {
	pct := rand.Intn(99) + 1
	pctStr := fmt.Sprintf("%3d%%", pct)
	fmt.Printf("%s[%s]%s Installing %s%s%s\n",
		o.color(ColorGray), pctStr, o.color(ColorReset),
		o.color(ColorBold), component, o.color(ColorReset))
}

// PrintWarning prints a GCC-style warning
func (o *Output) PrintWarning(file string, line int, message string, option string) {
	// Random chance of showing "In file included from..." line
	if rand.Float32() < 0.3 {
		includeFile := file
		includeLine := rand.Intn(50) + 1
		fmt.Printf("In file included from %s:%d:\n", includeFile, includeLine)
	}

	fmt.Printf("%s%s:%d: %swarning:%s %s [%s%s]\n",
		o.color(ColorBold), file, line,
		o.color(ColorYellow), o.color(ColorReset),
		message, o.color(ColorBold), option)

	// Show the line with caret
	indent := strings.Repeat(" ", len(fmt.Sprintf("%s:%d: ", file, line)))
	var codeLine string
	// The line of "code" with the variable
	words := []string{"int", "size_t", "char", "bool", "static", "const", "auto"}
	word := words[rand.Intn(len(words))]
	varName := message
	if strings.Contains(message, "'") {
		// Extract variable name from message
		start := strings.Index(message, "'")
		end := strings.LastIndex(message, "'")
		if start >= 0 && end > start {
			varName = message[start+1:end]
		}
	}
	if strings.Contains(message, "unused variable") || strings.Contains(message, "unused parameter") {
		codeLine = fmt.Sprintf("   %d | %s %s = %d;\n", line, word, varName, rand.Intn(10000))
	} else if strings.Contains(message, "implicit conversion") {
		codeLine = fmt.Sprintf("   %d | %s = (int)%s;\n", line, varName, varName)
	} else {
		codeLine = fmt.Sprintf("   %d |   %s;\n", line, varName)
	}
	fmt.Printf("%s", indent+codeLine)
	// Add the caret
	caretPos := len(indent) + len(fmt.Sprintf("   %d | ", line)) + len(word) + 1
	fmt.Printf("%s%s^%s\n", indent+strings.Repeat(" ", caretPos), o.color(ColorBold), o.color(ColorReset))
}

// PrintFinished prints the final "Built target" line
func (o *Output) PrintFinished(target string) {
	fmt.Printf("%s[100%%]%s Built target %s%s%s\n",
		o.color(ColorGray), o.color(ColorReset),
		o.color(ColorGreen)+o.color(ColorBold), target, o.color(ColorReset))
}

// RandomDelay sleeps for a random duration based on speed
func (o *Output) RandomDelay() {
	// Random delay between MinDelay and MaxDelay seconds
	delaySec := o.config.MinDelay + rand.Float64()*(o.config.MaxDelay - o.config.MinDelay)
	// Apply speed multiplier
	delaySec = delaySec / o.config.Speed
	delay := time.Duration(delaySec * float64(time.Second))
	time.Sleep(delay)
}

// OutputChan is a message type for the output channel
type OutputChan struct {
	Text string
}

// Printer runs in a goroutine printing from the channel
func (o *Output) Printer(ch <-chan string) {
	for msg := range ch {
		fmt.Print(msg)
	}
}
