# fakebuild

A program that pretends to compile large CMake projects in a colorful terminal.
Create fake productivity with endless realistic-looking CMake build output.

## Features

- Endless build logs (or specify a fixed number of files)
- Random C/C++/Rust compiler warnings (but **never** errors)
- Realistic CMake output format with **strictly increasing percentage**
- Supports C, C++, and **Rust** files
- ANSI colored output
- Parallel build simulation
- Configurable speed and warning frequency
- Configurable random compilation delay (simulates varying file sizes)
- Intermediate and final linking steps like real CMake
- 100% fake productivity

## Building

Requires Go 1.13+

```bash
make build
# or directly
go build -o fakebuild .
```

## Installing

```bash
# installs to /usr/local/bin/fakebuild
make install
```

## Usage

```
fakebuild [options]

Options:
  -s, --speed FLOAT       Speed multiplier (default: 1.0)
  -p, --parallel INT       Number of parallel jobs (default: number of CPUs)
  -e, --endless            Run forever (default: true)
  -t, --total INT          Total files to compile (0 = endless, default: 0)
  -w, --warnings FLOAT     Warning frequency 0.0 - 1.0 (default: 0.15)
  -m, --min-delay FLOAT    Minimum compilation delay in seconds (default: 0)
  -M, --max-delay FLOAT    Maximum compilation delay in seconds (default: 10)
  --no-color               Disable ANSI colored output
  -h, --help               Show this help message
```

### Examples

```bash
# Default: endless build, full fake productivity (0-10s random delay)
fakebuild

# 16 parallel jobs, realistic 1-15 second random delays
fakebuild --parallel 16 --min-delay 1 --max-delay 15

# 16 parallel jobs, double speed, 30% warning rate
fakebuild --parallel 16 --speed 2 --warnings 0.3

# Compile 1000 files then exit
fakebuild --total 1000
```

## Example Output

```
[  0%] Building Rust object thirdparty/zlib/hash_memory.rs.o
[ 10%] Building CXX object src/lexer/vector_socket.cc.o
[ 20%] Building CXX object modules/video/tar_decode.cc.o
[ 30%] Building CXX object thirdparty/match_math.cxx.o
[ 40%] Building Rust object examples/frame.rs.o
In file included from src/lexer/vector_socket.cc:25:
src/lexer/vector_socket.cc:41: warning: deprecated declaration of 'cls' [-Wdeprecated-declarations]
                                  41 |   cls;
                                           ^
[ 20%] Linking executable libdb
[ 26%] Building CXX object src/compiler/search.cpp.o
...
[ 99%] Linking executable fakebuild
[100%] Built target fakebuild
```

## License

MIT
