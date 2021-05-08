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
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "", http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	data := []byte(req.Email)
	emailHash := sha256.Sum256(data)

	password := []byte(req.Password)

	user := User{}
	if res := db.First(&user, "email_hash = ?", fmt.Sprintf("%x", emailHash[:])); res.Error != nil || user.ID == 0 {
		log.Println(res.Error)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), password)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	// Auth success
	loginAndRedirect(user, w, r)
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	gs, err := store.Get(r, cookieName)
	if err != nil {
		log.Println(err.Error())
	} else {
		if val, ok := gs.Values["session"]; ok {
			if ses, ok := val.(Session); ok {
				db.Where("auth_session_id = ?", ses.AuthSessionID).Delete(&AuthSession{})
				gs.Values["session"] = Session{}
				gs.Options.MaxAge = -1
				gs.Options.Path = "/"
				gs.Options.HttpOnly = true
				gs.Options.SameSite = http.SameSiteStrictMode
				gs.Options.Secure = secureCookies
				err = gs.Save(r, w)
				if err != nil {
					log.Println(err.Error())
				}
				// redirect and prevent further writes
				http.Redirect(w, r, "/", 302)
				return
			}
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Price    string `json:"priceId"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "", http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	// Validate password
	if len(req.Password) < 8 {
		writeJSONError(w, "password is too short", http.StatusInternalServerError)
		return
	}

	// Validate email
	if err := emailx.Validate(req.Email); err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	data := []byte(req.Email)
	hash := sha256.Sum256(data)

	// Make sure there isn't a user with that email
	user := User{}
	db.Where("email_hash = ?", fmt.Sprintf("%x", hash[:])).First(&user)
	if user.ID != 0 {
		writeJSONError(w, "", http.StatusInternalServerError)
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
		writeJSONError(w, "", http.StatusBadRequest)
		return
	}

	gs, err := store.Get(r, cookieName)
	if err != nil {
		log.Println(cookieName)
		log.Println(err.Error())
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	password := []byte(req.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
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
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	gs.Values["session"] = Session{
		TempSessionID:   tmp.SessionString,
		StripeSessionID: s.ID,
	}
	gs.Options.Path = "/"
	gs.Options.HttpOnly = true
	gs.Options.SameSite = http.SameSiteStrictMode
	gs.Options.Secure = secureCookies

	err = gs.Save(r, w)
	if err != nil {
		log.Println(gs.Values["session"])
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	writeJSON(w, struct {
		SessionID string `json:"sessionId"`
	}{
		SessionID: s.ID,
	}, 200)
}

func forgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "", http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	data := []byte(req.Email)
	emailHash := sha256.Sum256(data)

	user := User{}
	if res := db.First(&user, "email_hash = ?", fmt.Sprintf("%x", emailHash[:])); res.Error != nil || user.ID == 0 {
		log.Println(res.Error)
		// We don't want to disclose if the email exists or not
		writeJSON(w, nil, 200)
		return
	}

	passwordReset := PasswordReset{}
	if res := db.First(&passwordReset, "user_id = ? AND expires_at > ?", user.ID, time.Now()); res.Error != nil || user.ID == 0 {
		log.Println(res.Error)
		// We don't want to disclose if the email exists or not
		writeJSON(w, nil, 200)
		return
	}

	passwordReset = PasswordReset{
		UserID: user.ID,
		Code:   uuid.NewV4().String(),
		// Password reset links are valid for 1 hour
		ExpiresAt: time.Now().Add(time.Hour),
	}

	if res := db.Create(&passwordReset); res.Error != nil || passwordReset.ID == 0 {
		log.Println(res.Error)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	go sendEmail(req.Email, "Triangulate.xyz Password Reset", fmt.Sprintf("Please use the following link to reset your password:\n\nhttps://%s/reset-password/%s/", domain, passwordReset.Code))
	writeJSON(w, nil, 200)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "", http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	data := []byte(req.Email)
	emailHash := sha256.Sum256(data)

	user := User{}
	if res := db.First(&user, "email_hash = ?", fmt.Sprintf("%x", emailHash[:])); res.Error != nil || user.ID == 0 {
		log.Println(res.Error)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	passwordReset := PasswordReset{}
	if res := db.First(&passwordReset, "code = ? AND expires_at > ?", req.Code, time.Now()); res.Error != nil || passwordReset.ID == 0 {
		log.Println(res.Error)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	if passwordReset.UserID != user.ID {
		log.Println("password reset user does not match email")
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = string(hashedPassword)

	if res := db.Save(&user); res.Error != nil {
		log.Println(res.Error)
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	writeJSON(w, nil, 200)
}

func loginAndRedirect(user User, w http.ResponseWriter, r *http.Request) {
	gs, err := store.Get(r, cookieName)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	} else {
		sesID := uuid.NewV4().String()
		// store session id
		authSession := AuthSession{
			UserID:        user.ID,
			AuthSessionID: sesID,
		}
		db.Create(&authSession)
		if authSession.ID > 0 {
			// set cookie
			gs.Values["session"] = Session{
				AuthSessionID: sesID,
			}
			gs.Options.Path = "/"
			gs.Options.HttpOnly = true
			gs.Options.SameSite = http.SameSiteStrictMode
			gs.Options.Secure = secureCookies
			err = gs.Save(r, w)
			if err == nil {
				// redirect and prevent further writes
				http.Redirect(w, r, "/", 302)
				return
			}
		}
	}
}
