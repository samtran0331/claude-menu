package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type Target struct {
	Name    string
	Cmd     string
	Args    []string
	Env     map[string]string
	CheckFn func() error
}

func main() {
	kimiModel := "kimi-k2-5[1m]"
	qwenModel := "qwen3-coder-30b-a3b-instruct"
	glmModel := "zai-org/glm-4.7-flash"

	// Kimi token comes from the KIMI_API_KEY env var — exported from
	// ~/.zshrc / ~/.bashrc on macOS, or a system env var on Windows.
	kimiKey := os.Getenv("KIMI_API_KEY")

	options := []Target{
		{
			Name: "My Anthropic Pro Subscription",
			Cmd:  "claude",
		},
		{
			Name: "My Kimi Subscription",
			Cmd:  "claude",
			Env: map[string]string{
				"ANTHROPIC_AUTH_TOKEN":            kimiKey,
				"ANTHROPIC_BASE_URL":              "https://api.moonshot.ai/anthropic",
				"ANTHROPIC_MODEL":                 kimiModel,
				"ANTHROPIC_DEFAULT_OPUS_MODEL":    kimiModel,
				"ANTHROPIC_DEFAULT_SONNET_MODEL":  kimiModel,
				"ANTHROPIC_DEFAULT_HAIKU_MODEL":   kimiModel,
				"ANTHROPIC_DEFAULT_FABLE_MODEL":   kimiModel,
				"CLAUDE_CODE_SUBAGENT_MODEL":      kimiModel,
				"CLAUDE_CODE_AUTO_COMPACT_WINDOW": "1048576",
				"CLAUDE_CODE_EFFORT_LEVEL":        "max",
			},
			CheckFn: func() error { return checkKimiKey(kimiKey) },
		},
		{
			Name: "Sifi Vertex API Backend",
			Cmd:  "claude",
			Env: map[string]string{
				"CLAUDE_CODE_USE_VERTEX":      "1",
				"CLOUD_ML_REGION":             "global",
				"ANTHROPIC_MODEL":             "claude-sonnet-4-6",
				"ANTHROPIC_VERTEX_PROJECT_ID": "sbx-advsammac-dev-5ff9",
			},
		},
		{
			Name: "Ollama — Qwen3.6 27B",
			Cmd:  "claude",
			Env: map[string]string{
				"ANTHROPIC_BASE_URL":             "http://localhost:11434/v1",
				"ANTHROPIC_API_KEY":              "ollama",
				"ANTHROPIC_MODEL":                "qwen3.6:27b",
				"ANTHROPIC_DEFAULT_OPUS_MODEL":   "qwen3.6:27b",
				"ANTHROPIC_DEFAULT_SONNET_MODEL": "qwen3.6:27b",
				"ANTHROPIC_DEFAULT_HAIKU_MODEL":  "qwen3.6:27b",
				"ANTHROPIC_DEFAULT_FABLE_MODEL":  "qwen3.6:27b",
				"CLAUDE_CODE_SUBAGENT_MODEL":     "",
			},
			CheckFn: func() error { return checkOllama("qwen3.6:27b") },
		},
		{
			Name: "LM Studio — Qwen3 Coder 30B",
			Cmd:  "claude",
			Env: map[string]string{
				"ANTHROPIC_BASE_URL":             "http://127.0.0.1:1234",
				"ANTHROPIC_API_KEY":              "lm-studio-local",
				"ANTHROPIC_MODEL":                qwenModel,
				"ANTHROPIC_DEFAULT_OPUS_MODEL":   qwenModel,
				"ANTHROPIC_DEFAULT_SONNET_MODEL": qwenModel,
				"ANTHROPIC_DEFAULT_HAIKU_MODEL":  qwenModel,
				"ANTHROPIC_DEFAULT_FABLE_MODEL":  qwenModel,
				"CLAUDE_CODE_SUBAGENT_MODEL":     "",
			},
			CheckFn: func() error { return checkLMStudio(qwenModel) },
		},
		{
			Name: "LM Studio — GLM 4.7 Flash",
			Cmd:  "claude",
			Env: map[string]string{
				"ANTHROPIC_BASE_URL":             "http://127.0.0.1:1234",
				"ANTHROPIC_API_KEY":              "lm-studio-local",
				"ANTHROPIC_MODEL":                glmModel,
				"ANTHROPIC_DEFAULT_OPUS_MODEL":   glmModel,
				"ANTHROPIC_DEFAULT_SONNET_MODEL": glmModel,
				"ANTHROPIC_DEFAULT_HAIKU_MODEL":  glmModel,
				"ANTHROPIC_DEFAULT_FABLE_MODEL":  glmModel,
				"CLAUDE_CODE_SUBAGENT_MODEL":     glmModel,
			},
			CheckFn: func() error { return checkLMStudio(glmModel) },
		},
	}

	selectedIndex := 0

	if err := keyboard.Open(); err != nil {
		fmt.Printf("Failed to open keyboard event loop: %v\n", err)
		return
	}
	defer keyboard.Close()

	render(options, selectedIndex)

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			break
		}

		if key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
			clearScreen()
			fmt.Println("Selection cancelled.")
			return
		} else if key == keyboard.KeyArrowUp {
			selectedIndex = (selectedIndex - 1 + len(options)) % len(options)
			render(options, selectedIndex)
		} else if key == keyboard.KeyArrowDown {
			selectedIndex = (selectedIndex + 1) % len(options)
			render(options, selectedIndex)
		} else if key == keyboard.KeyEnter {
			clearScreen()
			runSelection(options[selectedIndex])
			return
		}
	}
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func render(options []Target, selected int) {
	clearScreen()
	fmt.Println("      💻 CLAUDE CODE ENGINE PICKER 💻     ")

	for idx, opt := range options {
		if idx == selected {
			fmt.Printf(" > \033[36m● %s\033[0m\n", opt.Name)
		} else {
			fmt.Printf("   ○ %s\n", opt.Name)
		}
	}
}

