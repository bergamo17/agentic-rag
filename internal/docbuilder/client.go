package docbuilder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Section map[string]interface{}

type BuildRequest struct {
	Theme      string    `json:"theme"`
	Title      string    `json:"title"`
	Sections   []Section `json:"sections"`
	OutputPath string    `json:"output_path"`
}

type docResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	OutputPath string `json:"output_path,omitempty"`
}

type Client struct {
	PythonPath string
	ScriptPath string
}

func NewClient(pythonPath, scriptPath string) *Client {
	return &Client{
		PythonPath: pythonPath,
		ScriptPath: scriptPath,
	}
}

func (c *Client) Build(req BuildRequest) (string, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("docbuilder: failed to marshal the request: %w", err)
	}

	cmd := exec.Command(c.PythonPath, c.ScriptPath)
	cmd.Stdin = bytes.NewReader(payload)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()

	var resp docResponse
	if jsonErr := json.Unmarshal(stdout.Bytes(), &resp); jsonErr != nil {
		return "", fmt.Errorf(
			"docbuilder: script produced no valid JSON output (exit err: %v): %s",
			runErr,
			stderr.String(),
		)
	}

	if !resp.Success {
		return "", fmt.Errorf("docbuilder: %s", resp.Error)
	}

	return resp.OutputPath, nil
}
