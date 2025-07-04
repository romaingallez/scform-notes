package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"scrapping/internals/scform"
	"scrapping/internals/web/session"

	"github.com/gofiber/fiber/v2"
)

// temporaryStudentData holds student data temporarily before it's saved to session
var (
	tempStudentData = make(map[string]*scform.Student)
	tempDataMux     sync.RWMutex
)

// GradeHandler holds the state and methods for handling grade-related requests
type GradeHandler struct {
	sessionManager *session.Manager
}

// NewGradeHandler creates a new instance of GradeHandler
func NewGradeHandler(sessionManager *session.Manager) *GradeHandler {
	return &GradeHandler{
		sessionManager: sessionManager,
	}
}

// getSessionID helper function to get session ID from Fiber context
func (h *GradeHandler) getSessionID(c *fiber.Ctx) string {
	return h.sessionManager.GetSessionID(c)
}

// getCurrentStudent retrieves the current student from session or temporary storage
func (h *GradeHandler) getCurrentStudent(c *fiber.Ctx) *scform.Student {
	sessionID := h.getSessionID(c)
	if sessionID == "" {
		return nil
	}

	// First check temporary storage
	tempDataMux.RLock()
	if student, exists := tempStudentData[sessionID]; exists {
		tempDataMux.RUnlock()
		// Move from temp storage to session storage
		h.setCurrentStudent(c, student)
		// Remove from temp storage
		tempDataMux.Lock()
		delete(tempStudentData, sessionID)
		tempDataMux.Unlock()
		return student
	}
	tempDataMux.RUnlock()

	// Then check session storage
	sess, err := h.sessionManager.Store.Get(c)
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		return nil
	}

	studentData := sess.Get("currentStudent")
	if studentData == nil {
		return nil
	}

	// Try to unmarshal the student data
	var student scform.Student
	if studentBytes, ok := studentData.([]byte); ok {
		if err := json.Unmarshal(studentBytes, &student); err != nil {
			log.Printf("Failed to unmarshal student data: %v", err)
			return nil
		}
		return &student
	}

	return nil
}

// setCurrentStudent stores the current student in session
func (h *GradeHandler) setCurrentStudent(c *fiber.Ctx, student *scform.Student) error {
	sess, err := h.sessionManager.Store.Get(c)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}

	// Marshal student data
	studentBytes, err := json.Marshal(student)
	if err != nil {
		return fmt.Errorf("failed to marshal student data: %v", err)
	}

	sess.Set("currentStudent", studentBytes)
	return sess.Save()
}

// setTempStudentData stores student data temporarily by session ID
func (h *GradeHandler) setTempStudentData(sessionID string, student *scform.Student) {
	tempDataMux.Lock()
	defer tempDataMux.Unlock()
	tempStudentData[sessionID] = student
}

// HandleIndex renders the index page with default credentials
func (h *GradeHandler) HandleIndex(c *fiber.Ctx) error {
	username := os.Getenv("SCFORM_USERNAME")
	password := os.Getenv("SCFORM_PASSWORD")
	scformURL := os.Getenv("SCFORM_URL")

	// Create a map with default values
	data := fiber.Map{
		"Title": "Visualiseur de Notes SCForm",
	}

	// Only add credentials if they exist in environment
	if username != "" {
		data["DefaultUsername"] = username
	}
	if password != "" {
		data["DefaultPassword"] = password
	}
	if scformURL != "" {
		data["DefaultURL"] = scformURL
	}

	return c.Render("index", data)
}

