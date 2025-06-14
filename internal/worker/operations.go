package worker

import (
	"fmt"
	"strings"

	"github.com/mukund/mediaconvert/internal/pipeline"
)

// OperationCommand represents a command to execute
type OperationCommand struct {
	Tool string
	Args []string
}

// MapOperation maps an abstract operation to a concrete command
func MapOperation(step pipeline.Step, context *ExecutionContext) (*OperationCommand, error) {
	switch step.Operation {
	case "transcode":
		return mapTranscode(step, context)
	case "resize":
		return mapResize(step, context)
	case "extract_text":
		return mapExtractText(step, context)
	case "extract_frame":
		return mapExtractFrame(step, context)
	case "convert":
		return mapConvert(step, context)
	case "generate_thumbnail":
		return mapGenerateThumbnail(step, context)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", step.Operation)
	}
}

func mapTranscode(step pipeline.Step, ctx *ExecutionContext) (*OperationCommand, error) {
	args := []string{"-i", substituteVars(step.Input, ctx)}

	// Map codec
	if codec, ok := step.Params["codec"].(string); ok {
		switch codec {
		case "h264":
			args = append(args, "-c:v", "libx264")
		case "h265":
			args = append(args, "-c:v", "libx265")
		case "vp9":
			args = append(args, "-c:v", "libvpx-vp9")
		default:
			args = append(args, "-c:v", codec)
		}
	}

	// Map quality (CRF for video)
	if quality, ok := step.Params["quality"]; ok {
		switch v := quality.(type) {
		case float64:
			args = append(args, "-crf", fmt.Sprintf("%.0f", v))
		case int:
			args = append(args, "-crf", fmt.Sprintf("%d", v))
		}
	}

	// Audio codec
	if audioCodec, ok := step.Params["audio_codec"].(string); ok {
		args = append(args, "-c:a", audioCodec)
	}

	// Audio bitrate
	if audioBitrate, ok := step.Params["audio_bitrate"].(string); ok {
		args = append(args, "-b:a", audioBitrate)
	}

	args = append(args, substituteVars(step.Output, ctx))

	return &OperationCommand{
		Tool: "ffmpeg",
		Args: args,
	}, nil
}

func mapResize(step pipeline.Step, ctx *ExecutionContext) (*OperationCommand, error) {
	args := []string{substituteVars(step.Input, ctx)}

	// Build resize argument
	var resizeArg string
	if width, ok := step.Params["width"]; ok {
		if height, ok := step.Params["height"]; ok {
			resizeArg = fmt.Sprintf("%vx%v", width, height)
		}
	}

	if resizeArg != "" {
		args = append(args, "-resize", resizeArg)
	}

	// Quality
	if quality, ok := step.Params["quality"]; ok {
		args = append(args, "-quality", fmt.Sprintf("%v", quality))
	}

	args = append(args, substituteVars(step.Output, ctx))

	return &OperationCommand{
		Tool: "convert",
		Args: args,
	}, nil
}

func mapExtractText(step pipeline.Step, ctx *ExecutionContext) (*OperationCommand, error) {
	args := []string{
		substituteVars(step.Input, ctx),
		substituteVars(step.Output, ctx),
	}

	return &OperationCommand{
		Tool: "pdftotext",
		Args: args,
	}, nil
}

func mapExtractFrame(step pipeline.Step, ctx *ExecutionContext) (*OperationCommand, error) {
	args := []string{"-i", substituteVars(step.Input, ctx)}

	// Timestamp
	if timestamp, ok := step.Params["timestamp"].(string); ok {
		args = append(args, "-ss", timestamp)
	}

	args = append(args, "-vframes", "1")
	args = append(args, substituteVars(step.Output, ctx))

	return &OperationCommand{
		Tool: "ffmpeg",
		Args: args,
	}, nil
}

func mapConvert(step pipeline.Step, ctx *ExecutionContext) (*OperationCommand, error) {
	// Similar to resize but for format conversion
	args := []string{
		substituteVars(step.Input, ctx),
		substituteVars(step.Output, ctx),
	}

	return &OperationCommand{
		Tool: "convert",
		Args: args,
	}, nil
}

func mapGenerateThumbnail(step pipeline.Step, ctx *ExecutionContext) (*OperationCommand, error) {
	input := substituteVars(step.Input, ctx)
	output := substituteVars(step.Output, ctx)

	// Determine input type from params or extension
	inputType := "video" // default
	if t, ok := step.Params["type"].(string); ok {
		inputType = t
	}

	switch inputType {
	case "video":
		// Use FFmpeg to extract frame from video
		args := []string{"-i", input}

		// Timestamp (default to 1 second)
		timestamp := "00:00:01"
		if t, ok := step.Params["timestamp"].(string); ok {
			timestamp = t
		}
		args = append(args, "-ss", timestamp)

		// Size
		if width, ok := step.Params["width"]; ok {
			if height, ok := step.Params["height"]; ok {
				args = append(args, "-vf", fmt.Sprintf("scale=%v:%v", width, height))
			}
		}

		args = append(args, "-vframes", "1", output)

		return &OperationCommand{
			Tool: "ffmpeg",
			Args: args,
		}, nil

	case "image":
		// Use ImageMagick to resize image
		args := []string{input}

		// Size
		if width, ok := step.Params["width"]; ok {
			if height, ok := step.Params["height"]; ok {
				args = append(args, "-resize", fmt.Sprintf("%vx%v", width, height))
			}
		}

		args = append(args, output)

		return &OperationCommand{
			Tool: "convert",
			Args: args,
		}, nil

	case "pdf":
		// Use ImageMagick to convert first page of PDF to image
		args := []string{input + "[0]"} // [0] selects first page

		// Size
		if width, ok := step.Params["width"]; ok {
			if height, ok := step.Params["height"]; ok {
				args = append(args, "-resize", fmt.Sprintf("%vx%v", width, height))
			}
		}

		args = append(args, output)

		return &OperationCommand{
			Tool: "convert",
			Args: args,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported thumbnail type: %s", inputType)
	}
}

func substituteVars(s string, ctx *ExecutionContext) string {
	s = strings.ReplaceAll(s, "${input}", ctx.InputFile)
	s = strings.ReplaceAll(s, "${output}", ctx.OutputDir)
	return s
}
