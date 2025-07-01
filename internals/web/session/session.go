package session

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Manager holds the session store and provides session utilities
type Manager struct {
	Store *session.Store
}

// NewManager creates a new session manager
func NewManager() *Manager {
	store := session.New(session.Config{
		KeyLookup:      "cookie:session_id",
		CookieDomain:   "",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		Expiration:     0, // Session cookie (expires when browser closes)
	})

	return &Manager{
		Store: store,
	}
}

// GenerateSessionID creates a unique session ID based on browser characteristics
func (m *Manager) GenerateSessionID(c *fiber.Ctx) string {
	// Get unique browser characteristics
	userAgent := c.Get("User-Agent")
	acceptLanguage := c.Get("Accept-Language")
	acceptEncoding := c.Get("Accept-Encoding")
	ip := c.IP()
	xForwardedFor := c.Get("X-Forwarded-For")

	// Create a unique fingerprint
	fingerprint := fmt.Sprintf("%s|%s|%s|%s|%s", userAgent, acceptLanguage, acceptEncoding, ip, xForwardedFor)

	// Generate a hash-based session ID
	hash := sha256.Sum256([]byte(fingerprint))
	sessionID := fmt.Sprintf("%x", hash)[:32] // Use first 32 characters

	return sessionID
}

// GetSessionID retrieves the session ID from the context
func (m *Manager) GetSessionID(c *fiber.Ctx) string {
	sess, err := m.Store.Get(c)
	if err != nil {
		log.Printf("GetSessionID: Failed to get session for path %s: %v", c.Path(), err)
		return ""
	}

	// First try to get our custom fingerprint ID
	sessionID := sess.Get("fingerprint_id")
	if sessionID != nil {
		if sessionIDStr, ok := sessionID.(string); ok {
			log.Printf("GetSessionID: Using fingerprint ID %s for path %s", sessionIDStr, c.Path())
			return sessionIDStr
		}
	}

	// Fallback to Fiber's session ID
	fiberSessionID := sess.ID()
	if fiberSessionID != "" {
		log.Printf("GetSessionID: Using Fiber session ID %s for path %s", fiberSessionID, c.Path())
		return fiberSessionID
	}

	log.Printf("GetSessionID: No session ID found for path %s", c.Path())
	return ""
}

// SetupSessionMiddleware configures session middleware for the app
func (m *Manager) SetupSessionMiddleware(app *fiber.App) {
	// Session middleware first (let Fiber handle session creation)
	app.Use(func(c *fiber.Ctx) error {
		// Get or create session
		sess, err := m.Store.Get(c)
		if err != nil {
			log.Printf("Failed to get session store for path: %s, error: %v", c.Path(), err)
			return c.Status(500).SendString("Failed to get session")
		}

		// Store browser fingerprint as additional session data for user identification
		fingerprint := m.GenerateSessionID(c)
		existingFingerprint := sess.Get("fingerprint_id")

		if existingFingerprint == nil {
			sess.Set("fingerprint_id", fingerprint)
			if err := sess.Save(); err != nil {
				log.Printf("Failed to save fingerprint for path: %s, error: %v", c.Path(), err)
			} else {
				log.Printf("Browser fingerprint saved: %s for path: %s", fingerprint, c.Path())
			}
		} else {
			log.Printf("Using existing fingerprint: %s for path: %s", existingFingerprint.(string), c.Path())
		}

		// Store session ID in context for easy access
		sessionID := m.GetSessionID(c)
		c.Locals("session_id", sessionID)
		log.Printf("Session ID set in context: %s for path: %s", sessionID, c.Path())

		return c.Next()
	})
}
