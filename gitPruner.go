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

	syncBranches()

	lines := getLocalBranches()

	re := regexp.MustCompile("\\s([^\\s]+?)\\s")

	reader := bufio.NewReader(os.Stdin)
	foundBranchToDelete := false
	for _, line := range lines {
		isMatch, err := regexp.MatchString("\\[.+: gone\\]", line)
		checkError(err)
		if isMatch {
			// skip the currently checked out branch
			branchName := re.FindStringSubmatch(line)[1]
			if line[0] == '*' {
				info("Branch \"%s\" is no longer a remote branch but skipping since it is checked out.", branchName)
				break
			}
			foundBranchToDelete = true
			info("Branch \"%s\" is no longer a remote branch. Delete? (y or n)", branchName)
			input, err := reader.ReadString('\n')
			checkError(err)
			input = strings.ToLower(input)
			if input[0] == 'y' {
				deleteBranch(branchName, reader)
			}
		}
	}
	if !foundBranchToDelete {
		info("No local branches to prune were found")
	}
}

func syncBranches() {
	info("\nSyncing branches with remote...[git fetch -p]")
	out, err := exec.Command("git", "fetch", "-p").CombinedOutput()
	if err != nil {
		fmt.Println("Unable to reach remote. Skipping sync step.\n")
	} else {
		fmt.Println(string(out))
	}
}

func getLocalBranches() []string {
	info("Getting information about local branches...[git branch -vv]")
	out, err := exec.Command("git", "branch", "-vv").CombinedOutput()
	checkError(err)
	branchInfo := string(out)
	fmt.Println(branchInfo)
	return strings.Split(branchInfo, "\n")
}

func deleteBranch(branchName string, reader *bufio.Reader) {
	info("Deleting branch...[git branch -d %s]", branchName)
	out, err := exec.Command("git", "branch", "-d", branchName).CombinedOutput()
	checkError(err)
	outText := string(out)
	fmt.Println(outText)
	isMatch, err := regexp.MatchString("error: The branch '.+' is not fully merged\\.", outText)
	checkError(err)

	if isMatch {
		info("Force delete branch \"%s\"? (y or n)", branchName)
		input, err := reader.ReadString('\n')
		checkError(err)
		input = strings.ToLower(input)
		if input[0] == 'y' {
			info("Force deleting branch...[git branch -D %s]", branchName)
			out, err = exec.Command("git", "branch", "-D", branchName).CombinedOutput()
			checkError(err)
			fmt.Println(string(out))
		}
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
