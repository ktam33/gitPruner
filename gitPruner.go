package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	path := os.Args[1]
	os.Chdir(path)

	info("\nSyncing branches with remote...[git fetch -p]")

	out, err := exec.Command("git", "fetch", "-p").CombinedOutput()
	if err != nil {
		fmt.Println("Unable to reach remote. Skipping sync step.\n")
	} else {
		fmt.Println(string(out))
	}

	info("\nGetting information about local branches...[git branch -vv]")
	out, err = exec.Command("git", "branch", "-vv").CombinedOutput()
	checkError(err)
	branchInfo := string(out)
	fmt.Println(branchInfo)

	lines := strings.Split(branchInfo, "\n")
	re := regexp.MustCompile("\\s([^\\s]+?)\\s")

	reader := bufio.NewReader(os.Stdin)
	foundBranchToDelete := false
	for _, line := range lines {
		isMatch, err := regexp.MatchString("\\[.+: gone\\]", line)
		checkError(err)
		if isMatch {
			foundBranchToDelete = true
			branchName := re.FindStringSubmatch(line)[1]
			info("\nBranch \"%s\" is no longer a remote branch. Delete? (y or n)", branchName)
			input, err := reader.ReadString('\n')
			checkError(err)
			input = strings.ToLower(input)
			if input[0] == 'y' {
				info("Deleting branch...[git branch -d %s]", branchName)
				out, err = exec.Command("git", "branch", "-d", branchName).CombinedOutput()
				fmt.Println(string(out))
			}
		}
	}
	if !foundBranchToDelete {
		info("No local branches to prune were found")
	}
}

func checkError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error: %s", err)
	os.Exit(1)
}

func info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
