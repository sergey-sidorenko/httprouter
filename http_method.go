package netrouter

type HttpMethod string

const (
	HttpGetMethod     HttpMethod = "GET"
	HttpPostMethod    HttpMethod = "POST"
	HttpPutMethod     HttpMethod = "PUT"
	HttpPatchMethod   HttpMethod = "PATCH"
	HttpDeleteMethod  HttpMethod = "DELETE"
	HttpHeadMethod    HttpMethod = "HEAD"
	HttpOptionsMethod HttpMethod = "OPTIONS"
)
