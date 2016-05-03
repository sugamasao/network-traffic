package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type NetworkInformation struct {
	InBytes  int
	OutBytes int
}

func main() {
	ifName := ""
	if len(os.Args) > 1 {
		ifName = os.Args[1]
	}

	baseStats := getNetworkStat(ifName)
	for {
		nowStats := getNetworkStat(ifName)
		for k, v := range nowStats {
			in_diff := v.InBytes - baseStats[k].InBytes
			out_diff := v.OutBytes - baseStats[k].OutBytes
			fmt.Println(k, in_diff, out_diff)
		}
		time.Sleep(time.Second)
		baseStats = nowStats
	}
}

func execCommand(ifName string) string {
	command := "netstat"
	options := []string{"-bni"}
	if ifName != "" {
		options = append(options, "-I", ifName)
	}

	out, err := exec.Command(command, options...).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(out)
}

func getNetworkStat(ifName string) map[string]NetworkInformation {
	result := make(map[string]NetworkInformation)
	rowStat := execCommand(ifName)
	rowStatLine := strings.Split(rowStat, "\n")

	for i := range rowStatLine {
		if i == 0 { // header.
			continue
		}
		if rowStatLine[i] == "" { // last line.
			continue
		}

		fields := strings.Fields(string(rowStatLine[i]))
		if regexp.MustCompile(`^<Link#\d+>$`).Match([]byte(fields[2])) {
			in, _ := strconv.Atoi(fields[6])
			out, _ := strconv.Atoi(fields[9])
			result[fields[0]] = NetworkInformation{InBytes: in, OutBytes: out}
			continue
		}
	}

	return result
}
