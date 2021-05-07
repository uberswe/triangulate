package triangulate

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"log"
)

func initSessions() {
	authKeyOne, err := base64Decode([]byte(sessionAuthKey))
	if err != nil {
		log.Fatal(err)
	}

	encryptionKeyOne, err := base64Decode([]byte(sessionEncryptionKey))
	if err != nil {
		log.Fatal(err)
	}

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 60 * 24 * 30, // 30 days
		HttpOnly: true,
	}

	gob.Register(Session{})
}
