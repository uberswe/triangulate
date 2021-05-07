package triangulate

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"log"
	"net/http"
)

func StripeSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := ""

		// We check for a get parameter with the session id
		if r.Method == "GET" {
			sessionID = r.URL.Query().Get(sessionIDParam)
		}

		if sessionID != "" {
			// Fetch session object from stripe
			s, err := session.Get(
				sessionID,
				nil,
			)
			// If no error then we seem to have a customer
			if err == nil {
				// Check if the customer paid
				if s.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid && s.Customer != nil {
					// Make sure a temporary session exists for the session id
					mutex.Lock()
					var t TempSession
					if res := db.First(&t, "stripe_session_id = ?", sessionID); res.Error == nil {
						user := User{
							EmailHash:        t.Email,
							PasswordHash:     t.Password,
							StripeCustomerID: s.Customer.ID,
						}
						db.Create(&user)
						if user.ID > 0 {
							db.Delete(&t)
							mutex.Unlock()

							// login && redirect
							gs, err := store.Get(r, cookieName)
							if err != nil {
								log.Println(err.Error())
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
									err = gs.Save(r, w)
									if err == nil {
										// redirect and prevent further writes
										http.Redirect(w, r, "/members", 302)
										return
									}
								}
							}

						} else {
							mutex.Unlock()
						}
					} else {
						mutex.Unlock()
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
