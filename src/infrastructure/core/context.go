package core

import (
	"github.com/myproject/api/pkg/aws"
	"github.com/myproject/api/pkg/multi_tenency"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/go-openapi/strfmt"

	"github.com/myproject/api/domain/entity"
	"github.com/myproject/api/infrastructure/storage"
	"github.com/myproject/api/pkg/logger"
)

type IContext interface {
	Writer() http.ResponseWriter
	Request() *http.Request

	JSON(status int, v any)
	HTML(status int, v string)
	PlainText(status int, v string)
	Data(status int, v []byte)

	BadRequest(v error)
	Forbidden()
	NotFound()
	Unauthorized()
	InternalServerError()

	OK(v any)
	Created(v any)
	NoContent()

	BindJSON(target any) error
	ShouldBindJSON(target interface {
		Validate(formats strfmt.Registry) error
	}) error

	URLParam(key string) string
	QueryParam(key string) string

	TenantSchemaName() string
	Storage() storage.IStorage
	Logger() logger.ILogger
	User() *entity.User

	AWS() *aws.Client
}

type context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) IContext {
	return &context{w: w, r: r}
}

func (c *context) Writer() http.ResponseWriter {
	return c.w
}

func (c *context) Request() *http.Request {
	return c.r
}

func (c *context) URLParam(key string) string {
	return chi.URLParam(c.r, key)
}

func (c *context) QueryParam(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *context) TenantSchemaName() string {
	schemaName := multi_tenency.SchemaNameFromContext(c.r.Context())
	return schemaName
}

func (c *context) Storage() storage.IStorage {
	app := AppFromContext(c)
	tenantPostgresClient := multi_tenency.DBFromContext(c.r.Context())

	return storage.NewStorage(
		app.redisClient,
		app.db.Shared(),
		tenantPostgresClient,
	)
}

func (c *context) Logger() logger.ILogger {
	return httplog.LogEntry(c.Request().Context())
}

func (c *context) User() *entity.User {
	return UserFromContext(c)
}

func (c *context) AWS() *aws.Client {
	app := AppFromContext(c)
	return app.awsClient
}
