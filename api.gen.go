// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

const (
	DeveloperTokenScopes = "DeveloperToken.Scopes"
)

// BuildPushOptions defines model for BuildPushOptions.
type BuildPushOptions struct {
	Password string `json:"password"`
	Repo     string `json:"repo"`
	Tag      string `json:"tag"`
	Username string `json:"username"`
}

// PublishComponent defines model for PublishComponent.
type PublishComponent struct {
	Description string    `json:"description"`
	Info        *string   `json:"info,omitempty"`
	Name        string    `json:"name"`
	Tags        *[]string `json:"tags,omitempty"`
}

// PublishModuleRequest defines model for PublishModuleRequest.
type PublishModuleRequest struct {
	Components  []PublishComponent `json:"components"`
	Description *string            `json:"description,omitempty"`
	Info        *string            `json:"info,omitempty"`
	Name        string             `json:"name"`
	Version     string             `json:"version"`
}

// PublishModuleResult defines model for PublishModuleResult.
type PublishModuleResult struct {
	Module  *PublishModuleVersion `json:"module,omitempty"`
	Options *BuildPushOptions     `json:"options,omitempty"`
}

// PublishModuleVersion defines model for PublishModuleVersion.
type PublishModuleVersion struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// UpdateModuleVersionRequest defines model for UpdateModuleVersionRequest.
type UpdateModuleVersionRequest struct {
	// Id Module Version ID
	Id string `json:"id"`

	// Repo Image repo
	Repo string `json:"repo"`

	// Tag Image tag
	Tag string `json:"tag"`
}

// PublishModuleJSONRequestBody defines body for PublishModule for application/json ContentType.
type PublishModuleJSONRequestBody = PublishModuleRequest

// UpdateModuleVersionJSONRequestBody defines body for UpdateModuleVersion for application/json ContentType.
type UpdateModuleVersionJSONRequestBody = UpdateModuleVersionRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PublishModuleWithBody request with any body
	PublishModuleWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PublishModule(ctx context.Context, body PublishModuleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateModuleVersionWithBody request with any body
	UpdateModuleVersionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateModuleVersion(ctx context.Context, body UpdateModuleVersionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// HealthCheck request
	HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PublishModuleWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPublishModuleRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PublishModule(ctx context.Context, body PublishModuleJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPublishModuleRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateModuleVersionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateModuleVersionRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateModuleVersion(ctx context.Context, body UpdateModuleVersionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateModuleVersionRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthCheckRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPublishModuleRequest calls the generic PublishModule builder with application/json body
func NewPublishModuleRequest(server string, body PublishModuleJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPublishModuleRequestWithBody(server, "application/json", bodyReader)
}

// NewPublishModuleRequestWithBody generates requests for PublishModule with any type of body
func NewPublishModuleRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/devtools/publish-module")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewUpdateModuleVersionRequest calls the generic UpdateModuleVersion builder with application/json body
func NewUpdateModuleVersionRequest(server string, body UpdateModuleVersionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateModuleVersionRequestWithBody(server, "application/json", bodyReader)
}

// NewUpdateModuleVersionRequestWithBody generates requests for UpdateModuleVersion with any type of body
func NewUpdateModuleVersionRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/devtools/update-module-version")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewHealthCheckRequest generates requests for HealthCheck
func NewHealthCheckRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/health")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PublishModuleWithBodyWithResponse request with any body
	PublishModuleWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PublishModuleResponse, error)

	PublishModuleWithResponse(ctx context.Context, body PublishModuleJSONRequestBody, reqEditors ...RequestEditorFn) (*PublishModuleResponse, error)

	// UpdateModuleVersionWithBodyWithResponse request with any body
	UpdateModuleVersionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateModuleVersionResponse, error)

	UpdateModuleVersionWithResponse(ctx context.Context, body UpdateModuleVersionJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateModuleVersionResponse, error)

	// HealthCheckWithResponse request
	HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error)
}

type PublishModuleResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PublishModuleResult
}

// Status returns HTTPResponse.Status
func (r PublishModuleResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PublishModuleResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateModuleVersionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r UpdateModuleVersionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateModuleVersionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type HealthCheckResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r HealthCheckResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthCheckResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PublishModuleWithBodyWithResponse request with arbitrary body returning *PublishModuleResponse
func (c *ClientWithResponses) PublishModuleWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PublishModuleResponse, error) {
	rsp, err := c.PublishModuleWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePublishModuleResponse(rsp)
}

func (c *ClientWithResponses) PublishModuleWithResponse(ctx context.Context, body PublishModuleJSONRequestBody, reqEditors ...RequestEditorFn) (*PublishModuleResponse, error) {
	rsp, err := c.PublishModule(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePublishModuleResponse(rsp)
}

// UpdateModuleVersionWithBodyWithResponse request with arbitrary body returning *UpdateModuleVersionResponse
func (c *ClientWithResponses) UpdateModuleVersionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateModuleVersionResponse, error) {
	rsp, err := c.UpdateModuleVersionWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateModuleVersionResponse(rsp)
}

