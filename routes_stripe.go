package triangulate

import (
	"github.com/stripe/stripe-go/v72"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/webhook"
	"io/ioutil"
	"log"
	"net/http"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "", http.StatusMethodNotAllowed)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, "", http.StatusBadRequest)
		log.Printf("ioutil.ReadAll: %v", err)
		return
	}

	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		writeJSONError(w, "", http.StatusBadRequest)
		log.Printf("webhook.ConstructEvent: %v", err)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		log.Println("Payment is successful and the subscription is created.")
		log.Println(event.ID)
		log.Println(event.Account)
		log.Println(event.Data.Raw)
		// Payment is successful and the subscription is created.
		// You should provision the subscription and save the customer ID to your database.
	case "invoice.paid":
		log.Println("Payment for a subscription is made.")
		log.Println(event.ID)
		log.Println(event.Account)
		log.Println(event.Data.Raw)
		// Continue to provision the subscription as payments continue to be made.
		// Store the status in your database and check when a user accesses your service.
		// This approach helps you avoid hitting rate limits.
	case "invoice.payment_failed":
		log.Println("The payment failed or the customer does not have a valid payment method. The subscription is now past due.")
		log.Println(event.ID)
		log.Println(event.Account)
		log.Println(event.Data.Raw)
		// The payment failed or the customer does not have a valid payment method.
		// The subscription becomes past_due. Notify your customer and send them to the
		// customer portal to update their payment information.
	default:
		// unhandled event type
		log.Println("Unhandled event type:")
		log.Println(event.Type)
	}
}

func portal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeJSONError(w, "", http.StatusMethodNotAllowed)
		return
	}

	uid := r.Context().Value(ContextUserKey)
	if uid == nil || uid.(uint) == 0 {
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}
	user := User{}
	db.Where("id = ?", uid).First(&user)

	if user.ID == 0 {
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	if user.StripeCustomerID == "" {
		log.Println("customer id missing for user")
		writeJSONError(w, "", http.StatusInternalServerError)
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(user.StripeCustomerID),
		ReturnURL: stripe.String(returnURL),
	}
	ps, _ := portalsession.New(params)

	writeJSON(w, struct {
		URL string `json:"url"`
	}{
		URL: ps.URL,
	}, 200)
}
