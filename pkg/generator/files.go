package generator

import (
	"math/rand"
)

var (
	directories = []string{
		"src", "src/core", "src/utils", "src/io", "src/net", "src/http", "src/db",
		"src/parser", "src/lexer", "src/compiler", "src/runtime", "src/gui",
		"src/platform", "src/thirdparty", "thirdparty", "thirdparty/zlib",
		"thirdparty/curl", "thirdparty/openssl", "thirdparty/libpng",
		"thirdparty/freetype", "thirdparty/glfw", "modules", "modules/audio",
		"modules/video", "modules/ui", "modules/renderer", "tests",
		"examples", "apps", "apps/cli", "apps/gui",
	}

	baseNames = []string{
		"parser", "lexer", "token", "config", "utils", "string", "buffer",
		"http", "client", "server", "request", "response", "json", "xml",
		"database", "connection", "query", "result", "network", "socket",
		"tcp", "udp", "file", "fs", "path", "directory", "thread", "mutex",
		"memory", "allocator", "cache", "hash", "map", "vector", "list",
		"tree", "graph", "algorithm", "sort", "search", "math", "matrix",
		"vector", "quaternion", "transform", "camera", "light", "material",
		"shader", "texture", "mesh", "model", "scene", "window", "input",
		"keyboard", "mouse", "event", "logger", "log", "debug", "error",
		"exception", "timer", "time", "clock", "aes", "crypto", "hash",
		"md5", "sha256", "random", "uuid", "color", "image", "png", "jpeg",
		"font", "text", "render", "draw", "shape", "rectangle", "circle",
		"line", "point", "audio", "sound", "music", "player", "wave", "mp3",
		"video", "frame", "encode", "decode", "stream", "compression", "zip",
		"gzip", "deflate", "inflate", "archive", "tar", "regex", "pattern",
		"match", "argparse", "cli", "command", "option", "argument",
	}

	extensions = []string{".c", ".cpp", ".cxx", ".cc", ".rs"}
	targetNames = []string{
		"main", "fakebuild", "app", "server", "client", "tool", "cli", "gui",
		"libcore", "libutils", "libnet", "libhttp", "libdb", "libparser",
		"libui", "librender", "libaudio", "libvideo", "testapp", "demo",
	}
)

// RandomFilePath generates a random fake C/C++ file path
func RandomFilePath() string {
	dir := directories[rand.Intn(len(directories))]
	base := randomBaseName()
	ext := extensions[rand.Intn(len(extensions))]
	return dir + "/" + base + ext
}

// RandomBaseName generates a random base name (possibly compound)
func randomBaseName() string {
	if rand.Float32() < 0.3 {
		// Compound name like http_client
		return baseNames[rand.Intn(len(baseNames))] + "_" + baseNames[rand.Intn(len(baseNames))]
	}
	return baseNames[rand.Intn(len(baseNames))]
}

// IsC returns whether the file should be treated as C
func IsC(filePath string) bool {
	return len(filePath) >= 2 && filePath[len(filePath)-2:] == ".c"
}

// IsRust returns whether the file should be treated as Rust
func IsRust(filePath string) bool {
	return len(filePath) >= 3 && filePath[len(filePath)-3:] == ".rs"
}

// RandomTargetName generates a random target name
func RandomTargetName() string {
	return targetNames[rand.Intn(len(targetNames))]
}

// RandomLineNumber generates a random line number (1-1000)
func RandomLineNumber() int {
	return 1 + rand.Intn(999)
}