func (c *ClientWithResponses) UpdateModuleVersionWithResponse(ctx context.Context, body UpdateModuleVersionJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateModuleVersionResponse, error) {
	rsp, err := c.UpdateModuleVersion(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateModuleVersionResponse(rsp)
}

// HealthCheckWithResponse request returning *HealthCheckResponse
func (c *ClientWithResponses) HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error) {
	rsp, err := c.HealthCheck(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthCheckResponse(rsp)
}

// ParsePublishModuleResponse parses an HTTP response from a PublishModuleWithResponse call
func ParsePublishModuleResponse(rsp *http.Response) (*PublishModuleResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PublishModuleResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PublishModuleResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateModuleVersionResponse parses an HTTP response from a UpdateModuleVersionWithResponse call
func ParseUpdateModuleVersionResponse(rsp *http.Response) (*UpdateModuleVersionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateModuleVersionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseHealthCheckResponse parses an HTTP response from a HealthCheckWithResponse call
func ParseHealthCheckResponse(rsp *http.Response) (*HealthCheckResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthCheckResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /devtools/publish-module)
	PublishModule(w http.ResponseWriter, r *http.Request)

	// (POST /devtools/update-module-version)
	UpdateModuleVersion(w http.ResponseWriter, r *http.Request)

	// (GET /health)
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PublishModule operation middleware
func (siw *ServerInterfaceWrapper) PublishModule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, DeveloperTokenScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PublishModule(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UpdateModuleVersion operation middleware
func (siw *ServerInterfaceWrapper) UpdateModuleVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, DeveloperTokenScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateModuleVersion(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// HealthCheck operation middleware
func (siw *ServerInterfaceWrapper) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.HealthCheck(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/devtools/publish-module", wrapper.PublishModule).Methods("POST")

	r.HandleFunc(options.BaseURL+"/devtools/update-module-version", wrapper.UpdateModuleVersion).Methods("POST")

	r.HandleFunc(options.BaseURL+"/health", wrapper.HealthCheck).Methods("GET")

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RVwW7rNhD8FWLbQwvoWcbrTbeXBEVdoGiQpOnB9YGRNhYTmWS4ZAIj0L8XJKVIimjX",
	"xYNvlshdzczOjt+hVDutJEpLULwDlTXuePh54URTXTuq/9RWKBneaaM0GiswPnGiN2Uq/9vuNUIBZI2Q",
	"W2gzMKhV8sDybfK9IzSS7zBxGNq9OGGwgmI93MwGCN0HY/tN1ndQD09YWt/+2j00gurLnu6cToVUGhHI",
	"JgEK+ZhmdAB1wBI6C4s7St+IL7gxfD/j2XEc4zrC7A9VuQZv8MUhJdhN5/yB6EeDj1DAD/lwnncmyGeS",
	"zQBnZxDtFQ2lu6Xl6e9nY4on6ESuSci0C6cnKhNb3XcA2gzUsCrHymer1bb/hfd+UGUKWFRnklf4pfqk",
	"cUrWv3TFLU5QHjRhBDvxDMRK1pWy1RVkh9NkWrra8S2yfvEPBE2qxB9lp/A/GiptBoSlM8Lub/1kI8cL",
	"5AbNN2dr//QQnn5VZsctFPD733eQxZT1neLpAKW2VnvoV/iKjRfuTj2j/AjmwyXtaNWmhIOLSra+E3LP",
	"bvfkF3/zk6+iIs+tkHuKLxdC5T+zb9cr31tYvwQwLuqOPvwDy8XXxTK6HiXXAgr4ZbFcLEMs2zqIkVf4",
	"apVqKNfRzV+GBdMqWmSK9wa3giwaYvEm47JiBq0zklhpsEJpBW+IPSrDtKP6H7/6XivuO6yqnnO/ORDH",
	"imQvVLWPUSht9w/AtW5EGSrzJ1KD1Px/RUBv+DCHhP5U92w6KDD2mjUOg/lIK0nRRV+Xy3NBDcGXQHrr",
	"yhKJJr6GYv0+c+N60276P7c19BOGjS8cBu5CLnTz/jKKnfTcY4z0MnXX2ZuwNRNhZ729/RYJJRMjT6TQ",
	"mQZ/JO8Sos7SLZvTOdUM5xhXjbyJSbXFxFR+C8esrLF8ZigrrYS0M/HjrUt/Cb4D+ghlhBU+64F+pjjO",
	"2EiP0HjPhNPpxz5nGItXIQNnmi5AfRRyLRbTOIR20/4bAAD//3TaZtMhCwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
