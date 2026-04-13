package claude

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	"github.com/yourname/claude-task-manager/pkg/models"
)

type OutputParser struct {
	taskID string
}

func NewOutputParser(taskID string) *OutputParser {
	return &OutputParser{taskID: taskID}
}

func (p *OutputParser) Parse(reader io.Reader) []*models.Event {
	var events []*models.Event
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "{") {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		event := p.parseLine(line, raw)
		if event != nil {
			events = append(events, event)
		}
	}

	return events
}

func (p *OutputParser) parseLine(line string, raw map[string]interface{}) *models.Event {
	eventType := p.detectEventType(raw)
	if eventType == "" {
		return nil
	}

	return models.NewEvent(p.taskID, eventType, json.RawMessage(line))
}

func (p *OutputParser) detectEventType(raw map[string]interface{}) models.EventType {
	if _, ok := raw["tool"]; ok {
		if _, ok := raw["result"]; ok {
			return models.EventToolResult
		}
		return models.EventToolUse
	}
	if _, ok := raw["message"]; ok {
		return models.EventMessage
	}
	if _, ok := raw["error"]; ok {
		return models.EventError
	}
	if _, ok := raw["permission"]; ok {
		return models.EventPermission
	}
	return ""
}
