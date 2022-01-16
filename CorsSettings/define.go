package CorsSettings

import (
	"github.com/lamgor666/goboot-common/util/castx"
	"time"
)

type Settings struct {
	allowedOrigins   []string
	allowedHeaders   []string
	allowedMethods   []string
	allowCredentials bool
	exposedHeaders   []string
	maxAge           time.Duration
}

func New(settings map[string]interface{}) *Settings {
	allowedOrigins := []string{"*"}

	if a1 := castx.ToStringSlice(settings["allowedOrigins"]); len(a1) > 0 {
		allowedOrigins = a1
	}

	allowedHeaders := []string{
		"Content-Type",
		"Content-Length",
		"Authorization",
		"Accept",
		"Accept-Encoding",
		"X-Requested-With",
	}

	if a1 := castx.ToStringSlice(settings["allowedHeaders"]); len(a1) > 0 {
		allowedHeaders = a1
	}

	allowedMethods := []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}

	if a1 := castx.ToStringSlice(settings["allowedMethods"]); len(a1) > 0 {
		allowedMethods = a1
	}

	exposedHeaders := []string{
		"Content-Length",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Cache-Control",
		"Content-Language",
		"Content-Type",
	}

	if a1 := castx.ToStringSlice(settings["exposedHeaders"]); len(a1) > 0 {
		exposedHeaders = a1
	}

	var maxAge time.Duration

	if d1, ok := settings["maxAge"].(time.Duration); ok && d1 >= 0 {
		maxAge = d1
	} else if s1, ok := settings["maxAge"].(string); ok && s1 != "" {
		maxAge = castx.ToDuration(s1)
	}

	return &Settings{
		allowedOrigins:   allowedOrigins,
		allowedHeaders:   allowedHeaders,
		allowedMethods:   allowedMethods,
		allowCredentials: castx.ToBool(settings["allowCredentials"]),
		exposedHeaders:   exposedHeaders,
		maxAge:           maxAge,
	}
}

func (st *Settings) AllowedOrigins() []string {
	return st.allowedOrigins
}

func (st *Settings) AllowedHeaders() []string {
	return st.allowedHeaders
}

func (st *Settings) AllowedMethods() []string {
	return st.allowedMethods
}

func (st *Settings) AllowCredentials() bool {
	return st.allowCredentials
}

func (st *Settings) ExposedHeaders() []string {
	return st.exposedHeaders
}

func (st *Settings) MaxAge() time.Duration {
	return st.maxAge
}