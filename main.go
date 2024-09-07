package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type AiResponse struct {
	Response string `json:"response"`
}

func main() {
	userMessage := "Write a git commit message for these changes. The commit message should contain no more than 30 words and no less 20 words. The response should contain only the message."
	prompt, err := buildPrompt(userMessage)
	if err != nil {
		log.Fatal(err)
	}

	responseBody, err := requestToModel(prompt)
	if err != nil {
		log.Fatalf("Failed to make HTTP request: %v", err)
	}

	var response AiResponse

	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		log.Fatal(err)
	}

	commitMessage := response.Response

	fmt.Println("AI commit message:", commitMessage)
	fmt.Println("Apply message for commit? [y/any]: ")

	var choice string

	fmt.Scanln(&choice)
	if choice == "y" {
		execGitCommand("add", ".")
		execGitCommand("commit", "-m", commitMessage)
	} else {
		fmt.Println("abort")
	}
}

func requestToModel(prompt string) ([]byte, error) {
	data := map[string]interface{}{"model": "gemma2:2b", "prompt": prompt, "stream": false}

	body, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	response, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func buildPrompt(message string) (string, error) {
	status, err := execGitCommand("status")
	if err != nil {
		return "", fmt.Errorf("%e", err)
	}

	diff, err := execGitCommand("diff")
	if err != nil {
		return "", fmt.Errorf("%e", err)
	}

	prompt := strings.ReplaceAll(message+status+diff, "\"", "")

	return prompt, nil
}

func execGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run git command: %v\nOutput: %s", err, output)
	}

	return string(output), nil
}
