package common

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DefaultPageLimit = int64(20)
	MaxPageLimit     = int64(50)
)

type PageRequest struct {
	Limit    int64
	AfterID  primitive.ObjectID
	HasAfter bool
}

func ParsePageRequest(c *gin.Context) (PageRequest, error) {
	limit := DefaultPageLimit
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil || parsedLimit < 1 {
			return PageRequest{}, fmt.Errorf("limit must be a positive integer")
		}
		if parsedLimit > MaxPageLimit {
			parsedLimit = MaxPageLimit
		}
		limit = parsedLimit
	}

	page := PageRequest{Limit: limit}
	if rawAfter := c.Query("after"); rawAfter != "" {
		afterID, err := primitive.ObjectIDFromHex(rawAfter)
		if err != nil {
			return PageRequest{}, fmt.Errorf("after must be a valid object id")
		}
		page.AfterID = afterID
		page.HasAfter = true
	}

	return page, nil
}

func ApplyCursorPage[T interface{ GetID() primitive.ObjectID }](items []T, limit int64) ([]T, gin.H) {
	hasMore := int64(len(items)) > limit
	if hasMore {
		items = items[:limit]
	}

	var nextCursor string
	if hasMore && len(items) > 0 {
		nextCursor = items[len(items)-1].GetID().Hex()
	}

	return items, gin.H{
		"limit":       limit,
		"has_more":    hasMore,
		"next_cursor": nextCursor,
	}
}
