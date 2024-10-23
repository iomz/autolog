package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/iomz/autolog"
	"github.com/spf13/viper"
)

func main() {
	conf := flag.String("config", "auto-pom.toml", "The auto-pom.[toml|yml] defining the config.")
	version := flag.Bool("version", false, "Print version.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		bi, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", bi.Main.Version)
		os.Exit(0)
	}

	// load config
	if *conf != "auto-pom.toml" {
		configPath, err := filepath.Abs(*conf)
		if err != nil {
			panic(err)
		}
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("auto-pom")
		viper.AddConfigPath(".")
		// add the path to the default config
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("no caller information")
		}
		viper.AddConfigPath(filepath.Join(filepath.Dir(filename), "../../"))
	}

	// read the config file
	if err := viper.ReadInConfig(); err != nil { // handle errors reading the config file
		log.Fatalf("config: %s \n", err)
	}

	log.Printf("Location: %s\n", autolog.Location)

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

	for _, ctr := range containers {
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
}
