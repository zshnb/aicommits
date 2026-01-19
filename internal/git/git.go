package git

import (
	"fmt"
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

func Commit(msg string) {
	commitCmd := exec.Command("git", "commit", "-m", msg)
	if out, err := commitCmd.CombinedOutput(); err != nil {
		fmt.Printf("❌ 提交失败:\n%s\n", string(out))
	} else {
		fmt.Println("✅ 提交成功!")
		fmt.Println(string(out))
	}
}
