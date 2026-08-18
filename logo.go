package main

import (
	"fmt"
	"strings"
)

var logo = [1]string{`                                 ++++++++               +++++
                              ++++++++++++++       ++++++++++++++
                           ++++++++++++++++++   +++++++++++++++++++
            ++++++++++++  ++++++++++++++++++++++++++++++++++++++++++
                        +++++++++      +++++  ++++++++       ++++++++
       ++++++++++++++++ +++++++    +++++++++++++++++          +++++++
               +++++++ ++++++++   ++++++++++++++++++          +++++++
                       ++++++++  +++++++++++++++++++          +++++++
                        +++++++        ++++++++++++++       ++++++++
                        ++++++++++++++++++++++++++++++++ ++++++++++
                         +++++++++++++++++++  ++++++++++++++++++++
                           +++++++++++++++      ++++++++++++++++
                             ++++++++++            ++++++++++`}

const (
	colorPurple = "\033[35m"
	colorReset  = "\033[0m"
)

func PrintLogo(infoLines []string) {
	logoLines := strings.Split(logo[0], "\n")

	for i, line := range logoLines {
		info := ""
		if i < len(infoLines) {
			info = infoLines[i]
		}
		padding := 78 - len(line)
		if padding < 0 {
			padding = 0
		}
		fmt.Printf("%s%s%s%s%s\n", colorPurple, line, colorReset, strings.Repeat(" ", padding), info)
	}
}
