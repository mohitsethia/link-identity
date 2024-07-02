package infrastructure

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// logMu is a mutex to control changes on watson.LogData
var logMu sync.RWMutex

// NewLogger create a new logger instance
func NewLogger(out io.Writer, logLevel string, env string) *logrus.Logger {
	log := logrus.New()
	log.Out = out
	log.Formatter = getJSONFormatter(env)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.DebugLevel
	}
	log.Level = level

	return log
}

// JSONFormatter Wrapper for logrus.JSONFormatter
type JSONFormatter struct {
	logrus.JSONFormatter
}

// getJSONFormatter
func getJSONFormatter(env string) *JSONFormatter {
	jsonFormatter := logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05-0700",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
		},
		PrettyPrint: env == "dev",
	}
	return &JSONFormatter{jsonFormatter}
}

type loggerMiddleware struct {
	logger *logrus.Entry
}

// Middleware specify a interface to http calls
type Middleware interface {
	Wrap(next http.Handler) http.Handler
}

// NewLoggerMiddleware ...
func NewLoggerMiddleware(logEntry *logrus.Entry) Middleware {
	return &loggerMiddleware{
		logger: logEntry,
	}
}

// Wrap ...
func (lmw *loggerMiddleware) Wrap(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// serve the request
		started := time.Now()
		lr := &LogResponse{
			ResponseWriter: w,
		}

		next.ServeHTTP(lr, r)
		elapsed := time.Since(started)

		// copy log data
		ld := LogData(r.Context())
		ld["status_code"] = lr.StatusCode
		ld["response_time"] = elapsed.Milliseconds()
		ld["request_path"] = r.RequestURI
		ld["remote_addr"] = r.RemoteAddr

		var logLevel logrus.Level
		switch statusCode := lr.StatusCode; {
		case statusCode >= http.StatusInternalServerError:
			logLevel = logrus.ErrorLevel
		case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusInternalServerError:
			logLevel = logrus.WarnLevel
		default:
			logLevel = logrus.InfoLevel
		}
		lmw.logger.WithFields(ld).Logf(logLevel, fmt.Sprintf("[%s] %s", r.Method, r.RequestURI))
	}

	return http.HandlerFunc(fn)
}

// CoreLogger ...
type CoreLogger struct {
	logger *logrus.Entry
}

// Printf ...
func (c *CoreLogger) Printf(format string, v ...interface{}) {
	c.logger.Infof(format, v...)
}

// Logger is the interface used internally to log
type Logger interface {
	Printf(format string, v ...interface{})
}

// RequestLogFormat is the default template used by the logger
var RequestLogFormat = "{{.RemoteAddr}} [{{.Response.Elapsed}}] \"{{.Method}} " +
	"{{.RequestURI}} {{.Proto}}\" {{.Response.StatusCode}} {{.Response.StatusText}} \"{{.UserAgent}}\""

// LogByRequestFunc specify a function that will be called everytime that is necessary
// log something
type LogByRequestFunc func(logReq *LogRequest)

// LogMiddleware is a implementation of Middleware with some additional methods to
// be configured: SetLogger() and SetLoggerFunc()
type LogMiddleware struct {
	fn LogByRequestFunc
}

// NewLogMiddleware create a log middleware
func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

// SetLogger set a fdhttp.Logger to send logs
func (m *LogMiddleware) SetLogger(log Logger) {
	tmpl := template.Must(template.New("log-template").Parse(RequestLogFormat))

	m.fn = func(logReq *LogRequest) {
		var b bytes.Buffer
		tmpl.Execute(&b, logReq)
		log.Printf(b.String())
	}
}

// SetLoggerFunc set a function that is called everytime that need to log
func (m *LogMiddleware) SetLoggerFunc(fn LogByRequestFunc) {
	m.fn = fn
}

// Wrap will be called in every request
func (m *LogMiddleware) Wrap(next http.Handler) http.Handler {
	if m.fn == nil {
		panic("Using LogMiddleware without set a log function (See: SetLogger or SetLoggerFunc)")
	}

	fn := func(w http.ResponseWriter, req *http.Request) {
		started := time.Now()

		lr := &LogResponse{
			ResponseWriter: w,
			req:            req,
		}
		next.ServeHTTP(lr, req)

		lr.Elapsed = time.Since(started)

		logReq := &LogRequest{
			Request:    *req,
			Response:   lr,
			RemoteAddr: getRemoteAddr(req),
		}

		m.fn(logReq)
	}

	return http.HandlerFunc(fn)
}

