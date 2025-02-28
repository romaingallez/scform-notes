package scform

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/launcher"
)

var debugEnabled bool

// Grade represents a single grade entry
type Grade struct {
	Value       float64   // The numerical grade value
	OutOf       float64   // The maximum possible grade (usually 20)
	Coefficient float64   // Grade coefficient
	Title       string    // Title/name of the grade
	Date        time.Time // Date of the grade
	Type        string    // Type of grade (exam, homework, etc.)
	Remarks     string    // Any remarks about the grade
	Observation string    // Any observations about the grade
}

// Course represents a course/subject with its grades
type Course struct {
	Name    string  // Course name
	Grades  []Grade // List of grades for this course
	Average float64 // Weighted average for the course
}

// CalculateAverage calculates the weighted average for the course
func (c *Course) CalculateAverage() {
	var totalWeightedGrade float64
	var totalCoefficient float64

	for _, grade := range c.Grades {
		if grade.Value > 0 && grade.Coefficient > 0 { // Only consider valid grades
			totalWeightedGrade += grade.Value * grade.Coefficient
			totalCoefficient += grade.Coefficient
		}
	}

	if totalCoefficient > 0 {
		c.Average = totalWeightedGrade / totalCoefficient
	}
}

type Student struct {
	Name         string   // Student name
	Grades       []Course // List of grades for this student
	TotalAverage float64  // Overall weighted average
}

// CalculateTotalAverage calculates the overall weighted average for all courses
func (s *Student) CalculateTotalAverage() {
	var totalWeightedGrade float64
	var totalCoefficient float64

	for i := range s.Grades {
		s.Grades[i].CalculateAverage() // Calculate average for each course

		// Calculate total weighted average
		for _, grade := range s.Grades[i].Grades {
			if grade.Value > 0 && grade.Coefficient > 0 {
				totalWeightedGrade += grade.Value * grade.Coefficient
				totalCoefficient += grade.Coefficient
			}
		}
	}

	if totalCoefficient > 0 {
		s.TotalAverage = totalWeightedGrade / totalCoefficient
	}
}

func init() {
	// setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set debug mode from environment variable
	debugEnabled = os.Getenv("SCFORM_DEBUG") == "true"
}

// debugLog logs a message only if debug mode is enabled
func debugLog(format string, v ...interface{}) {
	if debugEnabled {
		log.Printf(format, v...)
	}
}

// ProgressUpdate represents a progress update message
type ProgressUpdate struct {
	Status   string  `json:"status"`
	Message  string  `json:"message"`
	Progress float64 `json:"progress"`
}

