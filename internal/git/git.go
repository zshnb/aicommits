package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func Commit(msg string) {
	commitCmd := exec.Command("git", "commit", "-m", msg)
	if out, err := commitCmd.CombinedOutput(); err != nil {
		fmt.Printf("❌ 提交失败:\n%s\n", string(out))
	} else {
		fmt.Println(string(out))
	}
}

func StageAll() error {
	cmd := exec.Command("git", "add", ".")
	return cmd.Run()
}
