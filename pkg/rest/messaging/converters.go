package messaging

import (
	"context"
	"encoding/base64"
	"github.com/Cobalt0s/creme-brulee/pkg/rest/pagination"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"strings"
	"time"
)

func OptionalStringToPage(ctx context.Context, fieldName string, optional *string) (*pagination.PageCursor, error) {
	log := ctxlogrus.Extract(ctx)
	if optional != nil {
		pageCursor := &pagination.PageCursor{}
		invalidFieldErr := InvalidField{
			Name:   fieldName,
			Format: "page",
		}

		data, err := base64.StdEncoding.DecodeString(*optional)
		if err != nil {
			log.Debugf("field %v is not in base64 format", fieldName)
			return nil, invalidFieldErr
		}

		split := strings.Split(string(data), "|")
		if len(split) != 2 {
			log.Debugf("field %v is not in [timestamp|uuid] format", fieldName)
			return nil, invalidFieldErr
		}

		parsedTime, err := time.Parse(pagination.PageTimeFormat, split[0])
		if err != nil {
			log.Debugf("field %v has invalid timestamp format [%v]", fieldName, split[0])
			return nil, invalidFieldErr
		}
		log.Debugf("page time %v parsed as %v", split[0], parsedTime)
		pageCursor.Time = parsedTime

		pageNum, err := uuid.Parse(split[1])
		if err != nil {
			log.Debugf("field %v is not in uuid format", fieldName)
			return nil, invalidFieldErr
		}
		pageCursor.Num = pageNum

		return pageCursor, nil
	}
	return nil, nil
}

func OptionalStringToUUID(ctx context.Context, fieldName string, optional *string) (*uuid.UUID, error) {
	log := ctxlogrus.Extract(ctx)
	if optional != nil {
		invalidFieldErr := InvalidField{
			Name:   fieldName,
			Format: "uuid",
		}

		result, err := uuid.Parse(*optional)
		if err != nil {
			log.Debugf("field %v is not in uuid format", fieldName)
			return nil, invalidFieldErr
		}
		return &result, nil
	}
	return nil, nil
}

func OptionalStringListToUUIDList(ctx context.Context, fieldName string, stringList []string) ([]uuid.UUID, error) {
	var result []uuid.UUID
	if stringList != nil {
		result = make([]uuid.UUID, len(stringList))
		for i, v := range stringList {
			uuidResult, err := OptionalStringToUUID(ctx, fieldName, &v)
			if err != nil {
				return nil, err
			}
			result[i] = *uuidResult
		}
	}
	return result, nil
}

func OptionalStringToUUIDList(ctx context.Context, fieldName string, text *string) ([]uuid.UUID, error) {
	var stringList []string
	if text != nil {
		stringList = strings.Split(*text, ",")
	}
	return OptionalStringListToUUIDList(ctx, fieldName, stringList)
}

func OptionalStringToTime(ctx context.Context, fieldName string, text *string) (*time.Time, error) {
	if text != nil {
		parseTime, err := ParseNano3339Time(*text)
		if err != nil {
			invalidFieldErr := InvalidField{
				Name:   fieldName,
				Format: timeNanoFormat,
			}
			return nil, invalidFieldErr
		}
		return parseTime, nil
	}
	return nil, nil
}
