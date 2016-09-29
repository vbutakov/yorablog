package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRestorePasswordServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitRestorePasswordPageHandler(db, "/home/valya/myprogs/yorablog/templates")

	forms := make([]*url.Values, 0, 2)

	form := &url.Values{}
	form.Add("password", "")
	form.Add("passwordconfirm", "pass")

	forms = append(forms, form)

	form = &url.Values{}
	form.Add("password", "pass")
	form.Add("passwordconfirm", "pass1")

	forms = append(forms, form)

	c := &http.Cookie{}
	c.Name = "SessionID"
	c.Value = "99"

	for _, f := range forms {
		body := strings.NewReader(f.Encode())

		req := httptest.NewRequest(http.MethodPost, "http://localhost/restorepassword/?token=0123456789", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", c.String())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Wrong error code: %v\n", w.Code)
		}

		if !bytes.Contains(w.Body.Bytes(), []byte("</html>")) {
			t.Errorf("Page is not complete")
		}

		if !bytes.Contains(w.Body.Bytes(), []byte("<div class=\"errormessage\">")) {
			t.Errorf("Must be error in this case")
		}
	}

	form = &url.Values{}
	form.Add("password", "pass")
	form.Add("passwordconfirm", "pass")

	body := strings.NewReader(form.Encode())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/restorepassword/?token=0123456789", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", c.String())

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

}