// HandleGrades processes the grade retrieval request
func (h *GradeHandler) HandleGrades(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	if sessionID == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get session ID",
		})
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	scformURL := c.FormValue("url")

	// If any field is empty, use environment variables
	if username == "" {
		username = os.Getenv("SCFORM_USERNAME")
	}
	if password == "" {
		password = os.Getenv("SCFORM_PASSWORD")
	}
	if scformURL == "" {
		scformURL = os.Getenv("SCFORM_URL")
	}

	// Create a channel for progress updates
	progressChan := make(chan scform.ProgressUpdate)

	// Start a goroutine to get grades and broadcast progress
	go func() {
		defer close(progressChan)

		// Start a goroutine to handle progress updates for this specific session
		go func() {
			for progress := range progressChan {
				BroadcastProgressToSession(sessionID, progress)
			}
		}()

		// Retry logic with maximum 3 attempts
		maxRetries := 3
		var student *scform.Student
		var err error

		for attempt := 1; attempt <= maxRetries; attempt++ {
			func() {
				// Add panic recovery for each attempt
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Panic in grade retrieval goroutine (attempt %d/%d) for session %s: %v", attempt, maxRetries, sessionID, r)

						if attempt < maxRetries {
							log.Printf("Retrying in 2 seconds... (attempt %d/%d) for session %s", attempt+1, maxRetries, sessionID)
							BroadcastProgressToSession(sessionID, map[string]interface{}{
								"status":   "retrying",
								"message":  fmt.Sprintf("Attempt %d failed, retrying... (Error: %v)", attempt, r),
								"progress": float64(attempt) / float64(maxRetries),
							})
							time.Sleep(2 * time.Second)
						} else {
							BroadcastProgressToSession(sessionID, map[string]interface{}{
								"status":   "error",
								"message":  fmt.Sprintf("All %d attempts failed. Last error: %v", maxRetries, r),
								"progress": 1.0,
							})
						}
					}
				}()

				student, err = scform.GetStudentGrades(scformURL, username, password, progressChan)
			}()

			// If we got a student successfully, break out of retry loop
			if student != nil && err == nil {
				log.Printf("Successfully retrieved grades on attempt %d for session %s", attempt, sessionID)
				break
			}

			// If this is not the last attempt, log and retry
			if attempt < maxRetries {
				log.Printf("Error getting grades (attempt %d/%d) for session %s: %v", attempt, maxRetries, sessionID, err)
				log.Printf("Retrying in 2 seconds... (attempt %d/%d) for session %s", attempt+1, maxRetries, sessionID)
				BroadcastProgressToSession(sessionID, map[string]interface{}{
					"status":   "retrying",
					"message":  fmt.Sprintf("Attempt %d failed, retrying... (Error: %v)", attempt, err),
					"progress": float64(attempt) / float64(maxRetries),
				})
				time.Sleep(2 * time.Second)
			}
		}

		// Check final result
		if err != nil || student == nil {
			log.Printf("All %d attempts failed for session %s. Final error: %v", maxRetries, sessionID, err)
			BroadcastProgressToSession(sessionID, map[string]interface{}{
				"status":   "error",
				"message":  fmt.Sprintf("All %d attempts failed. Final error: %v", maxRetries, err),
				"progress": 1.0,
			})
			return
		}

		// Store student data in temporary storage (will be moved to session on next request)
		h.setTempStudentData(sessionID, student)

		BroadcastProgressToSession(sessionID, map[string]interface{}{
			"status":   "success",
			"message":  "Grades retrieved successfully",
			"progress": 1.0,
		})
	}()

	// Return success response immediately
	return c.JSON(fiber.Map{
		"status":  "processing",
		"message": "Grade retrieval started",
	})
}

// HandleSearch handles the search and sort functionality
func (h *GradeHandler) HandleSearch(c *fiber.Ctx) error {
	if h.getCurrentStudent(c) == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No grades data available",
		})
	}

	query := strings.ToLower(c.Query("q"))
	sortBy := c.Query("sort")
	sortDir := c.Query("dir")

	// Create a copy of the student data
	filteredStudent := &scform.Student{
		Name:         h.getCurrentStudent(c).Name,
		TotalAverage: h.getCurrentStudent(c).TotalAverage,
		Grades:       []scform.Course{},
	}

	// Filter courses
	for _, course := range h.getCurrentStudent(c).Grades {
		if query == "" || strings.Contains(strings.ToLower(course.Name), query) {
			// Create a copy of the course
			filteredCourse := scform.Course{
				Name:    course.Name,
				Average: course.Average,
				Grades:  []scform.Grade{},
			}

			// Add all grades for this course (no filtering at grade level for now)
			filteredCourse.Grades = append(filteredCourse.Grades, course.Grades...)
			filteredStudent.Grades = append(filteredStudent.Grades, filteredCourse)
		}
	}

	// Sort courses if requested
	if sortBy != "" {
		sort.Slice(filteredStudent.Grades, func(a, b int) bool {
			courseA := filteredStudent.Grades[a]
			courseB := filteredStudent.Grades[b]

			isAsc := sortDir != "desc"

			switch sortBy {
			case "course":
				if isAsc {
					return strings.ToLower(courseA.Name) < strings.ToLower(courseB.Name)
				}
				return strings.ToLower(courseA.Name) > strings.ToLower(courseB.Name)
			case "average":
				if isAsc {
					return courseA.Average < courseB.Average
				}
				return courseA.Average > courseB.Average
			case "gradeCount":
				if isAsc {
					return len(courseA.Grades) < len(courseB.Grades)
				}
				return len(courseA.Grades) > len(courseB.Grades)
			}
			return false
		})
	}

	return c.Render("partials/grades", fiber.Map{
		"Student": filteredStudent,
		"SortBy":  sortBy,
		"SortDir": sortDir,
	}, "")
}

// HandlePrint renders the print-friendly version of the grades
func (h *GradeHandler) HandlePrint(c *fiber.Ctx) error {
	if h.getCurrentStudent(c) == nil {
		return c.Redirect("/")
	}

	// Get current year for the academic year display
	currentYear := time.Now().Year()
	academicYear := fmt.Sprintf("%d-%d", currentYear-1, currentYear)

	return c.Render("print", fiber.Map{
		"Student":      h.getCurrentStudent(c),
		"AcademicYear": academicYear,
	}, "layouts/no_partial")
}

