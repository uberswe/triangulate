package triangulate

import (
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
	"sync"
)

var (
	store                   *sessions.CookieStore
	sourceDir               = "resources/source"
	outDir                  = "resources/out"
	queue                   []string
	images                  map[string]Image
	mutex                   = &sync.Mutex{}
	jobChan                 = make(chan Image, 999999)
	currentJob              Image
	db                      *gorm.DB
	priceID                 string
	stripePrivateKey        string
	stripePublicKey         string
	successUrl              string
	cancelUrl               string
	returnURL               string
	webhookSecret           string
	cookieName              string
	sessionIDParam          string
	sessionAuthKey          string
	sessionEncryptionKey    string
	addr                    string
	randomUnsplashImages    []UnsplashRandomImageResponse
	strictTransportSecurity = true
	secureCookies           = true
	unsplashAccessKey       string
)
