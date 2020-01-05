package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
	"io/ioutil"
	"net/http"
	"proiect_web/mongodb"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "./templates/*.html",
	}

	rnd = renderer.New(opts)
}

func UnauthorizedLoginHandler(res http.ResponseWriter, req *http.Request) {
	//if err := rnd.HTML(res, http.StatusUnauthorized, "login_bad", nil); err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	url := "/login"
	http.Redirect(res, req, url, 302)
	return
}

func LoginHandlerGet(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(Key)
	if err == nil {
		session, err := mongodb.GetSession(cookie.Value)
		if err == nil {
			url := fmt.Sprintf("/profile/%s", session.Username)
			http.Redirect(res, req, url, 302)
			return
		}
	}

	if err := rnd.HTML(res, http.StatusOK, "login_get", nil); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func LoginHandlerPost(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	username := req.Form["username"][0]
	password := req.Form["password"][0]

	if err := mongodb.CheckLogin(username, password); err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	session, err := mongodb.InsertSession(username)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	SetCookie(res, session.Token)
	url := fmt.Sprintf(`/profile/%s`, username)
	http.Redirect(res, req, url, 302)
}

func ProfileHandler(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(Key)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	// Username should match session username
	session, err := mongodb.GetSession(cookie.Value)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	vars := mux.Vars(req)
	username := vars["username"]
	if username != session.Username {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	if err := rnd.HTML(res, http.StatusOK, "profile", username); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func LogoutHandler(res http.ResponseWriter, req *http.Request) {
	UnsetCookie(res)
	url := "/login"
	http.Redirect(res, req, url, 302)
}

func ProfileAddQuoteGet(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(Key)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	// Username should match session username
	session, err := mongodb.GetSession(cookie.Value)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	vars := mux.Vars(req)
	username := vars["username"]
	if username != session.Username {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	if err := rnd.HTML(res, http.StatusOK, "article_add", username); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ProfileAddQuote(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cookie, err := req.Cookie(Key)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	// Username should match session username
	session, err := mongodb.GetSession(cookie.Value)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	username := session.Username
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var link mongodb.Quote
	if err = json.Unmarshal(body, &link); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := mongodb.AddProfileQuote(username, link); err != nil {
		fmt.Println("asdasd", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ProfileGetQuotes(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cookie, err := req.Cookie(Key)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	// Username should match session username
	session, err := mongodb.GetSession(cookie.Value)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	//vars := mux.Vars(req)
	//username := vars["username"]
	//if username != session.Username {
	//	url := "/unauthorized"
	//	http.Redirect(res, req, url, 302)
	//	return
	//}
	username := session.Username

	account, _ := mongodb.GetUser(username)
	if err := json.NewEncoder(res).Encode(account.Quotes); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func ProfileDeleteQuote(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cookie, err := req.Cookie(Key)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	// Username should match session username
	session, err := mongodb.GetSession(cookie.Value)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	//vars := mux.Vars(req)
	//username := vars["username"]
	//if username != session.Username {
	//	url := "/unauthorized"
	//	http.Redirect(res, req, url, 302)
	//	return
	//}
	username := session.Username

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var link mongodb.Quote
	if err = json.Unmarshal(body, &link); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := mongodb.DeleteProfileQuote(username, link.Data); err != nil {
		fmt.Println("asdasd", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ProfileUpdateQuote(res http.ResponseWriter, req *http.Request) {
	fmt.Println("asdasdasd")
	res.Header().Set("Content-Type", "application/json")
	cookie, err := req.Cookie(Key)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	// Username should match session username
	session, err := mongodb.GetSession(cookie.Value)
	if err != nil {
		url := "/unauthorized"
		http.Redirect(res, req, url, 302)
		return
	}

	username := session.Username

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var link mongodb.Quote
	if err = json.Unmarshal(body, &link); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := mongodb.UpdateQuotes(username, link.Data); err != nil {
		fmt.Println("asdasd", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
