package middleware

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MatomoProxy creates a middleware that proxies Matomo requests to bypass adblockers
func MatomoProxy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the request is for Matomo resources
		path := c.Path()

		// Debug: log all requests to see if middleware is triggered

		// Proxy matomo.php and matomo.js requests
		if strings.HasSuffix(path, "/matomo.php") || strings.HasSuffix(path, "/matomo.js") {
			log.Printf("Matomo proxy middleware triggered for path: %s", path)

			// Get Matomo URL from environment variable
			matomoURL := os.Getenv("MATOMO_URL")
			if matomoURL == "" {
				matomoURL = "https://matomo.romaingallez.fr" // fallback default
			}

			// Ensure the URL has a protocol scheme
			if !strings.HasPrefix(matomoURL, "http://") && !strings.HasPrefix(matomoURL, "https://") {
				matomoURL = "https://" + matomoURL
			}

			// Construct the target URL with query parameters
			targetURL := matomoURL + path

			// Add query parameters if they exist
			queryString := c.Context().QueryArgs().QueryString()
			if len(queryString) > 0 {
				targetURL += "?" + string(queryString)
			}

			// Debug logging
			log.Printf("Matomo proxy: %s %s -> %s", c.Method(), c.Path(), targetURL)

			// Create HTTP client
			client := &http.Client{}

			// Create request with proper body handling
			var req *http.Request
			var err error

			if c.Method() == "POST" {
				// For POST requests, use the body
				body := c.Body()
				log.Printf("Matomo proxy POST body: %s", string(body))
				req, err = http.NewRequest(c.Method(), targetURL, strings.NewReader(string(body)))
			} else {
				// For GET requests, no body needed
				req, err = http.NewRequest(c.Method(), targetURL, nil)
			}

			if err != nil {
				log.Printf("Matomo proxy error creating request: %v", err)
				return c.Status(500).SendString("Failed to create proxy request")
			}

			// Copy headers from original request
			for key, values := range c.GetReqHeaders() {
				// Skip some headers that shouldn't be forwarded
				if key != "Host" && key != "Connection" {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}
			}

			// Set User-Agent if not present
			if req.Header.Get("User-Agent") == "" {
				req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MatomoProxy/1.0)")
			}

			// Set proper content type for POST requests
			if c.Method() == "POST" {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			// Set additional headers that Matomo might expect
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("Pragma", "no-cache")

			// Set referer to match the original request
			if referer := c.Get("Referer"); referer != "" {
				req.Header.Set("Referer", referer)
			}

			// Debug headers
			log.Printf("Matomo proxy request headers: %v", req.Header)

			// Make the request
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Matomo proxy error making request: %v", err)
				return c.Status(502).SendString("Failed to proxy request to Matomo")
			}
			defer resp.Body.Close()

			// Debug response
			log.Printf("Matomo proxy response status: %d", resp.StatusCode)

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Matomo proxy error reading response: %v", err)
				return c.Status(502).SendString("Failed to read proxy response")
			}

			// Set response headers
			for key, values := range resp.Header {
				for _, value := range values {
					c.Set(key, value)
				}
			}

			// Set content type for JavaScript files
			if strings.HasSuffix(path, "/matomo.js") {
				c.Set("Content-Type", "application/javascript")
			}

			// Set cache control headers to prevent caching
			c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Set("Pragma", "no-cache")
			c.Set("Expires", "0")

			// Return the proxied response
			return c.Status(resp.StatusCode).Send(body)
		}

		// If not a Matomo request, continue to next middleware
		return c.Next()
	}
}
