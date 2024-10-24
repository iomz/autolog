package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
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
			fmt.Printf("[%s]%s %s (status: %s)\n", ctr.ID, ctr.Names, ctr.Image, ctr.Status)
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
			//Until:      today.Format("2006-01-02T") + "00:00:00",
			Until:      today.AddDate(0, 0, 1).Format("2006-01-02T") + "00:00:00",
			Timestamps: false,
			Follow:     false,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	result := []string{}
	re, err := regexp.Compile("^.*file saved.*$")
	for scanner.Scan() {
		s := re.FindString(scanner.Text())
		result = append(result, s)
		fmt.Println(s)
	}
	_ = result
}
