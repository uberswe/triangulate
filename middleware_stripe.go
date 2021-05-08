package triangulate

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
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
							loginAndRedirect(user, w, r)

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
