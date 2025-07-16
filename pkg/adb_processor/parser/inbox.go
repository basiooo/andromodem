package parser

import (
	"strings"
)

type Inbox struct {
	Messages []Message `json:"messages"`
}

func NewInbox() IParser {
	return &Inbox{}
}

func (i *Inbox) Parse(rawData string) error {
	rawData = strings.TrimSpace(rawData)

	if i.isEmpty(rawData) {
		return nil
	}

	splitMessages := i.splitData(rawData)
	messages := make([]Message, 0, len(splitMessages))
	for _, input := range splitMessages {
		msg := NewMessage()
		if err := msg.Parse(input); err != nil {
			continue
		}
		message := *msg.(*Message)
		if message.IsEmpty() {
			continue
		}
		messages = append(messages, message)
	}

	i.Messages = messages
	return nil
}

func (i *Inbox) isEmpty(rawData string) bool {
	return strings.EqualFold(strings.TrimSpace(rawData), "no result found.")
}

func (i *Inbox) splitData(rawData string) []string {
	rawParts := strings.Split(rawData, "Row:")
	var result []string
	for _, part := range rawParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, "Row:"+part)
	}
	return result
}
