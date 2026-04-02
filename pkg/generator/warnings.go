package generator

import (
	"math/rand"
	"strings"
)

// Warning represents a compiler warning
type Warning struct {
	File    string
	Line    int
	Message string
	Option  string
}

var warningTemplates = []struct {
	Message string
	Option  string
}{
	// C/C++ warnings
	{"unused variable '%s' [-Wunused-variable]", "Wunused-variable"},
	{"unused parameter '%s' [-Wunused-parameter]", "Wunused-parameter"},
	{"unused function '%s' [-Wunused-function]", "Wunused-function"},
	{"implicit conversion from 'int' to 'size_t' changes signedness [-Wsign-conversion]", "Wsign-conversion"},
	{"implicit conversion from 'long' to 'int' may change value [-Wimplicit-int-conversion]", "Wimplicit-int-conversion"},
	{"missing field '%s' initializer [-Wmissing-field-initializers]", "Wmissing-field-initializers"},
	{"comparison between signed and unsigned [-Wsign-compare]", "Wsign-compare"},
	{"enumeration value '%s' not handled in switch [-Wswitch]", "Wswitch"},
	{"overflow in implicit constant conversion [-Woverflow]", "Woverflow"},
	{"deprecated declaration of '%s' [-Wdeprecated-declarations]", "Wdeprecated-declarations"},
	{"type qualifier on return type is meaningless [-Wignored-qualifiers]", "Wignored-qualifiers"},
	{"value computed is not used [-Wunused-value]", "Wunused-value"},
	{"statement has no effect [-Wunused-value]", "Wunused-value"},
	{"shifting a negative value is undefined [-Wshift-negative-value]", "Wshift-negative-value"},
	{"the result of '<<' is undefined in this context [-Wshift-overflow]", "Wshift-overflow"},
	{"cast between pointer and integer of different sizes [-Wpointer-integer-cast]", "Wpointer-integer-cast"},
	{"uninitialized variable '%s' [-Wuninitialized]", "Wuninitialized"},
	{"non-static variable '%s' has linkage [-Wnon-static-variable-linkage]", "Wnon-static-variable-linkage"},
	{"function returns a local variable [-Wreturn-local-addr]", "Wreturn-local-addr"},
	{"format string is not a string literal [-Wformat-nonliteral]", "Wformat-nonliteral"},
	{"incompatible pointer types passing '%s' to parameter of type '%s' [-Wincompatible-pointer-types]", "Wincompatible-pointer-types"},
	// Assembly (as) warnings
	{"warning: label '%s' defined but not used", "Wunused-label"},
	{"warning: '%s' label defined multiple times", "Wmultiple-labels"},
	{"warning: ignoring unknown instruction `%s`", "Wunknown-instruction"},
	{"warning: assuming .implicit for section", "Wimplicit-section"},
	{"warning: relocation %s out of range", "Wrelocation-out-of-range"},
	{"warning: changing alignment of section %s", "Wchanged-alignment"},
	{"warning: flag %s is not supported for this target", "Wunsupported-flag"},
	{"warning: found instruction after unconditional jump", "Winstruction-after-jump"},
	{"warning: setting incorrect section attributes", "Wincorrect-section-attrs"},
}

var identifiers = []string{
	"buffer", "size", "count", "index", "tmp", "result", "data", "ptr",
	"value", "temp", "ctx", "context", "state", "flags", "option", "len",
	"capacity", "hash", "key", "node", "entry", "item", "element", "chunk",
	"block", "offset", "addr", "address", "handler", "callback", "func",
	"method", "obj", "object", "instance", "cls", "class", "type", "kind",
	"cache", "cache_size", "buffer_size", "max_size", "min_len", "timeout",
	"interval", "delay", "seconds", "milliseconds", "ret", "code", "err",
	"error", "status", "success", "failure", "flag", "mask", "bits", "byte",
	"word", "dword", "qword", "ptr1", "ptr2", "tmp1", "tmp2", "old_val",
	"new_val", "prev", "next", "current", "head", "tail", "root", "leaf",
}

var types = []string{
	"int", "size_t", "long", "unsigned int", "char", "bool", "float",
	"double", "void*", "char*", "int*", "size_t*", "struct State*",
	"struct Context*", "void (*callback)(void)", "FILE*",
}

// GenerateWarning generates a random compiler warning
func GenerateWarning(file string) *Warning {
	template := warningTemplates[rand.Intn(len(warningTemplates))]

	message := template.Message
	if strings.Contains(message, "%s") {
		// Count how many placeholders
		count := strings.Count(message, "%s")
		for i := 0; i < count; i++ {
			if len(types) > 0 && rand.Float32() < 0.3 {
				message = strings.Replace(message, "%s", types[rand.Intn(len(types))], 1)
			} else {
				message = strings.Replace(message, "%s", identifiers[rand.Intn(len(identifiers))], 1)
			}
		}
	}

	return &Warning{
		File:    file,
		Line:    RandomLineNumber(),
		Message: message,
		Option:  template.Option,
	}
}