func runSelection(target Target) {
	fmt.Printf("Launching profile: %s...\n\n", target.Name)

	if target.CheckFn != nil {
		if err := target.CheckFn(); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Press any key to exit.")
			fmt.Scanln()
			return
		}
	}

	cleanKeys := []string{
		"CLAUDE_CODE_USE_VERTEX", "CLOUD_ML_REGION", "ANTHROPIC_VERTEX_PROJECT_ID",
		"ANTHROPIC_MODEL", "ANTHROPIC_BASE_URL", "ANTHROPIC_API_KEY", "ANTHROPIC_AUTH_TOKEN",
		"ANTHROPIC_DEFAULT_OPUS_MODEL", "ANTHROPIC_DEFAULT_SONNET_MODEL",
		"ANTHROPIC_DEFAULT_HAIKU_MODEL", "ANTHROPIC_DEFAULT_FABLE_MODEL",
		"CLAUDE_CODE_SUBAGENT_MODEL", "CLAUDE_CODE_AUTO_COMPACT_WINDOW", "CLAUDE_CODE_EFFORT_LEVEL",
	}
	for _, k := range cleanKeys {
		os.Unsetenv(k)
	}

	for k, v := range target.Env {
		os.Setenv(k, v)
	}

	cmd := exec.Command(target.Cmd, target.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func checkKimiKey(key string) error {
	if strings.TrimSpace(key) == "" {
		if runtime.GOOS == "windows" {
			return fmt.Errorf("KIMI_API_KEY not set — add it via System Environment Variables " +
				"(Settings > System > About > Advanced system settings > Environment Variables), " +
				"then open a new terminal")
		}
		return fmt.Errorf("KIMI_API_KEY not set — add `export KIMI_API_KEY=\"sk-...\"` to " +
			"~/.zshrc (or ~/.bashrc), then run: source ~/.zshrc")
	}
	return nil
}

func isPortOpen(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func checkOllama(model string) error {
	if !isPortOpen("127.0.0.1:11434") {
		fmt.Println("Ollama not running. Starting...")
		cmd := exec.Command("ollama", "serve")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start ollama: %v", err)
		}
		for i := 1; i <= 15; i++ {
			time.Sleep(1 * time.Second)
			if isPortOpen("127.0.0.1:11434") {
				fmt.Println("Ollama ready.")
				break
			}
			fmt.Printf("Waiting for Ollama... (%ds)\n", i)
		}
		if !isPortOpen("127.0.0.1:11434") {
			return fmt.Errorf("ollama did not start within 15s")
		}
	}

	out, _ := exec.Command("ollama", "list").Output()
	if !strings.Contains(string(out), model) {
		return fmt.Errorf("model %q not found locally — run: ollama pull %s", model, model)
	}
	fmt.Printf("Model ready: %s\n", model)
	return nil
}

func checkLMStudio(model string) error {
	if !isPortOpen("127.0.0.1:1234") {
		fmt.Println("LM Studio server not running. Starting...")
		cmd := exec.Command("lms", "server", "start")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to start LM Studio server: %v", err)
		}
		for i := 1; i <= 15; i++ {
			time.Sleep(1 * time.Second)
			if isPortOpen("127.0.0.1:1234") {
				break
			}
			fmt.Printf("Waiting for LM Studio server... (%ds)\n", i)
		}
		if !isPortOpen("127.0.0.1:1234") {
			return fmt.Errorf("LM Studio server did not start within 15s")
		}
	}

	out, _ := exec.Command("lms", "ps").Output()
	if !strings.Contains(string(out), model) {
		fmt.Printf("Loading model: %s\n", model)
		cmd := exec.Command("lms", "load", model)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to load model %s: %v", model, err)
		}
	} else {
		fmt.Printf("Model already loaded: %s\n", model)
	}

	return nil
}
