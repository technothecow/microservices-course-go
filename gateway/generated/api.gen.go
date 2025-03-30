// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package gen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// CreatePostRequest defines model for CreatePostRequest.
type CreatePostRequest struct {
	Description string   `json:"description"`
	IsPrivate   string   `json:"isPrivate"`
	Tags        []string `json:"tags"`
	Title       string   `json:"title"`
}

// EditPostRequest defines model for EditPostRequest.
type EditPostRequest struct {
	Description string             `json:"description"`
	Id          openapi_types.UUID `json:"id"`
	IsPrivate   string             `json:"isPrivate"`
	Tags        []string           `json:"tags"`
	Title       string             `json:"title"`
}

// EditProfile defines model for EditProfile.
type EditProfile struct {
	DateOfBirth *openapi_types.Date `json:"date_of_birth,omitempty"`
	FullName    *string             `json:"full_name,omitempty"`
	PhoneNumber *string             `json:"phone_number,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PaginatedPostsRequest defines model for PaginatedPostsRequest.
type PaginatedPostsRequest struct {
	Page     float32   `json:"page"`
	Pagesize float32   `json:"pagesize"`
	Tags     *[]string `json:"tags,omitempty"`
}

// PingResponse defines model for PingResponse.
type PingResponse struct {
	Message string `json:"message"`
}

// Post defines model for Post.
type Post struct {
	CreatedAt   string             `json:"createdAt"`
	CreatorId   string             `json:"creatorId"`
	Description string             `json:"description"`
	Id          openapi_types.UUID `json:"id"`
	IsPrivate   string             `json:"isPrivate"`
	Tags        []string           `json:"tags"`
	Title       string             `json:"title"`
	UpdatedAt   string             `json:"updatedAt"`
}

// PostId defines model for PostId.
type PostId struct {
	Id openapi_types.UUID `json:"id"`
}

// PostsList defines model for PostsList.
type PostsList struct {
	Posts []Post `json:"posts"`
}

// ProfileResponse defines model for ProfileResponse.
type ProfileResponse struct {
	DateOfBirth *string            `json:"date_of_birth,omitempty"`
	Email       string             `json:"email"`
	FullName    *string            `json:"full_name,omitempty"`
	Id          openapi_types.UUID `json:"id"`
	LastLogin   string             `json:"last_login"`
	PhoneNumber *string            `json:"phone_number,omitempty"`
	Username    string             `json:"username"`
}

// UserRegistration defines model for UserRegistration.
type UserRegistration struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
	Username string              `json:"username"`
}

// UsernameAndPassword defines model for UsernameAndPassword.
type UsernameAndPassword struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// DeletePostJSONRequestBody defines body for DeletePost for application/json ContentType.
type DeletePostJSONRequestBody = PostId

// GetPostJSONRequestBody defines body for GetPost for application/json ContentType.
type GetPostJSONRequestBody = PostId

// EditPostJSONRequestBody defines body for EditPost for application/json ContentType.
type EditPostJSONRequestBody = EditPostRequest

// CreatePostJSONRequestBody defines body for CreatePost for application/json ContentType.
type CreatePostJSONRequestBody = CreatePostRequest

// GetPostsListJSONRequestBody defines body for GetPostsList for application/json ContentType.
type GetPostsListJSONRequestBody = PaginatedPostsRequest

// RegisterUserJSONRequestBody defines body for RegisterUser for application/json ContentType.
type RegisterUserJSONRequestBody = UserRegistration

// EditMyProfileJSONRequestBody defines body for EditMyProfile for application/json ContentType.
type EditMyProfileJSONRequestBody = EditProfile

// AuthUserJSONRequestBody defines body for AuthUser for application/json ContentType.
type AuthUserJSONRequestBody = UsernameAndPassword

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Ping endpoint
	// (GET /ping)
	Ping(c *gin.Context)
	// Delete a post
	// (DELETE /v1/posts)
	DeletePost(c *gin.Context)
	// Get a post
	// (GET /v1/posts)
	GetPost(c *gin.Context)
	// Edit a post
	// (PATCH /v1/posts)
	EditPost(c *gin.Context)
	// Create a post
	// (POST /v1/posts)
	CreatePost(c *gin.Context)
	// Get paginated list of posts
	// (GET /v1/posts/list)
	GetPostsList(c *gin.Context)
	// Register a new user
	// (POST /v1/users)
	RegisterUser(c *gin.Context)
	// Info about logged user
	// (GET /v1/users/me)
	GetMyProfile(c *gin.Context)
	// Edit user profile
	// (PATCH /v1/users/me)
	EditMyProfile(c *gin.Context)
	// Authenticates user
	// (POST /v1/users/me)
	AuthUser(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// Ping operation middleware
func (siw *ServerInterfaceWrapper) Ping(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Ping(c)
}

// DeletePost operation middleware
func (siw *ServerInterfaceWrapper) DeletePost(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeletePost(c)
}

// GetPost operation middleware
func (siw *ServerInterfaceWrapper) GetPost(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetPost(c)
}

// EditPost operation middleware
func (siw *ServerInterfaceWrapper) EditPost(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.EditPost(c)
}

// CreatePost operation middleware
func (siw *ServerInterfaceWrapper) CreatePost(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreatePost(c)
}

// GetPostsList operation middleware
func (siw *ServerInterfaceWrapper) GetPostsList(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetPostsList(c)
}

// RegisterUser operation middleware
func (siw *ServerInterfaceWrapper) RegisterUser(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RegisterUser(c)
}

// GetMyProfile operation middleware
func (siw *ServerInterfaceWrapper) GetMyProfile(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetMyProfile(c)
}

// EditMyProfile operation middleware
func (siw *ServerInterfaceWrapper) EditMyProfile(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.EditMyProfile(c)
}

// AuthUser operation middleware
func (siw *ServerInterfaceWrapper) AuthUser(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AuthUser(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/ping", wrapper.Ping)
	router.DELETE(options.BaseURL+"/v1/posts", wrapper.DeletePost)
	router.GET(options.BaseURL+"/v1/posts", wrapper.GetPost)
	router.PATCH(options.BaseURL+"/v1/posts", wrapper.EditPost)
	router.POST(options.BaseURL+"/v1/posts", wrapper.CreatePost)
	router.GET(options.BaseURL+"/v1/posts/list", wrapper.GetPostsList)
	router.POST(options.BaseURL+"/v1/users", wrapper.RegisterUser)
	router.GET(options.BaseURL+"/v1/users/me", wrapper.GetMyProfile)
	router.PATCH(options.BaseURL+"/v1/users/me", wrapper.EditMyProfile)
	router.POST(options.BaseURL+"/v1/users/me", wrapper.AuthUser)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xYX2/bthf9KgR/fVRsOS5+wBwEmNsNXbB1NRL0aRgMRryW2EqkSl41dQt/94GkZEsW",
	"5bhpvKXD3myK4j33nPuP+kITVZRKgkRDZ1+oSTIomPv5UgNDWCiD1/ChAoN2sdSqBI0C3BYOJtGiRKGk",
	"/YvrEuiMGtRCpnQTUWEWWnxkCMGnyFJ3ikAoTHiHX2Bas7X7LzAPnbWJqIYPldDA6eyPelvUgdcGU5v+",
	"c2tA3b6DBK2Fn7nAb3OZ2+WV0gVDOqNVJTiNngwzHsw30qPVSnhje9QwhKVaLW+FxswuwCdWlHYnPY/j",
	"52fx9Oz8BxrtyOG1tX3nVlWeLyUrwuSUmZKwlFVxCzrscR+21kr3ASeKh00UYAxLj+DTnbDbH6JswVIh",
	"GQK3YWUG46rs2qvds+6yFIz4HH74tYGyB98ZbZkI4hcyvQZTKmkCmh/N1EGKVIiRxNUfPsehUCKT89k0",
	"DgWQe1XpKx6k5DtL4YhWJX8YEw/L/h17tRtRS4s2miEtPe9dNY8itY93yIT5TQSzyD7q0P5Mw4rO6P/G",
	"uz43rpvc2MXdvTnijgzi8IVwODmOLYg9waFgIg+GwuHSeGTk5szgMlepkF+fWvdU34hWBvQAwFA4brc3",
	"XnfghWh/a0BfQyoMatbkcJf3LX07196pTI64gh/rpVGiinYramz3/WXG3CnNu8ct2LNnd7Hm7SO2O6PD",
	"pHRBLbmCexOhRVLLisc8RJHdP5d80cK/33F2T75BxRC2Pib7lpAr1Ruh6HxxRVZKE8yAvGIId2y9rVIz",
	"Wq+Q+eKKRvQjaOPfmoziUWyBqhIkKwWd0eloMoodCsycg+PS4p19oSlg3+41YKWlIYwYYfUgTYdyZ/rQ",
	"soXMNUBqffZp7o4+j2M/QEgE6Q5nZZmLxL01fmd8VPo6c28VajdYR1QX6JtfHeemKgqm1zUgApKXSkh0",
	"z8YfJ+Nt6eOQg29KXT9+cuuu5nkFweALxdeP54gv/ZtuhKCuYBOmr+9nRJ/Hk/6j3xUSVmEGEi004H7n",
	"dHin0uJzs+15f5uFSqRCslKV5Hv8eqIII2XdH+r46dL5CvApcfloVodC0DIZUO0F46Qm4InI9wqwpV3J",
	"MMn66jX3vBPJt3+N/E/Hr9fRctgWsr4kdHXcfaQ4kZL9ryD/Wi077Hu/t/y3m8w4r2fvQ3XRD+gnKo7B",
	"K/U/oIt38oA4DyDelq+ycZBYqolaEd/dGxnszLW97/RNeO3saCPhjtjNhCWJqiT2hhs/SIO2E+OJ1OrN",
	"60cJFWDOHkRMlSRgjL0HrUl9LaURzYDxmpMbwLOXSr0Xbgbq2olamHdzuJVjieo9yEt2m0zOpxfkF8Ty",
	"jczXF+QGkkrDBblhBdwIhMsb1CLBC7JgmF2OL8hr9ulsnsLl9P9x6Aq+2WXp47QW9y0rEHM263e50A6p",
	"RuRWQHRDaewH/aF8fr1uvvmdMp/2btMBD10AlDWUofx6K9uNp8PDlVwpwm5VhSRXaQq8puLglNB1/kSj",
	"QuPT31zC7qf8QCE7QLRr3dWeWgO1KoPkvbGZzG1FZLkhTHJiAA2RYFOd6TVJXD6bXvmaV5iduHTt36Mf",
	"er2Zt2u+D7zvpGz1y0xLrT3l216aptJsNn8FAAD//wo2OMVnGgAA",
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
