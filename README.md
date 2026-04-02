# fakebuild

A program that pretends to compile large CMake projects in a colorful terminal.
Create fake productivity with endless realistic-looking CMake build output.

## Features

- Endless build logs (or specify a fixed number of files)
- Random C/C++ compiler warnings (but **never** errors)
- Realistic CMake output format
- ANSI colored output
- Parallel build simulation
- Configurable speed and warning frequency
- Configurable random compilation delay (file size simulation)
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
[  0%] Building C object src/http/client.c.o
[  1%] Building CXX object src/utils/string.cpp.o
In file included from src/utils/string.cpp:12:
src/utils/string.hpp:45: warning: unused variable 'buffer_size' [-Wunused-variable]
   45 |   size_t buffer_size = 1024;
       |            ^~~~~~~~~~
[  2%] Building CXX object src/core/parser.cpp.o
[  3%] Linking CXX executable bin/fakebuild
...
[100%] Built target fakebuild
```

## License

MIT
