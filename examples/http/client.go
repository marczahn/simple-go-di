package http

import (
	"github.com/marczahn/simple-go-di/pkg/di"
	"net/http"
	"time"
)

var httpSingleton = di.NewSingleton[*http.Client]()

func HttpClient() *http.Client {
	return httpSingleton.GetOrSet(
		func() *http.Client { return &http.Client{Timeout: 10 * time.Second} },
		false,
	)
}
