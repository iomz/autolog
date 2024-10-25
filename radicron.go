package autolog

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func LogRadicron(filename string) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(
		context.Background(),
		container.ListOptions{All: true},
	)
	if err != nil {
		panic(err)
	}

	var targetID string
	for _, ctr := range containers {
		if strings.Contains(ctr.Image, "radicron") {
			targetID = ctr.ID
			//fmt.Printf("[%s]%s %s (status: %s)\n", ctr.ID, ctr.Names, ctr.Image, ctr.Status)
			break
		}
	}
	if targetID == "" {
		// no container was found
		os.Exit(0)
	}

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	reader, err := apiClient.ContainerLogs(
		context.Background(),
		targetID,
		container.LogsOptions{ShowStderr: true,
			ShowStdout: true,
			Since:      yesterday.Format("2006-01-02T") + "00:00:00",
			Until:      today.Format("2006-01-02T") + "00:00:00",
			//Until:      today.AddDate(0, 0, 1).Format("2006-01-02T") + "00:00:00",
			Timestamps: false,
			Follow:     false,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	reFile := regexp.MustCompile("[0-9]{12}_[A-Z]{3,}_.*$")
	reSave := regexp.MustCompile("file saved")

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	m := make(map[string]string)
	for scanner.Scan() {
		text := scanner.Text()
		if reSave.FindString(text) != "" {
			ss := strings.Split(strings.TrimSuffix(reFile.FindString(text), ".aac"), "_")
			ts, _ := time.Parse("200601021504", ss[0])
			station := ss[1]
			program := ss[2]
			m[ts.Format("2006-01-02T 15:04 ")+station] = program
			// fmt.Println(scanner.Text())
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, k := range slices.Sorted(maps.Keys(m)) {
		if _, err = f.WriteString(fmt.Sprintf("%s %s\n", k, m[k])); err != nil {
			panic(err)
		}
	}
}
