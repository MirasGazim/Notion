package service

import (
	"context"
	"errors"
	"fmt"
	"notion/internal/models/blocks"

	"github.com/google/uuid"
)

var (
	ErrInvalidContent = errors.New("invalid block content")
)

type BlockService interface {
	CreateBlock(ctx context.Context, req blocks.CreateBlockRequest, WorkspaceID, UserID uuid.UUID) (*blocks.Block, error)
}

type ValidationError struct {
	Err error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid block content: %v", e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

func (e *ValidationError) Is(target error) bool {
	return target == ErrInvalidContent
}

func (s *blockService) CreateBlock(ctx context.Context, req blocks.CreateBlockRequest, workspaceID, userID uuid.UUID) (*blocks.Block, error) {
	const op = "service.CreateBlock"
	err := ValidateBlockContent(req.Type, req.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, &ValidationError{Err: err})
	}

	block, err := s.repo.Create(ctx, req, workspaceID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return block, nil
}

func ValidateBlockContent(blocktype string, content map[string]interface{}) error {
	switch blocktype {
	case "checkbox":
		return validateCheckboxContent(content)
	case "text":
		return validateTextContent(content)
	case "date":
		return validateDateContent(content)
	case "files":
		return validateFilesContent(content)
	case "number":
		return validateNumberContent(content)
	case "person":
		return validatePersonContent(content)
	case "email":
		return validateEmailContent(content)
	case "url":
		return validateURLContent(content)
	case "phone":
		return validatePhoneContent(content)
	case "multiple_choice":
		return validateMultipleChoiceContent(content)
	default:
		return fmt.Errorf("unknown block type: %s", blocktype)
	}
}

func validateCheckboxContent(content map[string]interface{}) error {
	val, ok := content["checked"]
	if !ok {
		return errors.New("checkbox: missing 'checked' field")
	}
	if _, ok := val.(bool); !ok {
		return errors.New("checkbox: 'checked' must be boolean")
	}
	return nil
}

func validateTextContent(content map[string]interface{}) error {
	val, ok := content["text"]
	if !ok {
		return errors.New("text: missing 'text' field")
	}
	if _, ok := val.(string); !ok {
		return errors.New("text: 'text' must be string")
	}
	return nil
}

func validateDateContent(content map[string]interface{}) error {
	val, ok := content["date"]
	if !ok {
		return errors.New("date: missing 'date' field")
	}
	if _, ok := val.(string); !ok {
		return errors.New("date: 'date' must be string")
	}
	return nil
}

func validateNumberContent(content map[string]interface{}) error {
	val, ok := content["number"]
	if !ok {
		return errors.New("number: missing 'number' field")
	}
	if _, ok := val.(float64); !ok {
		return errors.New("number: 'number' must be float")
	}
	return nil
}

func validateEmailContent(content map[string]interface{}) error {
	val, ok := content["email"]
	if !ok {
		return errors.New("email: missing 'email' field")
	}
	if _, ok := val.(string); !ok {
		return errors.New("email: 'email' must be string")
	}
	return nil
}

func validateURLContent(content map[string]interface{}) error {
	val, ok := content["url"]
	if !ok {
		return errors.New("url: missing 'url' field")
	}
	if _, ok := val.(string); !ok {
		return errors.New("url: 'url' must be string")
	}
	return nil
}

func validatePhoneContent(content map[string]interface{}) error {
	val, ok := content["phone"]
	if !ok {
		return errors.New("phone: missing 'phone' field")
	}
	if _, ok = val.(string); !ok {
		return errors.New("phone: 'phone' must be string")
	}
	return nil
}

func validatePersonContent(content map[string]interface{}) error {
	val, ok := content["user_ids"]
	if !ok {
		return errors.New("user_ids: missing 'user_ids' field")
	}
	persons, ok := val.([]interface{})
	if !ok {
		return errors.New("user_ids: 'user_ids' must be an array")
	}
	for _, person := range persons {
		_, ok := person.(string)
		if !ok {
			return errors.New("person: each user id must be a string")
		}
	}
	return nil
}

func validateFilesContent(content map[string]interface{}) error {
	val, ok := content["files"]
	if !ok {
		return errors.New("files: missing 'files' field")
	}
	files, ok := val.([]interface{})
	if !ok {
		return errors.New("files: 'files' must be an array")
	}
	for _, file := range files {
		fileMap, ok := file.(map[string]interface{})
		if !ok {
			return errors.New("files: each file must be an object")
		}
		if _, ok := fileMap["url"].(string); !ok {
			return errors.New("files: each file must have a 'url' field of type string")
		}
	}
	return nil
}

func validateMultipleChoiceContent(content map[string]interface{}) error {
	val, ok := content["options"]
	if !ok {
		return errors.New("options: missing 'options' field")
	}
	options, ok := val.([]interface{})
	if !ok {
		return errors.New("options: 'options' must be an array")
	}
	for _, option := range options {
		_, ok := option.(string)
		if !ok {
			return errors.New("multiple_choice: each option must be a string")
		}
	}
	return nil
}
