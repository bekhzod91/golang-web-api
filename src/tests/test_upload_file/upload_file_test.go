package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/hzmat24/api/infrastructure/api/dto"
	"github.com/hzmat24/api/tests"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	t.Skip("File upload")
	testApp := tests.NewTestApp(t)
	testApp.LoadFixtureTenant([]string{
		"fixtures/TestUploadFile/user.json",
	})

	token := testApp.Authenticate("admin@example.com")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "myfile.txt")
	_, _ = io.Copy(part, strings.NewReader("Hello world"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/file/upload/", body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Tenant", "test_tenant")

	rr := testApp.ExecuteRequest(req)

	assert.Equal(t, rr.Code, http.StatusOK)

	response := dto.UploadFileResponseDTO{}
	tests.BindJSON(rr, &response)

	assert.True(t, strings.Contains(response.FileURL, "test_tenant"))
	assert.True(t, strings.Contains(response.FileURL, "myfile.txt"))
}
