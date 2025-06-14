package pipeline

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Pipeline represents a processing pipeline
type Pipeline struct {
	Name  string `json:"name" yaml:"name"`
	Steps []Step `json:"steps" yaml:"steps"`
}

// Step represents a single processing step
type Step struct {
	Operation string                 `json:"operation" yaml:"operation"`
	Input     string                 `json:"input" yaml:"input"`
	Output    string                 `json:"output" yaml:"output"`
	Params    map[string]interface{} `json:"params,omitempty" yaml:"params,omitempty"`
}

// ParseYAML parses a YAML pipeline definition
func ParseYAML(data []byte) (*Pipeline, error) {
	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &p, nil
}

// ParseJSON parses a JSON pipeline definition
func ParseJSON(data []byte) (*Pipeline, error) {
	var p Pipeline
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &p, nil
}

// ToYAML converts pipeline to YAML
func (p *Pipeline) ToYAML() ([]byte, error) {
	return yaml.Marshal(p)
}

// ToJSON converts pipeline to JSON
func (p *Pipeline) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// Validate performs basic validation on the pipeline
func (p *Pipeline) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("pipeline name is required")
	}
	if len(p.Steps) == 0 {
		return fmt.Errorf("pipeline must have at least one step")
	}
	for i, step := range p.Steps {
		if step.Operation == "" {
			return fmt.Errorf("step %d: operation is required", i)
		}
		if step.Input == "" {
			return fmt.Errorf("step %d: input is required", i)
		}
		if step.Output == "" {
			return fmt.Errorf("step %d: output is required", i)
		}
	}
	return nil
}
