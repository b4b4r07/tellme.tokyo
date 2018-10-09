package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sort"

	"github.com/b4b4r07/go-finder/source"
	"github.com/k0kubun/pp"
)

// EditCommand is one of the subcommands
type EditCommand struct {
	CLI
	Option EditOption
}

// EditOption is the options for EditCommand
type EditOption struct {
	Tag  bool
	Open bool
}

func (c *EditCommand) flagSet() *flag.FlagSet {
	flags := flag.NewFlagSet("edit", flag.ExitOnError)
	flags.BoolVar(&c.Option.Tag, "tag", false, "edit article with tag")
	flags.BoolVar(&c.Option.Open, "open", false, "open article with browser when editing")
	return flags
}

// Run run edit command
func (c *EditCommand) Run(args []string) int {
	flags := c.flagSet()
	if err := flags.Parse(args); err != nil {
		return c.exit(err)
	}

	var files []string
	var err error
	if c.Option.Tag {
		files, err = c.selectFilesWithTag()
	} else {
		files, err = c.selectFiles()
	}
	if err != nil {
		return c.exit(err)
	}

	return c.exit(c.edit(files))
}

// Synopsis returns synopsis
func (c *EditCommand) Synopsis() string {
	return "Edit blog articles"
}

// Help returns help message
func (c *EditCommand) Help() string {
	var b bytes.Buffer
	flags := c.flagSet()
	flags.SetOutput(&b)
	flags.PrintDefaults()
	return fmt.Sprintf(
		"Usage of %s:\n\nOptions:\n%s", flags.Name(), b.String(),
	)
}

func (c *EditCommand) selectFilesWithTag() ([]string, error) {
	var files []string
	articles, err := walk(c.Config.BlogDir, 1)
	if err != nil {
		return files, err
	}

	var tags []string
	for _, article := range articles {
		tags = append(tags, article.Body.Tags...)
	}
	sort.Strings(tags)
	tags = uniqSlice(tags)
	c.Finder.Read(source.Slice(tags))
	items, err := c.Finder.Run()
	if err != nil {
		return files, err
	}
	for _, item := range items {
		for _, article := range articles.Filter(item) {
			files = append(files, article.Path)
		}
	}
	return files, nil
}

func (c *EditCommand) selectFiles() ([]string, error) {
	articles, err := walk(filepath.Join(c.Config.BlogDir, "content", "post"), 1)
	if err != nil {
		return []string{}, err
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})
	for _, article := range articles {
		c.Finder.Add(article.Body.Title, article)
	}
	var files []string
	items, err := c.Finder.Select()
	for _, item := range items {
		files = append(files, item.(Article).Path)
	}
	return files, err
}

func (c *EditCommand) edit(files []string) error {
	if len(files) == 0 {
		return nil
	}

	// os.Chdir()

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer signal.Stop(ch)
	defer cancel()
	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()

	go newHugo(c.Config.BlogDir, "server", "-D").Run(ctx)

	if c.Option.Open {
		quit := make(chan bool)
		go func() {
			// discard error
			runCommand("open", envHostURL)
			quit <- true
		}()
		<-quit
	}

	vim := newShell("vim", files...)
	return vim.Run(context.Background())
}

func newHugo(dir string, args ...string) shell {
	pp.Println(dir)
	return shell{
		stdin:  os.Stdin,
		stdout: ioutil.Discard, // to /dev/null
		// stderr:  ioutil.Discard, // to /dev/null
		stderr:  os.Stderr,
		env:     map[string]string{},
		command: "hugo",
		args:    args,
		dir:     dir,
	}
}
