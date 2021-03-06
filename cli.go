package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error reading command pipe: %s\n", err)

		if os.Getenv("TESTING") == "" {
			os.Exit(1)
		}
		return
	}

	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage:")
		fmt.Println("  cat yourlogfile.log | zapper")
	} else if info.Size() > 0 {
		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(os.Stdout)

		zap := zapper{reader, writer}
		zap.pipe()

		if os.Getenv("TESTING") == "" {
			os.Exit(0)
		}
		return
	}
}
