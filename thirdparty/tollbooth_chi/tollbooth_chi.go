package tollbooth_chi

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/config"
	"net/http"
)

func LimitHandler(limiter *config.Limiter) func(http.Handler) http.Handler {
	wrapper := &limiterWrapper{
		limiter: limiter,
	}

	return func(handler http.Handler) http.Handler {
		wrapper.handler = handler
		return wrapper
	}
}

type limiterWrapper struct {
	limiter *config.Limiter
	handler http.Handler
}

func (l *limiterWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		http.Error(w, "Context was canceled", http.StatusServiceUnavailable)
		return

	default:
		httpError := tollbooth.LimitByRequest(l.limiter, r)
		if httpError != nil {
			w.Header().Add("Content-Type", l.limiter.MessageContentType)
			w.WriteHeader(httpError.StatusCode)
			w.Write([]byte(httpError.Message))
			return
		}

		l.handler.ServeHTTP(w, r)
	}
}