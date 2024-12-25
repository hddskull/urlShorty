package middleware

import (
	"errors"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	"net/http"
)

const cookieName = "ShortenerCookie"

func WithJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get cookie
		cookie, err := getShortenerCookie(r)
		if r.Method == http.MethodGet {
			utils.SugaredLogger.Debugln("getShortenerCookie(...) cookie:", cookie, " | err:", err)
			if err != nil {
				var cookieErr *custom.CookieError
				if errors.As(err, &cookieErr) {
					http.Error(w, cookieErr.Error(), cookieErr.HTTPStatus)
					return
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			//if not valid return error
			if !isCookieValid(cookie) {
				utils.SugaredLogger.Debugln("Cookie is not valid")
				http.Error(w, custom.ErrUnauthorized.Error(), http.StatusUnauthorized)
				return
			}

			//if valid - continue

		} else if r.Method == http.MethodPost {
			//if cookie error or cookie invalid -> create a new one
			if err != nil || !isCookieValid(cookie) {
				err = createNewCookie(w, r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			//if cookie exists and is valid -> continue
		}
		//get cookie again in case if it was re/created
		cookie, err = getShortenerCookie(r)
		//get sessionID and insert into context
		//ignore err, it was checked in isCookieValid() or cookie was recently created
		sessionID, _ := utils.GetSessionID(cookie.Value)

		//set value into context
		ctxWithValue := model.NewContextWithSessionID(r.Context(), sessionID)

		//replace ctx
		newR := r.WithContext(ctxWithValue)

		h.ServeHTTP(w, newR)
	})
}

func getShortenerCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, custom.NewCookieError(err, http.StatusUnauthorized)
		}

		return nil, custom.NewCookieError(err, http.StatusBadRequest)
	}

	return cookie, nil
}

func createNewCookie(w http.ResponseWriter, r *http.Request) error {
	jwtStr, err := utils.NewJWTString()
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: jwtStr,
	}

	r.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: jwtStr,
	})

	http.SetCookie(w, cookie)

	return nil
}

func isCookieValid(cookie *http.Cookie) bool {
	if cookie == nil {
		utils.SugaredLogger.Debugln("cookie.Value == \"\"")
		return false
	}
	//check if jwt is valid
	if cookie.Value == "" {
		utils.SugaredLogger.Debugln("cookie.Value == \"\"")
		return false
	}

	if utils.IsExpired(cookie.Value) {
		utils.SugaredLogger.Debugln("IsExpired")
		return false
	}

	sessionID, err := utils.GetSessionID(cookie.Value)
	if err != nil {
		utils.SugaredLogger.Debugln("err", err)
		return false
	}

	if sessionID == "" {
		utils.SugaredLogger.Debugln("sessionID == \"\". sessionID:", sessionID)
		return false
	}

	return true
}
