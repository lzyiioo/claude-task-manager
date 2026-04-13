// internal/batch/config.go
package batch

import (
	"encoding/json"
	"os"

	"github.com/yourname/claude-task-manager/pkg/models"
)

// LoadConfig loads a batch configuration from a file
func LoadConfig(path string) (*models.BatchConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config models.BatchConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves a batch configuration to a file
func SaveConfig(path string, config *models.BatchConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// DefaultConfig returns a default batch configuration
func DefaultConfig() *models.BatchConfig {
	return &models.BatchConfig{
		Iterations:     1,
		PermissionMode: models.PermissionAuto,
		DelayBetween:   0,
		StopOnError:    false,
		OnComplete:     models.OnCompleteCommit,
	}
}
