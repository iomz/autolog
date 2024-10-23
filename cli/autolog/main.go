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

	. "github.com/go-git/go-git/v5/_examples"
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

	directory := "/Users/iomz/ghq/github.com/iomz/logs"
	r, err := git.PlainOpen(directory)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	_, err = w.Add(".")
	CheckIfError(err)

	Info("git status --porcelain")
	status, err := w.Status()
	CheckIfError(err)

	fmt.Println(status)

	Info("git commit -m \"example go-git commit\"")
	now := time.Now().In(autolog.Location)
	commit, err := w.Commit(
		fmt.Sprintf("%s %s", now.Format("2006-01-02"), "radicron"),
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "Iori Mizutani",
				Email: "iomz@sazanka.io",
				When:  now,
			},
		},
	)

	CheckIfError(err)

	// Prints the current HEAD to verify that all worked well.
	Info("git show -s")
	obj, err := r.CommitObject(commit)
	CheckIfError(err)

	fmt.Println(obj)

	Info("git push")
	// push using default options
	err = r.Push(&git.PushOptions{})
	CheckIfError(err)
}
