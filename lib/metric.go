package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/marguerite/go-stdlib/ioutils"
)

func mkMetricCompatibility(f io.Reader) string {
	var metric string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		metric += line + "\n"
		if strings.Contains(line, "<alias ") {
			metric += "\t  <test name=\"search_metric_aliases\"><bool>true</bool></test>\n"
		}
		if strings.Contains(line, "<!DOCTYPE ") {
			metric += "\n<!-- DO NOT EDIT; this is a generated file -->\n<!-- modify /etc/sysconfig/fonts-config && run /usr/bin/fonts-config instead -->\n\n"
		}
	}

	return metric
}

// GenMetricCompatibility generate 30-metric-aliases.conf
func GenMetricCompatibility(verbosity int) {
	// replace fontconfig's /etc/fonts/conf.d/30-metric-aliases.conf
	// by fonts-config's one

	avail := "/usr/share/fontconfig/conf.avail/30-metric-aliases.conf"
	file := "/etc/fonts/conf.d/30-metric-aliases.conf"

	text := mkMetricCompatibility(ioutils.NewReaderFromFile(avail))

	Dbg(verbosity, Debug, fmt.Sprintf("Writing %s\n", file))

	err := overwriteOrRemoveFile(file, []byte(text))
	if err != nil {
		log.Fatalf("Can not write %s: %s\n", file, err.Error())
	}
}
