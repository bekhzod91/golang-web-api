package command

import (
	"context"
	"fmt"
	"github.com/myproject/api/domain/exception"
	"github.com/myproject/api/infrastructure/api/dto"
	"github.com/myproject/api/infrastructure/core"
	"path/filepath"
	"strings"
	"time"
)

func UploadFile(c core.IContext) (*dto.UploadFileResponseDTO, error) {
	file, metadata, err := c.Request().FormFile("file")
	if err != nil {
		return nil, exception.New("provide file")
	}
	ext := filepath.Ext(metadata.Filename)
	if ext == "." {
		return nil, exception.New("provide file extension")
	}

	currentYearMonth := time.Now().Format("2006-01")
	key := fmt.Sprintf("%s/%s/%d-%s", c.TenantSchemaName(), currentYearMonth, time.Now().UnixMilli(), strings.ReplaceAll(metadata.Filename, " ", "_"))
	fileURL, err := c.AWS().UploadS3(context.Background(), key, file)

	if err != nil {
		return nil, err
	}

	return &dto.UploadFileResponseDTO{FileURL: fileURL}, nil
}
