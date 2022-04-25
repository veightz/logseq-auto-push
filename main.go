package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	LogseqWordDir = "/Users/veightz/my-logseq"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	log.Println("init ...")
	DebugLogger = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	for true {
		doCommitTask()
		doPushTask()
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func doCommitTask() {
	needCommit, err := hasUncommittedChanges()
	if err != nil {
		ErrorLogger.Fatalf("exec `git status` failed.error=%v", err)
		return
	}
	if !needCommit {
		DebugLogger.Println("need not to commit.")
	} else {
		InfoLogger.Println("start committing.")
		err := execCommit()
		if err != nil {
			ErrorLogger.Fatalf("check failed.error=%v", err)
		} else {
			InfoLogger.Println("end committing.")
		}
	}
}

func doPushTask() {
	isAheadOf, err := localVersionIsAheadOfOrigin()
	if err != nil {
		ErrorLogger.Fatalf("exec `git status` failed.error=%v", err)
		return
	}
	if !isAheadOf {
		DebugLogger.Println("need not to push.")
	} else {
		InfoLogger.Println("start pushing.")
		err := execPush()
		if err != nil {
			ErrorLogger.Fatalf("check failed.error=%v", err)
		} else {
			InfoLogger.Println("end pushing.")
		}
	}
}

func execCommit() error {
	cmd := exec.Command("git", "commit", "-am", "Auto saved")
	cmd.Dir = LogseqWordDir
	_, err := cmd.Output()
	return err
}

func execPush() error {
	cmd := exec.Command("git", "push", "origin")
	cmd.Dir = LogseqWordDir
	_, err := cmd.Output()
	return err
}

func hasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "-b")
	cmd.Dir = LogseqWordDir
	out, err := cmd.Output()
	if err != nil {
		return false, err
	} else {
		result := string(out)
		return len(result) > 0 && containsChangesToCommit(result), nil
	}
}

func localVersionIsAheadOfOrigin() (bool, error) {
	cmd := exec.Command("git", "status", "-b")
	cmd.Dir = LogseqWordDir
	out, err := cmd.Output()
	if err != nil {
		return false, err
	} else {
		result := string(out)
		return len(result) > 0 && containsAheadOfCommits(result), nil
	}
}

func containsAheadOfCommits(r string) bool {
	return strings.Contains(r, "Your branch is ahead of")
}

func containsChangesToCommit(r string) bool {
	return strings.Contains(r, "Changes to be committed:") ||
		strings.Contains(r, "Changes not staged for commit:")
}
