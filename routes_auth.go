package triangulate

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/uberswe/emailx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {

}

func logout(w http.ResponseWriter, r *http.Request) {

}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Price    string `json:"priceId"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	// Validate email
	if err := emailx.Validate(req.Email); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// See https://stripe.com/docs/api/checkout/sessions/create
	// for additional parameters to pass.
	// {CHECKOUT_SESSION_ID} is a string literal; do not change it!
	// the actual Session ID is returned in the query parameter when your customer
	// is redirected to the success page.
	params := &stripe.CheckoutSessionParams{
		SuccessURL: &successUrl,
		CancelURL:  &cancelUrl,
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price: stripe.String(req.Price),
				// For metered billing, do not pass quantity
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := session.New(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, struct {
			ErrorData string `json:"error"`
		}{
			ErrorData: "test",
		})
		return
	}

	gs, err := store.Get(r, cookieName)
	if err != nil {
		log.Println(cookieName)
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := []byte(req.Email)
	hash := sha256.Sum256(data)

	password := []byte(req.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sesID := uuid.NewV4()

	tmp := TempSession{
		SessionString:   sesID.String(),
		Email:           fmt.Sprintf("%x", hash[:]),
		Password:        string(hashedPassword),
		StripeSessionID: s.ID,
	}

	res := db.Create(&tmp)
	if res.Error != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gs.Values["session"] = Session{
		TempSessionID:   tmp.SessionString,
		StripeSessionID: s.ID,
	}

	err = gs.Save(r, w)
	if err != nil {
		log.Println(gs.Values["session"])
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, struct {
		SessionID string `json:"sessionId"`
	}{
		SessionID: s.ID,
	})
}

func forgotPassword(w http.ResponseWriter, r *http.Request) {

}

func resetPassword(w http.ResponseWriter, r *http.Request) {

}
