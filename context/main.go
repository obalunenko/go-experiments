package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/objx"
)

type user struct {
	id   int64
	name string
}

type contextWrap struct {
	context.Context
	user   *user
	id     int64
	logger *log.Entry
}

func newContextWrap(context context.Context) *contextWrap {
	id := rand.Int63()
	l := log.WithFields(log.Fields{
		"traceID": id,
	})

	return &contextWrap{Context: context, id: id, logger: l}
}

func (cw contextWrap) GetUser() *user {
	return cw.user
}

func (cw *contextWrap) SetUser(u user) {
	cw.user = &u
}

func (cw contextWrap) GetID() int64 {
	return cw.id
}

func (cw *contextWrap) SetID(id int64) {
	cw.id = id
}

func (cw *contextWrap) LogInfo(msg ...interface{}) {
	cw.logger.Info(messageFromMsgs(msg...))
}

func (cw *contextWrap) LogDebug(msg ...interface{}) {
	cw.logger.Debug(messageFromMsgs(msg...))
}
func (cw *contextWrap) LogError(msg ...interface{}) {
	cw.logger.Error(messageFromMsgs(msg...))
}

func messageFromMsgs(msgs ...interface{}) string {
	if len(msgs) == 0 {
		return ""
	}

	if len(msgs) == 1 {
		m := msgs[0]
		if msgAsStr, ok := m.(string); ok {
			return msgAsStr
		}

		return fmt.Sprintf("%v", m)
	}

	if len(msgs) > 1 {
		var msg string

		for _, m := range msgs {
			if msg != "" {
				msg += " "
			}

			if msgAsStr, ok := m.(string); ok {
				msg += msgAsStr
				continue
			}

			msg += fmt.Sprintf("%v", m)
		}

		return msg
	}
	return ""
}

func (cw *contextWrap) SetLoggerFields(fields log.Fields) {
	cw.logger = cw.logger.WithFields(fields)
}

var users = map[cred]user{
	{
		username: "user1",
		password: "password1",
	}: {
		id:   1,
		name: "user1",
	},
	{
		username: "user2",
		password: "password2",
	}: {
		id:   2,
		name: "user2",
	},
}

type cred struct {
	username string
	password string
}

type middleware func(http.Handler) http.Handler

// chainMiddleware provides syntactic sugar to create a new middleware
// which will be the result of chaining the ones received as parameters.
func chainMiddleware(mw ...middleware) middleware {
	return func(final http.Handler) http.Handler {
		last := final
		for i := len(mw) - 1; i >= 0; i-- {
			last = mw[i](last)
		}

		return last
	}
}

type traceHandler struct {
	next http.Handler
}

func (t *traceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, ok := r.Context().(*contextWrap)
	if !ok {
		ctx = newContextWrap(r.Context())
	}

	ctx.SetLoggerFields(log.Fields{
		"method": r.Method,
		"url":    r.RequestURI,
	})

	r = r.WithContext(ctx)

	ctx.LogDebug("Trace request")

	// log.Printf("Trace request [%d]: [%s]%s \n", ctx.GetID(), r.Method, r.RequestURI)
	t.next.ServeHTTP(w, r)
}

func withTracing(next http.Handler) http.Handler {
	return &traceHandler{next: next}

}

func withLogin(next http.Handler) http.Handler {
	return &authHandler{next: next}
}

type authHandler struct {
	next http.Handler
}

func logingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().(*contextWrap)

	ctx.LogDebug("login reached")
	switch r.Method {
	case http.MethodPost:
		unm := r.FormValue("username")
		pswd := r.FormValue("password")

		u, ok := users[cred{
			username: unm,
			password: pswd,
		}]
		if !ok {
			http.Error(w, "invalid credentials", http.StatusBadRequest)
			ctx.LogError("invalid credentials")
			return
		}

		ctx.SetUser(u)

		r = r.WithContext(ctx)

		_, ok = r.Context().(*contextWrap)
		if !ok {
			http.Error(w, "invalid context login", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "not supported method", http.StatusMethodNotAllowed)
		ctx.LogError("not supported method")
		return
	}

	ctx, ok := r.Context().(*contextWrap)
	if !ok {
		http.Error(w, "invalid context loginHandler", http.StatusInternalServerError)
		return
	}

	authCookieValue := objx.New(map[string]interface{}{
		"user_id": ctx.GetUser().id,
		"name":    ctx.GetUser().name,
	}).MustBase64()

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().(*contextWrap)

	ctx.LogDebug("login reached")
	switch r.Method {
	case http.MethodPost:
		unm := r.FormValue("username")
		pswd := r.FormValue("password")

		u, ok := users[cred{
			username: unm,
			password: pswd,
		}]
		if !ok {
			http.Error(w, "invalid credentials", http.StatusBadRequest)
			ctx.LogError("invalid credentials")
			return
		}

		ctx.SetUser(u)

		r = r.WithContext(ctx)

		_, ok = r.Context().(*contextWrap)
		if !ok {
			http.Error(w, "invalid context login", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "not supported method", http.StatusMethodNotAllowed)
		ctx.LogError("not supported method")
		return
	}

	ctx, ok := r.Context().(*contextWrap)
	if !ok {
		http.Error(w, "invalid context loginHandler", http.StatusInternalServerError)
		return
	}

	authCookieValue := objx.New(map[string]interface{}{
		"user_id": ctx.GetUser().id,
		"name":    ctx.GetUser().name,
	}).MustBase64()

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	a.next.ServeHTTP(w, r)
}

func home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().(*contextWrap)
	ctx.LogDebug("reached home ")

	us := ctx.GetUser()
	if us == nil {
		http.Error(w, "not logged in user", http.StatusUnauthorized)
		ctx.LogError("not logged in user")
		return
	}

	ctx.LogInfo("logged in ", us.name)

	if _, err := fmt.Fprintf(w, "welcome %s", us.name); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		ctx.LogError("failed to write response")
		return
	}
	// http.Redirect()
}

var (
	logLevel = flag.String("log_level", "DEBUG", "Set level of output logs")
)

func setLogger() {
	l, err := log.ParseLevel(*logLevel)
	if err != nil {
		l = log.InfoLevel
	}

	log.SetLevel(l)

	formatter := &log.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          true,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    true,
		QuoteEmptyFields:          true,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	}

	log.SetFormatter(formatter)
}

func main() {
	flag.Parse()

	setLogger()

	rand.Seed(time.Now().UTC().UnixNano())

	mw := chainMiddleware(withTracing, withLogin)
	mux := http.NewServeMux()

	mux.Handle("/", mw(http.HandlerFunc(home)))
	mux.HandleFunc("/login", logingHandler)

	if err := http.ListenAndServe("", mux); err != nil {
		log.Fatal(err)
	}
}
