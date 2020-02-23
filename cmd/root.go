package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// Used for flags.
var (
	readAddress int64
	size int64
	hexdump bool
	memorymap bool
	quiet bool
	pid string
	permission string
	outputPath string
	targetProcess string
)

var rootCmd = &cobra.Command{
	Use:   "dumproid",
	Short: "Android memory dump tool without ndk",
	Run: func(cmd *cobra.Command, args []string) {
		if !quiet {
			asciiArt()
		}
		
		if pid == "" {
			parsedPid, err := getPID()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			pid = parsedPid
		}

		if memorymap {
			displayMemoryMap(pid)
			os.Exit(0)
		}

		if hexdump {
			if readAddress != 0 {
				displayMemoryBytes(pid, readAddress, size)
				os.Exit(0)
			}	
		}

		dumpToFile(pid, permission, targetProcess, outputPath)
	},
}

func init() {
	rootCmd.Flags().Int64VarP(&readAddress, "address", "a", 0, "begin address of target memory")
	rootCmd.Flags().Int64VarP(&size, "number", "n", 0x10, "number of bytes")
	rootCmd.Flags().BoolVarP(&memorymap, "maps", "m", false, "output memory mapping")
	rootCmd.Flags().BoolVarP(&hexdump, "dump", "d", false, "output hexdump memory")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "output hexdump memory")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "/data/local/tmp", "output path")
	rootCmd.Flags().StringVarP(&permission, "filter", "f", "", "")
	rootCmd.Flags().StringVarP(&pid, "pid", "p", "", "target app's pid")
	rootCmd.Flags().StringVarP(&targetProcess, "target", "t", "", "target process name")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func asciiArt() {
	art := `
██████╗ ██╗   ██╗███╗   ███╗██████╗ ██████╗  ██████╗ ██╗██████╗
██╔══██╗██║   ██║████╗ ████║██╔══██╗██╔══██╗██╔═══██╗██║██╔══██╗
██║  ██║██║   ██║██╔████╔██║██████╔╝██████╔╝██║   ██║██║██║  ██║
██║  ██║██║   ██║██║╚██╔╝██║██╔═══╝ ██╔══██╗██║   ██║██║██║  ██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║     ██║  ██║╚██████╔╝██║██████╔╝
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝╚═════╝
`
	fmt.Println(art) 
}

func getPID() (string, error) {
	cmd := exec.Command("ps", "-e")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	re := regexp.MustCompile(`\s+`)
	line, err := out.ReadString('\n')
	pids := []string{}
	for err == nil && len(line) != 0 {
		s := strings.Split(re.ReplaceAllString(string(line), " "), " ")
		pid := s[1]
		cmd := s[8]
		if pid != "PID" && cmd != "" && cmd != "ps" && cmd != "sh" && cmd != "dumproid" {
			pids = append(pids, pid)
		}
		line, err = out.ReadString('\n')
	}

	if len(pids) == 1 {
		fmt.Printf("Target PID: %s\n", pids[0])
		return pids[0], nil
	}

	return "", fmt.Errorf("Failed to identify PID")
}
