package worker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mukund/mediaconvert/internal/pipeline"
)

// ExecutionContext holds context for pipeline execution
type ExecutionContext struct {
	InputFile  string
	OutputDir  string
	WorkDir    string
	Variables  map[string]string
}

// ExecutePipeline executes all steps in a pipeline
func ExecutePipeline(p *pipeline.Pipeline, inputFile, workDir string) ([]string, error) {
	// Create output directory
	outputDir := filepath.Join(workDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	ctx := &ExecutionContext{
		InputFile: inputFile,
		OutputDir: outputDir,
		WorkDir:   workDir,
		Variables: make(map[string]string),
	}

	var outputFiles []string

	// Execute each step sequentially
	for i, step := range p.Steps {
		fmt.Printf("Executing step %d: %s (%s)\n", i+1, step.Operation, step.Output)

		// Map operation to command
		cmd, err := MapOperation(step, ctx)
		if err != nil {
			return nil, fmt.Errorf("step %d: failed to map operation: %w", i+1, err)
		}

		// Execute command
		if err := executeCommand(cmd); err != nil {
			return nil, fmt.Errorf("step %d: command failed: %w", i+1, err)
		}

		// Track output file
		outputPath := substituteVars(step.Output, ctx)
		outputFiles = append(outputFiles, outputPath)

		fmt.Printf("Step %d completed: %s\n", i+1, outputPath)
	}

	return outputFiles, nil
}

func executeCommand(cmd *OperationCommand) error {
	fmt.Printf("Running: %s %v\n", cmd.Tool, cmd.Args)

	command := exec.Command(cmd.Tool, cmd.Args...)

	// Capture output
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s failed: %w\nOutput: %s", cmd.Tool, err, string(output))
	}

	return nil
}
