package code

import (
	"isy-cli/internal/openai"
)

// CodeModificationStep represents a step with various operations to be performed.
type CodeModificationStep struct {
	OperationType string       `json:"operation_type" jsonschema:"description=Type of operation: create, delete, edit"`
	FilePath      string       `json:"file_path" jsonschema:"description=The path to the file for the operation, type=string"`
	Edits         []EditDetail `json:"edits" jsonschema:"description=List of edits (optional, used when operation is edit), type=array"`
}

// EditDetail details a specific edit within a file.
type EditDetail struct {
	StartLine int    `json:"start_line" jsonschema:"description=The starting line number for the edit, type=integer"`
	EndLine   int    `json:"end_line" jsonschema:"description=The ending line number for the edit, type=integer"`
	NewCode   string `json:"new_code" jsonschema:"description=The new code to insert, type=string"`
}

// CodeModificationResponse represents a complete set of steps to be executed.
type CodeModificationResponse struct {
	Steps []CodeModificationStep `json:"steps" jsonschema:"description=An ordered list of modification steps to perform, type=array"`
}

var CodeModificationResponseSchema = openai.GenerateSchema[CodeModificationResponse]()

const SYSTEM_PROMPT = `You are an AI assistant for a CLI tool, assisting developers with precise, structured code modifications. You perform operations such as:

1. **Create**: Generate new files at specified paths with provided content.
2. **Delete**: Remove files or specified content within files.
3. **Edit**: Modify specific parts of files, described by line numbers and new content.

Each task should consist of an ordered list of steps. Each step can contain multiple operations. Particularly for edits, provide detailed line information (start and end lines) and the code to substitute.

While engaging with the developer:
- Provide clear steps that outline file paths and operations type.
- Confirm exact modifications before applying to prevent errors.
- Suggest any shell commands needed after modifications for task validation or better integration.

Maintain technical accuracy and clear communication to facilitate efficient and error-free code adjustments.
`
