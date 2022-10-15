package messaging

import (
	"context"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"strings"
	"time"
)

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