// LogRequest contain all necessary fields to be logged
type LogRequest struct {
	http.Request
	Response   *LogResponse
	RemoteAddr string
}

// LogResponse it's a wrap to be able read the status code
type LogResponse struct {
	http.ResponseWriter
	req        *http.Request
	StatusCode int
	Elapsed    time.Duration
}

// WriteHeader ...
func (lr *LogResponse) WriteHeader(code int) {
	lr.StatusCode = code
	lr.ResponseWriter.WriteHeader(code)
}

// StatusText ...
func (lr *LogResponse) StatusText() string {
	return http.StatusText(lr.StatusCode)
}

func getRemoteAddr(req *http.Request) string {
	remoteAddr := req.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr, _, _ = net.SplitHostPort(req.RemoteAddr)
	}

	return remoteAddr
}

// LogData return a new LoggerData with all fields.
// The copy is returned to avoid clients change the original values
func LogData(ctx context.Context) map[string]interface{} {
	logMu.RLock()
	defer logMu.RUnlock()

	logData, ok := ctx.Value(ContextKeyLogData).(map[string]interface{})
	if !ok {
		return make(map[string]interface{}, 0)
	}

	d := make(map[string]interface{}, len(logData))
	for k, v := range logData {
		d[k] = v
	}

	return d
}

// AddLogData adds to the log data field that will be kept during the whole request.
// Passing data as nill would delete the corresponding field.
func AddLogData(ctx context.Context, field string, data interface{}) error {
	logMu.Lock()
	defer logMu.Unlock()

	logData, ok := ctx.Value(ContextKeyLogData).(map[string]interface{})
	if !ok {
		logData = make(map[string]interface{})
	}

	if data == nil {
		delete(logData, field)
		return nil
	}
	logData[field] = data

	if !ok {
		return errors.New("LogData was not initialized in this context")
	}
	return nil
}

type contextKey struct {
	name string
}

func (c contextKey) String() string {
	return "pdkit context key " + c.name
}

var (
	// ContextKeyLogData ...
	ContextKeyLogData = &contextKey{"log-data"}

	// ContextKeyPlatform ...
	ContextKeyPlatform = &contextKey{"platform"}

	// ContextKeyCountry ...
	ContextKeyCountry = &contextKey{"country"}

	// ContextKeyRegion ...
	ContextKeyRegion = &contextKey{"region"}

	// ContextKeyRequestID ...
	ContextKeyRequestID = &contextKey{"request-id"}

	// ContextKeyRequestIP ...
	ContextKeyRequestIP = &contextKey{"request-ip"}

	// ContextKeyJWTClaims ...
	ContextKeyJWTClaims = &contextKey{"jwt-claims"}

	// ContextKeyBearerToken ...
	ContextKeyBearerToken = &contextKey{"jwt-token"}
)

const (
	// HeaderAPIOAuthToken is the Oauth Bearer token header name
	HeaderAPIOAuthToken string = "Authorization"

	// LogAttrChannel standarize the log attr in stdout / DD
	LogAttrChannel string = "channel"
	// LogAttrRequest standarize the log attr in stdout / DD
	LogAttrRequest string = "request"
	// LogAttrRequestID standarize the log attr in stdout / DD
	LogAttrRequestID string = "request_id"
	// LogAttrRequestMethod standarize the log attr in stdout / DD
	LogAttrRequestMethod string = "request_method"
	// LogAttrRequestPath standarize the log attr in stdout / DD
	LogAttrRequestPath string = "request_path"
	// LogAttrRequestLength standarize the log attr in stdout / DD
	LogAttrRequestLength string = "request_length"
	// LogAttrURI standarize the log attr in stdout / DD
	LogAttrURI string = "uri"
	// LogAttrUserAgent standarize the log attr in stdout / DD
	LogAttrUserAgent string = "user_agent"
	// LogAttrHTTPReferer standarize the log attr in stdout / DD
	LogAttrHTTPReferer string = "http_referer"
	// LogAttrApplication standarize the log attr in stdout / DD
	LogAttrApplication string = "application"
	// LogAttrCountry standarize the log attr in stdout / DD
	LogAttrCountry string = "country"
	// LogAttrRegion standarize the log attr in stdout / DD
	LogAttrRegion string = "region"
	// LogAttrStatus standarize the log attr in stdout / DD
	LogAttrStatus string = "status"
	// LogAttrPlatform standarize the log attr in stdout / DD
	LogAttrPlatform string = "platform"
	// LogAttrClientIP standarize the log attr in stdout / DD
	LogAttrClientIP string = "client_ip"
)

// StringInSlice checks if the string is in the array
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
