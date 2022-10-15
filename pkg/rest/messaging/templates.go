package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
)

var (
	ErrFlatteningFailed = errors.New("failed flattening event")
)

type EventTemplate struct {
	Template interface{}
}

type EventsTemplate struct {
	Template interface{}
	Flatten  func(interface{}) []interface{}
}

func NewEventFromTemplate(ctx context.Context, data []byte, eventTemplates map[string]EventTemplate) (interface{}, error) {
	events, err := NewEventsFromTemplate(ctx, data, templateToTemplates(eventTemplates))
	if err != nil {
		return nil, err
	}
	return events[0], nil
}

func NewEventsFromTemplate(ctx context.Context, data []byte, eventTemplates map[string]EventsTemplate) ([]interface{}, error) {
	log := ctxlogrus.Extract(ctx)

	var objMap map[string]interface{}
	err := json.Unmarshal(data, &objMap)
	if err != nil {
		return nil, err
	}

	body := string(data)
	log.Debugf("incoming data from kafka %v", body)
	if eventType, ok := objMap["type"]; ok {
		event, ok := eventTemplates[eventType.(string)]
		if !ok {
			return nil, fmt.Errorf("unknown event type %v", eventType)
		}
		concreteEvent, err := newEvent(data, eventType, body, event.Template)
		if err != nil {
			return nil, err
		}
		resultingEvents := event.Flatten(concreteEvent)
		if resultingEvents == nil {
			return nil, ErrFlatteningFailed
		}
		return resultingEvents, nil
	}
	return nil, fmt.Errorf("missing key [type] in %v", body)
}

func templateToTemplates(eventTemplates map[string]EventTemplate) map[string]EventsTemplate {
	list := make(map[string]EventsTemplate, len(eventTemplates))
	for k, v := range eventTemplates {
		list[k] = EventsTemplate{
			Template: v.Template,
			Flatten: func(data interface{}) []interface{} {
				return []interface{}{data}
			},
		}
	}
	return list
}

func newEvent(data []byte, eventType interface{}, body string, template interface{}) (interface{}, error) {
	if err := json.Unmarshal(data, template); err != nil {
		return nil, fmt.Errorf(
			"failed to decode userEvent of type %v with value %v because %v", eventType, body, err)
	}
	return template, nil
}
