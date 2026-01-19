package git

import (
	"os/exec"
	"strings"
)

// GetStagedDiff 获取暂存区的差异
func GetStagedDiff() (string, error) {
	// 执行 git diff --cached --diff-algorithm=minimal
	// minimal算法能产生更紧凑的diff，节省Token
	cmd := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
