package pagination

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Cobalt0s/creme-brulee/pkg/rest/messaging"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/google/uuid"
)

const (
	PageTimeFormat     = "2006-01-02 15:04:05.000000"
	defaultPageSize    = 3
	defaultMaxPageSize = 100
)

type PageCursor struct {
	Num  uuid.UUID
	Time time.Time
}

type Summary struct {
	Current    *PageCursor
	Next       *PageCursor
	Size       int
	NumResults int
}

type Pagination struct {
	Current    *string `json:"current"`
	Next       *string `json:"next"`
	Size       int     `json:"size"`
	NumResults int     `json:"numResults"`
}

type QP struct {
	PageNum  *string `form:"pageNum"`
	PageSize *int    `form:"pageSize"`
}

func OptionalStringToPage(ctx context.Context, fieldName string, optional *string) (*PageCursor, error) {
	log := ctxlogrus.Extract(ctx)
	if optional != nil {
		pageCursor := &PageCursor{}
		invalidFieldErr := messaging.InvalidField{
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

		parsedTime, err := time.Parse(PageTimeFormat, split[0])
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

func ResolvePageSize(size *int) int {
	if size == nil {
		return defaultPageSize
	}
	if *size < 1 {
		// Negative page size or zero will result into usage of default page size
		return defaultPageSize
	}
	if *size > defaultMaxPageSize {
		return defaultMaxPageSize
	}
	return *size
}

func formatPageCursor(pageCursor *PageCursor) *string {
	if pageCursor != nil {
		formattedTime := pageCursor.Time.Format(PageTimeFormat)
		stringCursor := []byte(
			fmt.Sprintf("%v|%v", formattedTime, pageCursor.Num.String()),
		)
		result := base64.StdEncoding.EncodeToString(stringCursor)
		return &result
	}
	return nil
}

func FromPageSummary(pageSummary *Summary) *Pagination {
	return &Pagination{
		Current:    formatPageCursor(pageSummary.Current),
		Next:       formatPageCursor(pageSummary.Next),
		Size:       pageSummary.Size,
		NumResults: pageSummary.NumResults,
	}
}

type Page struct {
	Cursor *PageCursor
	Size   int
}

func CreatePage(ctx context.Context, qp QP) (*Page, error) {
	pageCursor, err := OptionalStringToPage(ctx, "pageNum", qp.PageNum)
	if err != nil {
		return nil, err
	}
	return &Page{
		Cursor: pageCursor,
		Size:   ResolvePageSize(qp.PageSize),
	}, nil
}

func (page *Page) MakeQueryWithCol(db *gorm.DB, table schema.Tabler, timeColName, identifierCol string) *gorm.DB {
	timeCol := fmt.Sprintf("%v.%v", table.TableName(), timeColName)
	idCol := fmt.Sprintf("%v.%v", table.TableName(), identifierCol)
	resultDB := db.
		Limit(page.Size + 1).
		Order(fmt.Sprintf("%v DESC, %v DESC", timeCol, idCol))
	if page.Cursor != nil {
		t := page.Cursor.Time.Format(PageTimeFormat)
		resultDB = resultDB.
			Where(fmt.Sprintf("(%v < ? OR ( %v = ? AND %v <= ? ))", timeCol, timeCol, idCol), t, t, page.Cursor.Num)
	}
	return resultDB
}
