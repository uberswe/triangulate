package triangulate

import (
	"context"
	"log"
	"net/http"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gs, err := store.Get(r, cookieName)
		if err != nil {
			log.Println(err.Error())
		} else {
			if val, ok := gs.Values["session"]; ok {
				if ses, ok := val.(Session); ok {
					if ses.AuthSessionID != "" {
						a := AuthSession{}
						if res := db.First(&a, "auth_session_id = ?", ses.AuthSessionID); res.Error == nil {
							if a.ID > 0 {
								ctx := context.WithValue(r.Context(), ContextUserKey, a.UserID)
								r.WithContext(ctx)
							}
						}
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
