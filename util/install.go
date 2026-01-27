package util

import (
	"fmt"
	"os/exec"
	"os"
	"strings"
)

func InstallDeps() {
	installHomebrew()
	installDuti()
}

func installHomebrew() {
	fmt.Println("Check Homebrew Environment")
	if !commandExists("brew") {
		fmt.Println("Homebrew not exists, installing ...")
		cmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
		_, _ = cmd.Output()
		updatePathForHomebrew()
	}

	cmd := exec.Command("brew", "--help")
	_, err := cmd.Output()

	if err != nil {
		fmt.Printf("Homebrew error: %v\n", err)
	} else {
		fmt.Println("Homebrew works fine")
	}
}

func updatePathForHomebrew() {
	brewPath := "/home/linuxbrew/.linuxbrew/bin"
	path := os.Getenv("PATH")
	if !strings.Contains(path, brewPath) {
		newPath := fmt.Sprintf("%s:%s", brewPath, path)
		os.Setenv("PATH", newPath)
	}
}

func installDuti() {
	fmt.Println("Check Duti Environment")
	if !commandExists("duti") {
		fmt.Println("Duti not exists, installing ...")
		cmd := exec.Command("brew", "install", "duti")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			fmt.Println("Error installing duti:", err)
			return
		}
	}

	// Verify duti is executable by checking if command exists
	if commandExists("duti") {
		fmt.Println("Duti works fine")
	} else {
		fmt.Println("Error: duti command not found after installation")
	}
}
