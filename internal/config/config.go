package config

import (
	"flag"
	"os"
	"runtime"
	"strings"
)

var (
	sizeFlag = flag.Int("size", 0, "Number of chanks in every axis. Chunk size is 16 voxels.")
	cpuFlag  = flag.Int("cpu", 0, "Number of parallel workers to call lua")
)

type Config struct {
	Script  string
	Size    int
	Workers int
}

func ParseConfig() Config {
	permutateArgs(os.Args)
	flag.Parse()
	s := size()
	return Config{
		Script:  script(),
		Workers: cpu(s),
		Size:    s,
	}
}

func script() string {
	args := flag.Args()
	script := "demo"
	if len(args) >= 1 {
		arg := args[0]
		if strings.HasSuffix(arg, ".lua") {
			return arg
		}
		script = strings.TrimSuffix(arg, ".lua")
	}
	return "demo/" + script + ".lua"
}

func cpu(size int) int {
	if *cpuFlag > 0 {
		return *cpuFlag
	}
	cpu := runtime.NumCPU()
	if chunks := size * size * size; cpu >= chunks {
		return chunks
	}
	if cpu == 0 {
		return 1
	}
	return cpu
}

func size() int {
	if *sizeFlag > 0 {
		return *sizeFlag
	}
	return 2
}

// permutateArgs permutates args such that options are in front, leaving the program name untouched.
// Thanks to fuz in https://stackoverflow.com/questions/25113313
func permutateArgs(args []string) int {
	args = args[1:]
	optind := 0

	for i := range args {
		if args[i][0] == '-' {
			tmp := args[i]
			args[i] = args[optind]
			args[optind] = tmp
			optind++
		}
	}

	return optind + 1
}
