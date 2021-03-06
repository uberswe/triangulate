package triangulate

import (
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
	"log"
	"os"
	"strconv"
)

func initSettings() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error reading .env file")
	}

	addr = os.Getenv("ADDR")
	stripePrivateKey = os.Getenv("STRIPE_PRIVATE_KEY")
	stripePublicKey = os.Getenv("STRIPE_PUBLIC_KEY")
	unsplashAccessKey = os.Getenv("UNSPLASH_ACCESS_KEY")
	successUrl = os.Getenv("SUCCESS_URL")
	cancelUrl = os.Getenv("CANCEL_URL")
	returnURL = os.Getenv("RETURN_URL")
	webhookSecret = os.Getenv("WEBHOOK_SECRET")
	priceID = os.Getenv("PRICE_ID")
	cookieName = os.Getenv("COOKIE_NAME")
	sessionIDParam = os.Getenv("SESSION_ID_PARAM")
	sessionAuthKey = os.Getenv("SESSION_AUTH_KEY")
	sessionEncryptionKey = os.Getenv("SESSION_ENCRYPTION_KEY")
	smtpUsername = os.Getenv("SMTP_USERNAME")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	fromEmail = os.Getenv("FROM_EMAIL")
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	domain = os.Getenv("DOMAIN")

	workerCountString := os.Getenv("WORKER_COUNT")
	workerCount, err = strconv.Atoi(workerCountString)
	if err != nil {
		log.Fatal(err)
	}

	premiumWorkerCountString := os.Getenv("PREMIUM_WORKER_COUNT")
	premiumWorkerCount, err = strconv.Atoi(premiumWorkerCountString)
	if err != nil {
		log.Fatal(err)
	}
	sts := os.Getenv("STRICT_TRANSPORT_SECURITY")
	strictTransportSecurity = sts != "false"
	sc := os.Getenv("SECURE_COOKIES")
	secureCookies = sc != "false"

	if sessionAuthKey == "" || sessionEncryptionKey == "" {
		sessionAuthKey = string(base64Encode(securecookie.GenerateRandomKey(64)))
		sessionEncryptionKey = string(base64Encode(securecookie.GenerateRandomKey(32)))
		_ = godotenv.Write(mergeMaps(envMap, map[string]string{
			"SESSION_AUTH_KEY":       sessionAuthKey,
			"SESSION_ENCRYPTION_KEY": sessionEncryptionKey,
		}), ".env")
	}

	stripe.Key = stripePrivateKey
}