func (h *GradeHandler) HandlePrintDemo(c *fiber.Ctx) error {
	// geneareate a fake student
	student := &scform.Student{
		Name:   "John Doe",
		Grades: []scform.Course{},
	}

	for i := 0; i < 10; i++ {
		student.Grades = append(student.Grades, scform.Course{
			Name:   fmt.Sprintf("Course %d", i),
			Grades: []scform.Grade{},
		})
		for j := 0; j < 3; j++ {
			student.Grades[i].Grades = append(student.Grades[i].Grades, scform.Grade{
				Title: fmt.Sprintf("Midterm %d", j),
				Value: float64(j),
				Date:  time.Now(),
			})
		}
	}

	// Get current year for the academic year display
	currentYear := time.Now().Year()
	academicYear := fmt.Sprintf("%d-%d", currentYear-1, currentYear)

	return c.Render("print", fiber.Map{
		"Student":      student,
		"AcademicYear": academicYear,
	}, "layouts/no_partial")
}

// HandleExport handles the export of grades to JSON
func (h *GradeHandler) HandleExport(c *fiber.Ctx) error {
	if h.getCurrentStudent(c) == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No grades data available",
		})
	}

	jsonData, err := scform.ExportToJSON(h.getCurrentStudent(c))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Disposition", "attachment; filename=grades.json")
	c.Set("Content-Type", "application/json")
	return c.SendStream(bytes.NewReader(jsonData))
}

// HandleExcelExport handles the export of grades to Excel
func (h *GradeHandler) HandleExcelExport(c *fiber.Ctx) error {
	if h.getCurrentStudent(c) == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No grades data available",
		})
	}

	f, err := scform.ExportToExcel(h.getCurrentStudent(c))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=grades.xlsx")

	fReader, err := f.WriteToBuffer()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fSize := fReader.Len()

	return c.SendStream(fReader, fSize)
}

// HandleImport handles the import of grades from JSON file
func (h *GradeHandler) HandleImport(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("json_file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No file uploaded or file upload failed",
		})
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".json") {
		return c.Status(400).JSON(fiber.Map{
			"error": "File must be a JSON file",
		})
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Read file content
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to read file content",
		})
	}

	// Parse JSON into Student struct
	var student scform.Student
	if err := json.Unmarshal(buf.Bytes(), &student); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid JSON format or incompatible structure",
		})
	}

	// Recalculate averages to ensure consistency
	student.CalculateTotalAverage()

	// Set as current student
	h.setCurrentStudent(c, &student)

	// Log successful import
	log.Printf("Successfully imported grades for student: %s", student.Name)

	// Return success response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("Successfully imported grades for %s", student.Name),
		"student": student.Name,
	})
}

// HandleGradesAPI returns grades data as JSON for the Excel-like table
func (h *GradeHandler) HandleGradesAPI(c *fiber.Ctx) error {
	if h.getCurrentStudent(c) == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No grades data available",
		})
	}

	query := strings.ToLower(c.Query("q"))
	sortBy := c.Query("sort")
	sortDir := c.Query("dir")

	// Create a grouped structure by course
	var groupedCourses []map[string]interface{}

	for _, course := range h.getCurrentStudent(c).Grades {
		if query == "" || strings.Contains(strings.ToLower(course.Name), query) {
			// Create course object with its grades
			courseData := map[string]interface{}{
				"course":     course.Name,
				"courseAvg":  course.Average,
				"gradeCount": len(course.Grades),
				"grades":     []map[string]interface{}{},
			}

			// Add grades for this course
			for _, grade := range course.Grades {
				courseData["grades"] = append(courseData["grades"].([]map[string]interface{}), map[string]interface{}{
					"title":         grade.Title,
					"value":         grade.Value,
					"outOf":         grade.OutOf,
					"coefficient":   grade.Coefficient,
					"date":          grade.Date.Format("2006-01-02"),
					"dateFormatted": grade.Date.Format("02/01/06"),
					"type":          grade.Type,
					"remarks":       grade.Remarks,
					"observation":   grade.Observation,
				})
			}

			groupedCourses = append(groupedCourses, courseData)
		}
	}

	// Sort the grouped courses
	if sortBy != "" {
		sort.Slice(groupedCourses, func(a, b int) bool {
			courseA := groupedCourses[a]
			courseB := groupedCourses[b]

			isAsc := sortDir != "desc"

			switch sortBy {
			case "course":
				courseNameA := courseA["course"].(string)
				courseNameB := courseB["course"].(string)
				if isAsc {
					return strings.ToLower(courseNameA) < strings.ToLower(courseNameB)
				}
				return strings.ToLower(courseNameA) > strings.ToLower(courseNameB)
			case "average":
				avgA := courseA["courseAvg"].(float64)
				avgB := courseB["courseAvg"].(float64)
				if isAsc {
					return avgA < avgB
				}
				return avgA > avgB
			case "gradeCount":
				countA := courseA["gradeCount"].(int)
				countB := courseB["gradeCount"].(int)
				if isAsc {
					return countA < countB
				}
				return countA > countB
			}
			return false
		})
	}

	// Calculate total grades across all courses
	totalGrades := 0
	for _, course := range groupedCourses {
		totalGrades += course["gradeCount"].(int)
	}

	return c.JSON(fiber.Map{
		"student": map[string]interface{}{
			"name":         h.getCurrentStudent(c).Name,
			"totalAverage": h.getCurrentStudent(c).TotalAverage,
		},
		"courses": groupedCourses,
		"total":   totalGrades,
	})
}
