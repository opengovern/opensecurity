package utils

import (
	"fmt"
	"strconv"
)

const DefaultPageSize = int64(20)

func PageConfigFromStrings(page, size string) (pageNumber int64, pageSize int64, err error) {
	pageSize = DefaultPageSize
	if size != "" {
		pageSize, err = strconv.ParseInt(size, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("pageSize is not a valid integer")
		}
	}
	pageNumber = int64(1)
	if page != "" {
		pageNumber, err = strconv.ParseInt(page, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("pageNumber is not a valid integer")
		}
	}
	return pageNumber, pageSize, nil
}
