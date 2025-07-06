package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var base string
var stepsFile = "steps.json"

type StepsConfig struct {
	Base    string   `json:"base"`
	Apply   []string `json:"apply"`
	Destroy []string `json:"destroy"`
}

var (
	mode     string
	logFile  *os.File
	logDir   = "logs"
)

func init() {
	flag.StringVar(&mode, "mode", "", "Operation mode: apply or destroy")
}

func main() {
	flag.Parse()
	if mode != "apply" && mode != "destroy" {
		fmt.Println(colorRed("[ERROR]") + " --mode must be 'apply' or 'destroy'")
		os.Exit(1)
	}

	initLogging()
	defer logFile.Close()

	steps := loadSteps(mode)

	for _, step := range steps {
		runTerragrunt(filepath.Join(base, step), mode)
	}

	consoleLog("[INFO]", fmt.Sprintf("✅ Terragrunt '%s' completed for all steps.", mode))
}

func initLogging() {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	logPath := filepath.Join(logDir, fmt.Sprintf("terragrunt-run-%s.log", timestamp))

	var err error
	logFile, err = os.Create(logPath)
	if err != nil {
		fmt.Println(colorRed("[ERROR]") + " Failed to create log file:", err)
		os.Exit(1)
	}

	consoleLog("[INFO]", fmt.Sprintf("Logging to %s", logPath))
}

func loadSteps(mode string) []string {
	data, err := os.ReadFile(stepsFile)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read %s: %v", stepsFile, err)
	}

	var cfg StepsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("[ERROR] Failed to parse %s: %v", stepsFile, err)
	}

	base = cfg.Base // set global base from JSON

	switch mode {
	case "apply":
		return cfg.Apply
	case "destroy":
		return cfg.Destroy
	default:
		return nil
	}
}

func runTerragrunt(dir, mode string) {
	consoleLog("[INFO]", fmt.Sprintf("Running %s in %s", mode, dir))

	cmd := exec.Command("terragrunt", mode, "--non-interactive", "-auto-approve")
	cmd.Dir = dir

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		logAndPrint("[ERROR]", fmt.Sprintf("Failed to start terragrunt in %s: %v", dir, err))
		os.Exit(1)
	}

	go streamOutput(stdout, "[TG]")
	go streamOutput(stderr, "[TG-ERR]")

	if err := cmd.Wait(); err != nil {
		logAndPrint("[ERROR]", fmt.Sprintf("Terragrunt failed in %s: %v", dir, err))
		os.Exit(1)
	}
}

func streamOutput(pipe io.Reader, prefix string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()

		// Save raw log output
		fmt.Fprintf(logFile, "%s %s\n", prefix, line)

		// Print colored to console
		switch {
		case strings.Contains(line, "Apply complete"), strings.Contains(line, "Destroy complete"):
			fmt.Println(colorCyan("[SUMMARY]") + " " + line)
		case prefix == "[TG-ERR]":
			fmt.Println(colorRed(prefix) + " " + line)
		default:
			fmt.Println(colorWhite(prefix) + " " + line)
		}
	}
}

// ===== Utility logging and colors =====

func consoleLog(level, msg string) {
	color := colorFor(level)
	fmt.Println(color(level) + " " + msg)
	fmt.Fprintf(logFile, "%s %s\n", level, msg)
}

func logAndPrint(level, msg string) {
	consoleLog(level, msg)
}

func colorFor(level string) func(string) string {
	switch level {
	case "[ERROR]":
		return colorRed
	case "[INFO]":
		return colorGreen
	case "[SUMMARY]":
		return colorCyan
	default:
		return colorWhite
	}
}

func colorRed(s string) string    { return "\033[1;31m" + s + "\033[0m" }
func colorGreen(s string) string  { return "\033[1;32m" + s + "\033[0m" }
func colorCyan(s string) string   { return "\033[1;36m" + s + "\033[0m" }
func colorWhite(s string) string  { return "\033[1;37m" + s + "\033[0m" }
