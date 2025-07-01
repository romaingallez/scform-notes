package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MatomoInjector creates a middleware that injects Matomo tracking script into HTML responses
func MatomoInjector() fiber.Handler {
	log.Println("MatomoInjector middleware triggered")
	return func(c *fiber.Ctx) error {
		// Continue to next middleware/handler
		err := c.Next()
		if err != nil {
			return err
		}

		// Check if response is HTML
		contentType := c.Response().Header.ContentType()
		log.Println("Content-Type: ", string(contentType))

		// If content type is not set or not HTML, check if the response body looks like HTML
		if len(contentType) == 0 || !strings.Contains(string(contentType), "text/html") {
			// Get the response body to check if it looks like HTML
			body := c.Response().Body()
			if len(body) == 0 {
				return nil
			}

			bodyStr := string(body)
			// Check if the response starts with HTML tags
			if !strings.HasPrefix(strings.TrimSpace(bodyStr), "<!DOCTYPE") &&
				!strings.HasPrefix(strings.TrimSpace(bodyStr), "<html") {
				log.Println("Response is not HTML, skipping Matomo injection")
				return nil
			}
		}

		// Get Matomo configuration from environment
		matomoURL := os.Getenv("MATOMO_URL")
		matomoSiteID := os.Getenv("MATOMO_SITE_ID")

		// If Matomo is not configured, don't inject anything
		if matomoURL == "" || matomoSiteID == "" {
			return nil
		}

		// Get the response body (we already checked it's not empty above)
		body := c.Response().Body()
		bodyStr := string(body)

		// Check if Matomo script is already present
		if strings.Contains(bodyStr, "_paq") {
			return nil
		}

		// Create Matomo script
		matomoScript := `
<!-- Matomo -->
<script>
  var _paq = window._paq = window._paq || [];
  /* tracker methods like "setCustomDimension" should be called before "trackPageView" */
  _paq.push(['trackPageView']);
  _paq.push(['enableLinkTracking']);
  (function() {
    var u=window.location.origin+"/";
    _paq.push(['setTrackerUrl', u+'matomo.php']);
    _paq.push(['setSiteId', '` + matomoSiteID + `']);
    var d=document, g=d.createElement('script'), s=d.getElementsByTagName('script')[0];
    g.async=true; g.src=u+'matomo.js'; s.parentNode.insertBefore(g,s);
  })();
</script>
<!-- End Matomo Code -->`

		// Find the closing </head> tag and insert the script before it
		headEndIndex := strings.Index(bodyStr, "</head>")
		if headEndIndex == -1 {
			// If no </head> tag found, try to insert before </body>
			bodyEndIndex := strings.Index(bodyStr, "</body>")
			if bodyEndIndex == -1 {
				// If no </body> tag either, append at the end
				newBody := bodyStr + matomoScript
				c.Response().SetBody([]byte(newBody))
				return nil
			}
			// Insert before </body>
			newBody := bodyStr[:bodyEndIndex] + matomoScript + bodyStr[bodyEndIndex:]
			c.Response().SetBody([]byte(newBody))
			return nil
		}

		// Insert before </head>
		newBody := bodyStr[:headEndIndex] + matomoScript + bodyStr[headEndIndex:]
		c.Response().SetBody([]byte(newBody))

		return nil
	}
}
