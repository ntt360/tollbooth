## Tollbooth

基于滴滴团队的限流器 <github.com/didip/tollbooth>，扩展添加支持基于CookieKey限流


## Five Minute Tutorial
```go
package main

import (
    "github.com/didip/tollbooth"
    "net/http"
)

func HelloHandler(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("Hello, World!"))
}

func main() {
    // Create a request limiter per handler.
    http.Handle("/", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(1, nil), HelloHandler))
    http.ListenAndServe(":12345", nil)
}
```

## Features

1. Rate-limit by request's remote IP, path, methods, custom headers, & basic auth usernames.
    ```go
    import (
        "time"
        "github.com/didip/tollbooth/limiter"
    )

    lmt := tollbooth.NewLimiter(1, nil)

    // or create a limiter with expirable token buckets
    // This setting means:
    // create a 1 request/second limiter and
    // every token bucket in it will expire 1 hour after it was initially set.
    lmt = tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

    // Configure list of places to look for IP address.
    // By default it's: "RemoteAddr", "X-Forwarded-For", "X-Real-IP"
    // If your application is behind a proxy, set "X-Forwarded-For" first.
    lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})

    // Limit only GET and POST requests.
    lmt.SetMethods([]string{"GET", "POST"})

    // Limit based on basic auth usernames.
    // You add them on-load, or later as you handle requests.
    lmt.SetBasicAuthUsers([]string{"bob", "jane", "didip", "vip"})
    // You can remove them later as well.
    lmt.RemoveBasicAuthUsers([]string{"vip"})

    // Limit request headers containing certain values.
    // You add them on-load, or later as you handle requests.
    lmt.SetHeader("X-Access-Token", []string{"abc123", "xyz098"})
    // You can remove all entries at once.
    lmt.RemoveHeader("X-Access-Token")
    // Or remove specific ones.
    lmt.RemoveHeaderEntries("X-Access-Token", []string{"limitless-token"})

    // By the way, the setters are chainable. Example:
    lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
        SetMethods([]string{"GET", "POST"}).
        SetBasicAuthUsers([]string{"sansa"}).
        SetBasicAuthUsers([]string{"tyrion"})
    ```

2. Compose your own middleware by using `LimitByKeys()`.

3. Header entries and basic auth users can expire over time (to conserve memory).

    ```go
    import "time"

    lmt := tollbooth.NewLimiter(1, nil)

    // Set a custom expiration TTL for token bucket.
    lmt.SetTokenBucketExpirationTTL(time.Hour)

    // Set a custom expiration TTL for basic auth users.
    lmt.SetBasicAuthExpirationTTL(time.Hour)

    // Set a custom expiration TTL for header entries.
    lmt.SetHeaderEntryExpirationTTL(time.Hour)
    ```

4. Upon rejection, the following HTTP response headers are available to users:

    * `X-Rate-Limit-Limit` The maximum request limit.

    * `X-Rate-Limit-Duration` The rate-limiter duration.

    * `X-Rate-Limit-Request-Forwarded-For` The rejected request `X-Forwarded-For`.

    * `X-Rate-Limit-Request-Remote-Addr` The rejected request `RemoteAddr`.


5. Customize your own message or function when limit is reached.

    ```go
    lmt := tollbooth.NewLimiter(1, nil)

    // Set a custom message.
    lmt.SetMessage("You have reached maximum request limit.")

    // Set a custom content-type.
    lmt.SetMessageContentType("text/plain; charset=utf-8")

    // Set a custom function for rejection.
    lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) { fmt.Println("A request was rejected") })
    ```

6. Tollbooth does not require external storage since it uses an algorithm called [Token Bucket](http://en.wikipedia.org/wiki/Token_bucket) [(Go library: golang.org/x/time/rate)](https://godoc.org/golang.org/x/time/rate).
