package triangulate

import (
	"encoding/json"

	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/build/index.html")
}

func settings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	loggedIn := false
	user := r.Context().Value(ContextUserKey)
	if user != nil && user.(uint) > 0 {
		loggedIn = true
	}

	s := Settings{
		PriceId:   priceID,
		StripeKey: stripePublicKey,
		LoggedIn:  loggedIn,
	}

	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
