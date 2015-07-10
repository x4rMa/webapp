package context

// 20.02.15 18:43
// (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

import (
	"encoding/gob"
	"encoding/hex"
	"errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Context struct {
	Input      *Input
	Output     *Output
	XsrfToken string

	user       User
}

type Session struct {
	response http.ResponseWriter
	request  *http.Request
	session  *sessions.Session
	*sessions.Session
}

type User struct {
	Name   string
	Gid    int
	Uid    int
	Active int

	Login string
	Pass  string
	Email string
	Phone string
}

var (
	store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))

	ErrNoSession  = errors.New("No value of Session. User should be exists")
	ErrNoTypeUser = errors.New("Type mismatch for User")
)

const (
	GROUP_ADMIN = 1
)

// Context возвращает текущий контекст
func New(w http.ResponseWriter, r *http.Request) *Context {

	ctx := Context{}
	ctx.Input = NewInput(r)
	ctx.Output = NewOutput(w)

	// get User
	if session, err := ctx.Session("User"); err == nil {
		user, _ := auth(session)
		if user != nil {
			ctx.user = *user
		}
	}

	ctx.initXsrfToken()

	return &ctx
}

// Request возвращает *http.Request
func (c *Context) Request() *http.Request {
	return c.Input.Request()
}

// Response возвращает http.ResponseWriter
func (c *Context) Response() http.ResponseWriter {
	return c.Output.Response()
}

// RouteVars возвращает переменные, заданные в Routes, ex: "/articles/{category}/{id:[0-9]+}"
func (c *Context) RouteVars() map[string]string {
	return mux.Vars(c.Request())
}

// NotFound sends page with 404 http code from template tpls/404.tpl
func (c *Context) NotFound() {
	http.Error(c.Response(), "Not found", http.StatusNotFound)
}

// Forbidden sends response with 403 code
func (c *Context) Forbidden(msg ...string) {
	data := ""
	if len(msg) > 0 {
		log.Println(msg[0])
		data = msg[0]
	}

	http.Error(c.Response(), data, http.StatusForbidden)
}

// InternalError sends response with 501 code
func (c *Context) InternalError(msg ...string) {
	data := ""
	if len(msg) > 0 {
		log.Println(msg[0])
		data = msg[0]
	}

	http.Error(c.Response(), data, http.StatusInternalServerError)
}

// Redirect sends http redirect with 301 code
func (c *Context) Redirect(url string) {
	http.Redirect(c.Output.Response(), c.Input.Request(), url, http.StatusMovedPermanently)
}

// Redirect303 sends http redirect with 303 code
func (c *Context) Redirect303(url string) {
	http.Redirect(c.Output.Response(), c.Input.Request(), url, http.StatusSeeOther)
}

// MethodNotAllowed sends http redirect to root with 405 code
func (c *Context) MethodNotAllowed() {
	http.Redirect(c.Output.Response(), c.Input.Request(), "/", http.StatusMethodNotAllowed)
}

// RedirectTemp send http redirect with 307 code
func (c *Context) RedirectTemp(url string) {
	http.Redirect(c.Output.Response(), c.Input.Request(), url, http.StatusTemporaryRedirect)
}

// SendJSON sends json-content (data)
func (c *Context) SendJSON(data string) int {
	c.Output.Header("Content-Type", "application/json; charset=utf-8")
	return c.RawSend(data)
}

// SendXML sends xml-content (data)
func (c *Context) SendXML(data string) int {
	c.Output.Header("Content-Type", "application/xml; charset=utf-8")
	return c.RawSend(data)
}

// SendJSONP sends jsonp-content (data)
func (c *Context) SendJSONP(data string) int {
	c.Output.Header("Content-Type", "application/javascript; charset=utf-8")
	return c.RawSend(data)
}

// SendHTML sends content (data)
func (c *Context) SendHTML(data string) int {
	c.Output.Header("Content-Type", "text/html; charset=utf-8")
	return c.RawSend(data)
}

// RawSend sends content (data)
func (c *Context) RawSend(data string) int {
	c.Output.Response().WriteHeader(c.Output.Status)

	res, err := c.Output.Response().Write([]byte(data))
	if err != nil {
		log.Println(err)
	}

	return res
}

// GetCookie return cookie from request by a given key.
func (c *Context) GetCookie(key string) string {
	return c.Input.Cookie(key)
}

// SetCookie set cookie for response.
func (c *Context) SetCookie(name string, value string, others ...interface{}) {
	c.Output.Cookie(name, value, others...)
}

// initXsrfToken Генерирует случайный token (для предотвращения csrf)
// Сохраняет в Cookie "XSRF-TOKEN" (для совместимости с AngularJS)
func (c *Context) initXsrfToken() {
	token := c.GetCookie("XSRF-TOKEN")
	if token != "" {
		c.XsrfToken = token
		return
	}

	// Create new token
	key := securecookie.GenerateRandomKey(32)
	c.XsrfToken = hex.EncodeToString(key)
	c.Output.Cookie("XSRF-TOKEN", c.XsrfToken)
}

// GetXsrfToken возвращает token из Form & Headers
func (c *Context) GetXsrfToken() string {
	token := c.Input.Query("XSRF-TOKEN")

	if token == "" {
		token = c.Input.Header("X-XSRF-TOKEN")
	}

	if token == "" {
		token = c.Input.Header("X-CSRF-TOKEN")
	}

	return token
}

// CheckXsrfToken проверяет token
func (c *Context) CheckXsrfToken() bool {
	return c.XsrfToken == c.GetXsrfToken()
}

// User возвращает текущего пользователя
func (c *Context) User() User {
	return c.user
}

// Session открыват сессию
func (c *Context) Session(name string) (*Session, error) {
	session, err := store.Get(c.Input.Request(), name)
	return &Session{c.Output.Response(), c.Input.Request(), session, session}, err
}

// Clear очищает открытую сессию
func (s *Session) Clear() {
	if s != nil {
		s.Options.MaxAge = -1
		s.Save()
	}
}

// Save сохраняет сессию
func (s *Session) Save() {
	(*sessions.Session)(s.session).Save(s.request, s.response)
}

// Exist вернет true, если пользователь авторизован
func (user *User) Exist() bool {
	if user == nil {
		return false
	}

	return user.Uid != 0
}

func (user *User) IsAdmin() bool {
	return user.Gid == GROUP_ADMIN
}

func init() {
	gob.Register(&User{})
}

// Auth возвращает пользователя по сессии
func auth(session *Session) (*User, error) {

	if session == nil {
		return nil, ErrNoSession
	}

	if session.IsNew {
		return nil, nil
	}

	user, ok := (session.Values["User"]).(*User)
	if !ok {
		return user, ErrNoTypeUser
	}

	return user, nil
}


