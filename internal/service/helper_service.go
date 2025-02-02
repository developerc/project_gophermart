package service

import (
	"net/http"
)

type User struct {
	Name string
}

func (s *Service) SetUserCookie(usr string) (*http.Cookie, error) {
	//var usr string
	var cookie *http.Cookie
	//var err error
	u := &User{
		Name: usr,
	}

	//if len(cookieValue) == 0 {
	u.Name = usr
	if encoded, err := s.secure.Encode("user", u); err == nil {
		cookie = &http.Cookie{
			Name:  "user",
			Value: encoded,
		}
		return cookie, nil
	} else {
		return nil, err
	}
	//}
	//return nil,  nil
}
