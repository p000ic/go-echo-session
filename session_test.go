package echosession

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-session/session/v3"
	"github.com/labstack/echo/v4"
)

func TestSession(t *testing.T) {
	cookieName := "test_echo_session"
	e := echo.New()
	e.Use(New(
		session.SetCookieName(cookieName),
		session.SetSign([]byte("sign")),
	))
	e.GET("/", func(ctx echo.Context) error {
		store := FromContext(ctx)
		if ctx.QueryParam("login") == "1" {
			foo, ok := store.Get("foo")
			_, err := fmt.Fprintf(ctx.Response(), "%s:%v", foo, ok)
			if err != nil {
				return err
			}
			return nil
		}

		store.Set("foo", "bar")
		err := store.Save()
		if err != nil {
			t.Error(err)
			return nil
		}
		_, err = fmt.Fprint(ctx.Response(), "ok")
		if err != nil {
			return err
		}
		return nil
	})
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	e.ServeHTTP(w, req)
	res := w.Result()
	cookie := res.Cookies()[0]
	if cookie.Name != cookieName {
		t.Error("Not expected value:", cookie.Name)
		return
	}
	buf, _ := ioutil.ReadAll(res.Body)
	err = res.Body.Close()
	if err != nil {
		return
	}
	if string(buf) != "ok" {
		t.Error("Not expected value:", string(buf))
		return
	}
	req, err = http.NewRequest("GET", "/?login=1", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	e.ServeHTTP(w, req)
	res = w.Result()
	buf, _ = ioutil.ReadAll(res.Body)
	err = res.Body.Close()
	if err != nil {
		return
	}
	if string(buf) != "bar:true" {
		t.Error("Not expected value:", string(buf))
		return
	}
}