func GetStudentGrades(scformURL, username, password string, progressChan chan<- ProgressUpdate) (*Student, error) {
	// Send initial progress update
	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Status:   "connecting",
			Message:  "Connecting to browser...",
			Progress: 0.1,
		}
	}

	remoteURL := os.Getenv("SCFORM_REMOTE_URL")
	debugLog("remoteURL: %s", remoteURL)

	var browser *rod.Browser
	var err error

	// Set default timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if remoteURL != "" {
		// Add error handling for remote connection
		if strings.HasPrefix(remoteURL, "ws://") || strings.HasPrefix(remoteURL, "wss://") {
			debugLog("Connecting to remote browser via WebSocket at: %s", remoteURL)
			ws := NewWebSocket(remoteURL)
			client := cdp.New().Start(ws)
			browser = rod.New().Client(client).Context(ctx)
			err = browser.Connect()
			if err != nil {
				debugLog("Failed to connect browser via WebSocket: %v", err)
			} else {
				debugLog("Successfully connected browser via WebSocket")
			}
		} else {
			debugLog("Connecting to remote browser via direct URL: %s", remoteURL)
			browser = rod.New().ControlURL(remoteURL).Context(ctx).MustConnect()
		}
		if err != nil {
			return nil, fmt.Errorf("failed to connect to remote browser: %v", err)
		}
		browser = browser.NoDefaultDevice()
	} else {
		// use already installed chrome browser
		chromePath := os.Getenv("CHROME_PATH")

		// Use headless by default, only use headed mode for debugging
		useHeadless := os.Getenv("SCFORM_HEADLESS") != "false"

		var l *launcher.Launcher
		if chromePath == "" {
			path, _ := launcher.LookPath()
			l = launcher.New().Bin(path).Headless(useHeadless)
		} else {
			l = launcher.New().Bin(chromePath).Headless(useHeadless)
		}

		// Set browser flags for better performance
		l = l.Set("disable-gpu", "true").
			Set("disable-dev-shm-usage", "true").
			Set("disable-web-security", "true").
			Set("disable-features", "IsolateOrigins,site-per-process").
			Set("disable-site-isolation-trials", "true").
			Set("disable-blink-features", "AutomationControlled").
			Set("blink-settings", "imagesEnabled=false")

		// Launch and connect to the browser
		url := l.MustLaunch()
		browser = rod.New().ControlURL(url).Context(ctx).MustConnect().NoDefaultDevice()
	}

	// Set a default shorter timeout for all browser operations
	browser = browser.Timeout(30 * time.Second)

	// Send progress update
	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Status:   "navigating",
			Message:  "Navigating to login page...",
			Progress: 0.2,
		}
	}

	// Add error handling for browser operations
	defer func() {
		if browser != nil {
			browser.Close()
		}
	}()

	// Set a shorter timeout for page navigation (15 seconds)
	page := browser.Timeout(15 * time.Second).MustPage(scformURL)

	// Use more targeted wait with a timeout instead of waiting for DOM stable
	// Wait up to 5 seconds for the username input field
	usernameInput := page.Timeout(5 * time.Second).MustElement("input[id='MainContent_LoginUser_UserName']")
	usernameInput.MustInput(username)

	// Wait up to 5 seconds for the password input field
	passwordInput := page.Timeout(5 * time.Second).MustElement("input[id='MainContent_LoginUser_Password']")
	passwordInput.MustInput(password)

	// Send progress update
	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Status:   "logging_in",
			Message:  "Logging in...",
			Progress: 0.3,
		}
	}

	// Use a shorter wait timeout of 5 seconds
	page.Timeout(5 * time.Second).MustWaitStable()

	// Click the login button using JavaScript
	page.MustEval(`() => {
		LoginBt();
	}`)

	// Use a shorter wait timeout of 5 seconds
	page.Timeout(5 * time.Second).MustWaitStable()

	// Send progress update
	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Status:   "navigating_grades",
			Message:  "Navigating to grades page...",
			Progress: 0.4,
		}
	}

	// Navigate to grades page using JavaScript
	page.MustEval(`() => {
		GoTo('Eleve/MesNotes.aspx');
	}`)

	// Use a shorter wait timeout of 5 seconds
	page.Timeout(5 * time.Second).MustWaitStable()

	// Use a more efficient selector with timeout of 5 seconds
	page.Timeout(5 * time.Second).MustElement("input[id='MainContent_RadioButtonAffichage_1']").MustClick()

	// Use a shorter wait timeout of 5 seconds
	page.Timeout(5 * time.Second).MustWaitStable()

	// Send progress update
	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Status:   "fetching_grades",
			Message:  "Fetching grades...",
			Progress: 0.5,
		}
	}

	// Get all course tables with a timeout of 10 seconds
	courseTables, err := page.Timeout(10 * time.Second).Elements("table.AfficheInfoEnMieux")
	if err != nil {
		return nil, fmt.Errorf("failed to find course tables: %v", err)
	}

	var courses []Course
	totalTables := len(courseTables)

	for i, table := range courseTables {
		// Send progress update for each course
		if progressChan != nil {
			progress := 0.5 + (float64(i) / float64(totalTables) * 0.4)
			progressChan <- ProgressUpdate{
				Status:   "processing_course",
				Message:  fmt.Sprintf("Processing course %d of %d...", i+1, totalTables),
				Progress: progress,
			}
		}

		// Extract course name with a timeout of 5 seconds
		nameElement, err := table.Timeout(5 * time.Second).Element("span[id*='NomCompletLabel']")
		if err != nil {
			debugLog("Failed to find course name element, skipping table")
			continue
		}

		courseName, err := nameElement.Text()
		if err != nil {
			debugLog("Failed to get course name text, skipping table")
			continue
		}
		courseName = strings.TrimSpace(courseName)

		// Create new course
		course := Course{
			Name:   courseName,
			Grades: []Grade{},
		}

		// Find all grade divs within the table with a timeout of 5 seconds
		gradeDivs, err := table.Timeout(5 * time.Second).Elements("div[id='DivNOTE']")
		if err != nil {
			debugLog("Failed to find grade divs for course %s: %v", courseName, err)
			continue
		}

		for _, gradeDiv := range gradeDivs {
			grade := Grade{}

			// Extract grade value and maximum with a timeout of 2 seconds
			if valueSpan, err := gradeDiv.Timeout(2 * time.Second).Element("span[id*='Label1']"); err == nil {
				if value, err := valueSpan.Text(); err == nil {
					value = strings.TrimSpace(value)
					if value != "" {
						grade.Value = parseFloat(value)
						grade.OutOf = 20 // Default out of 20
					}
				}
			}

			// Extract coefficient with a timeout of 2 seconds
			if coeffSpan, err := gradeDiv.Timeout(2 * time.Second).Element("span[id*='Label3']"); err == nil {
				if coeffText, err := coeffSpan.Text(); err == nil {
					coeffText = strings.TrimSpace(coeffText)
					coeffText = strings.TrimPrefix(coeffText, "coeff. ")
					grade.Coefficient = parseFloat(coeffText)
				}
			}

			// Extract title and date with a timeout of 2 seconds
			if titleSpan, err := gradeDiv.Timeout(2 * time.Second).Element("span[id*='Label7']"); err == nil {
				if titleText, err := titleSpan.Text(); err == nil {
					titleText = strings.TrimSpace(titleText)
					// Use regex to match the pattern: any text followed by "du" and a date
					re := regexp.MustCompile(`(.*?)\s+du\s+(\d{2}/\d{2}/\d{4})`)
					matches := re.FindStringSubmatch(titleText)
					if len(matches) == 3 {
						grade.Title = strings.TrimSpace(matches[1])
						dateStr := strings.TrimSpace(matches[2])
						grade.Date = parseDate(dateStr)
					} else {
						// If no date pattern found, use the entire text as title
						grade.Title = titleText
					}
				}
			}

			// Extract type with a timeout of 2 seconds
			if typeSpan, err := gradeDiv.Timeout(2 * time.Second).Element("span[id*='Label8']"); err == nil {
				if typeText, err := typeSpan.Text(); err == nil {
					grade.Type = strings.TrimSpace(typeText)
				}
			}

			// Extract remarks if present with a timeout of 2 seconds
			if remarksSpan, err := gradeDiv.Timeout(2 * time.Second).Element("span[id*='Label9']"); err == nil {
				if remarks, err := remarksSpan.Text(); err == nil {
					grade.Remarks = strings.TrimSpace(remarks)
					grade.Remarks = strings.TrimPrefix(grade.Remarks, "Remarque : ")
				}
			}

			// Extract observations if present with a timeout of 2 seconds
			if obsSpan, err := gradeDiv.Timeout(2 * time.Second).Element("span[id*='Label10']"); err == nil {
				if obs, err := obsSpan.Text(); err == nil {
					grade.Observation = strings.TrimSpace(obs)
					grade.Observation = strings.TrimPrefix(grade.Observation, "Observation : ")
				}
			}

			// Only append grade if we have at least some basic information
			if grade.Value > 0 || grade.Title != "" {
				course.Grades = append(course.Grades, grade)
			}
		}

		// Only append course if it has a name and at least one grade
		if course.Name != "" && len(course.Grades) > 0 {
			courses = append(courses, course)
		}
	}

	// Send progress update
	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Status:   "calculating",
			Message:  "Calculating averages...",
			Progress: 0.9,
		}
	}

	// Create a student and calculate averages
	student := &Student{
		Name:   username,
		Grades: courses,
	}
	student.CalculateTotalAverage()

	// Send completion progress update
	if progressChan != nil {
		progressChan <- ProgressUpdate{
			Status:   "complete",
			Message:  "Done!",
			Progress: 1.0,
		}
	}

	return student, nil
}

// Helper function to parse float values
func parseFloat(s string) float64 {
	var result float64
	_, err := fmt.Sscanf(strings.Replace(s, ",", ".", 1), "%f", &result)
	if err != nil {
		return 0
	}
	return result
}

// Helper function to parse dates
func parseDate(s string) time.Time {
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return time.Time{}
	}
	return t
}
