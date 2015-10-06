package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/github/hub/cmd"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

const (
	name          = "git-fire"
	version       = "0.0.1"
	branchFmt     = "gf-%d-%s-%s"
	leaveBuilding = "\n\n\n＿人人人人人人人人人＿\n＞ Leave building ! ＜\n￣Y^Y^Y^Y^Y^Y^Y^Y^Y￣\n\n:runner: :dash: :dash: :dash: \t :fire: :office:\n"
)

var message = emoji.Sprint(":name_badge:") + " The fire, run!"

func main() {

	cmd := &cobra.Command{
		Use: name,
		Run: gitFire,
	}

	cmd.Execute()
}

func gitFire(cmd *cobra.Command, args []string) {

	color.New(color.FgCyan, color.Bold).Println("\nIn case of fire ...\n")

	wg := sync.WaitGroup{}
	email := "anonymous"
	ref := "noref"

	wg.Add(2)

	// get user email
	go func() {

		defer wg.Done()

		resp, err := git("config", "user.email")
		exitIf(err)

		if len(resp) != 0 {
			email = resp[0]
		}
	}()

	// get head ref
	go func() {

		defer wg.Done()

		resp, err := git("rev-parse", "--short", "-q", "head")
		exitIf(err)

		if len(resp) != 0 {
			ref = resp[0]
		}
	}()

	wg.Wait()

	// checkout new branch
	branch := fmt.Sprintf(branchFmt, time.Now().Unix(), email, ref)
	_, err := git("checkout", "-b", branch)
	exitIf(err)

	fmt.Printf("Branch:\n\t%s\n", branch)

	// commit all files
	git("add", "-A")

	if len(args) != 0 {
		message = args[0]
	}
	_, err = git("commit", "--allow-empty", "-m", message)
	exitIf(err)

	fmt.Printf("Message:\n\t%s\n\n\n", message)

	// push remotes
	remotes, err := git("remote")
	exitIf(err)

	count := len(remotes)
	wg.Add(count)
	fmt.Printf("Push (%d remote(s)):\n", count)
	for i := range remotes {
		go func(remote string) {

			defer wg.Done()

			_, err = git("push", remote, branch)
			isPushed(remote, err)

		}(remotes[i])
	}

	wg.Wait()

	// leave building
	emoji.Println(leaveBuilding)
}

func git(input ...string) (outputs []string, err error) {

	cmd := cmd.New("git")

	for _, i := range input {
		cmd.WithArg(i)
	}

	out, err := cmd.CombinedOutput()
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			outputs = append(outputs, string(line))
		}
	}

	return
}

func isPushed(remote string, err error) {

	if err == nil {
		fmt.Printf("\t%s: %s\n", remote, color.GreenString("OK"))
		return
	}

	fmt.Printf("\t%s: %s\n", remote, color.RedString("NG"))
}

func exitIf(err error) {

	if err == nil {
		return
	}

	fmt.Println(color.RedString("ERROR:"), err, emoji.Sprint(":fire:"))
	os.Exit(1)
}
