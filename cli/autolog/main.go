package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/iomz/autolog"
	"github.com/spf13/viper"
)

func main() {
	conf := flag.String("config", "autolog.toml", "The autolog.[toml|yml] defining the config.")
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
	if *conf != "autolog.toml" {
		configPath, err := filepath.Abs(*conf)
		if err != nil {
			panic(err)
		}
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("autolog")
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

	directory := viper.GetString("directory")

	// radicron
	radicronFile := filepath.Join(directory, "radicron.log")
	autolog.LogRadicron(radicronFile)

	// Git
	r, err := git.PlainOpen(directory)
	if err != nil {
		panic(err)
	}
	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}
	_, err = w.Add(".")
	if err != nil {
		panic(err)
	}

	/*
		Info("git status --porcelain")
		status, err := w.Status()
		CheckIfError(err)
		fmt.Println(status)
	*/

	// git commit
	now := time.Now().In(autolog.Location)
	_, err = w.Commit(
		fmt.Sprintf("%s %s", now.Format("2006-01-02"), "radicron"),
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "Iori Mizutani",
				Email: "iomz@sazanka.io",
				When:  now,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	/*
			// Prints the current HEAD to verify that all worked well.
			Info("git show -s")
			obj, err := r.CommitObject(commit)
		    if err != nil {
		        panic(err)
		    }
			fmt.Println(obj)
	*/

	// push using default options
	err = r.Push(&git.PushOptions{})
	if err != nil {
		panic(err)
	}
}
