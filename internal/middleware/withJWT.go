package middleware

import (
	"errors"
	"github.com/hddskull/urlShorty/internal/model"
	"github.com/hddskull/urlShorty/internal/utils"
	"github.com/hddskull/urlShorty/tools/custom"
	"net/http"
	"strings"
)

const cookieName = "ShortenerCookie"

func WithJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var err error
		var cookieErr *custom.CookieError

		arr := strings.Split(r.URL.Path, "/")
		utils.SugaredLogger.Debugf("arr: %v", arr)

		if strings.ToLower(arr[len(arr)-1]) == "urls" {
			utils.SugaredLogger.Debugf("userPath")
			err = userPath(r)
		} else {
			utils.SugaredLogger.Debugf("otherPaths")
			err = otherPaths(w, r)
		}

		if err != nil {
			if errors.As(err, &cookieErr) {
				custom.JSONError(w, cookieErr, cookieErr.HTTPStatus)
				return
			} else {
				custom.JSONError(w, err, http.StatusInternalServerError)
				return
			}
		}

		//put sessionID into context
		cookie, err := getShortenerCookie(r)
		if err != nil {
			custom.JSONError(w, err, http.StatusInternalServerError)
			return
		}

		sessionID, err := utils.GetSessionID(cookie.Value)
		if err != nil {
			custom.JSONError(w, err, http.StatusInternalServerError)
			return
		}

		ctxWithValue := model.NewContextWithSessionID(r.Context(), sessionID)
		newR := r.WithContext(ctxWithValue)

		h.ServeHTTP(w, newR)
		//___________________________________________________________________________________________________________

		////Check if cookie exists -> if not create
		//cookie, err := getShortenerCookie(r)
		//if err != nil {
		//	err = createNewCookie(w, r)
		//	if err != nil {
		//		custom.JSONError(w, err, http.StatusInternalServerError)
		//		return
		//	}
		//}
		//
		////[for "/user/urls" endpoint] Check if endpoint is  and that cookie is valid
		//arr := strings.Split(r.URL.Path, "/")
		//if strings.ToLower(arr[len(arr)-1]) == "urls" && !isCookieValid(cookie) {
		//	custom.JSONError(w, custom.ErrUnauthorized, http.StatusUnauthorized)
		//	return
		//}
		//
		////[for all other endpoints] Check if cookie is valid -> if not create new and return 401
		//if !isCookieValid(cookie) {
		//	err = createNewCookie(w, r)
		//	if err != nil {
		//		custom.JSONError(w, err, http.StatusInternalServerError)
		//		return
		//	}
		//}
		//
		////put cookie into context
		////ignore err, it was checked in isCookieValid() or cookie was recently created
		//sessionID, _ := utils.GetSessionID(cookie.Value)
		//ctxWithValue := model.NewContextWithSessionID(r.Context(), sessionID)
		//newR := r.WithContext(ctxWithValue)
		//
		//h.ServeHTTP(w, newR)

		//OLD
		//___________________________________________________________________________________________________________
		//get cookie

		//cookie, err := getShortenerCookie(r)
		//if r.Method == http.MethodGet {
		//	utils.SugaredLogger.Debugln("getShortenerCookie(...) cookie:", cookie, " | err:", err)
		//	if err != nil {
		//		var cookieErr *custom.CookieError
		//		if errors.As(err, &cookieErr) {
		//			http.Error(w, cookieErr.Error(), cookieErr.HTTPStatus)
		//			return
		//		} else {
		//			http.Error(w, err.Error(), http.StatusInternalServerError)
		//			return
		//		}
		//	}
		//
		//	//if not valid return error
		//	if !isCookieValid(cookie) {
		//		utils.SugaredLogger.Debugln("Cookie is not valid")
		//		http.Error(w, custom.ErrUnauthorized.Error(), http.StatusUnauthorized)
		//		return
		//	}
		//
		//	//if valid - continue
		//
		//} else if r.Method == http.MethodPost {
		//	//if cookie error or cookie invalid -> create a new one
		//	if err != nil || !isCookieValid(cookie) {
		//		err = createNewCookie(w, r)
		//		if err != nil {
		//			http.Error(w, err.Error(), http.StatusInternalServerError)
		//			return
		//		}
		//	}
		//
		//	//if cookie exists and is valid -> continue
		//}
		////get cookie again in case if it was re/created
		//cookie, err = getShortenerCookie(r)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//
		////get sessionID and insert into context
		////ignore err, it was checked in isCookieValid() or cookie was recently created
		//sessionID, _ := utils.GetSessionID(cookie.Value)
		//
		////set value into context
		//ctxWithValue := model.NewContextWithSessionID(r.Context(), sessionID)
		//
		////replace ctx
		//newR := r.WithContext(ctxWithValue)
		//
		//h.ServeHTTP(w, newR)
	})
}

func userPath(r *http.Request) error {
	cookie, err := getShortenerCookie(r)
	if err != nil {
		return custom.NewCookieError(err, http.StatusUnauthorized)
	}

	if !isCookieValid(cookie) {
		return custom.NewCookieError(custom.ErrUnauthorized, http.StatusUnauthorized)
	}

	return nil
}

func otherPaths(w http.ResponseWriter, r *http.Request) error {
	cookie, err := getShortenerCookie(r)
	if err != nil || !isCookieValid(cookie) {
		err = createNewCookie(w, r)
		if err != nil {
			return custom.NewCookieError(err, http.StatusInternalServerError)
		}
	}

	return nil
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
