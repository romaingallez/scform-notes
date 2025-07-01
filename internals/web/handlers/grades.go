package handlers

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"scrapping/internals/scform"

	"github.com/gofiber/fiber/v2"
)

// GradeHandler holds the state and methods for handling grade-related requests
type GradeHandler struct {
	currentStudent *scform.Student
}

// NewGradeHandler creates a new instance of GradeHandler
func NewGradeHandler() *GradeHandler {
	return &GradeHandler{}
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

		// Add panic recovery
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic in grade retrieval goroutine: %v", r)
				BroadcastProgress(map[string]interface{}{
					"status":   "error",
					"message":  fmt.Sprintf("Panic occurred: %v", r),
					"progress": 1.0,
				})
			}
		}()

		// Start a goroutine to handle progress updates
		go func() {
			for progress := range progressChan {
				BroadcastProgress(progress)
			}
		}()

		student, err := scform.GetStudentGrades(scformURL, username, password, progressChan)
		if err != nil {
			log.Printf("Error getting grades: %v", err)
			BroadcastProgress(map[string]interface{}{
				"status":   "error",
				"message":  "Error: " + err.Error(),
				"progress": 1.0,
			})
			return
		}
		h.currentStudent = student
	}()

	// Return success response immediately
	return c.JSON(fiber.Map{
		"status":  "processing",
		"message": "Grade retrieval started",
	})
}

// HandleSearch handles the search and sort functionality
func (h *GradeHandler) HandleSearch(c *fiber.Ctx) error {
	if h.currentStudent == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No grades data available",
		})
	}

	query := strings.ToLower(c.Query("q"))
	sortBy := c.Query("sort")
	sortDir := c.Query("dir")

	// Create a copy of the student data
	filteredStudent := &scform.Student{
		Name:         h.currentStudent.Name,
		TotalAverage: h.currentStudent.TotalAverage,
		Grades:       []scform.Course{},
	}

	// Filter courses
	for _, course := range h.currentStudent.Grades {
		if query == "" || strings.Contains(strings.ToLower(course.Name), query) {
			filteredStudent.Grades = append(filteredStudent.Grades, course)
		}
	}

	// Sort grades within each course if requested
	for i := range filteredStudent.Grades {
		if sortBy != "" {
			sort.Slice(filteredStudent.Grades[i].Grades, func(a, b int) bool {
				gradeA := filteredStudent.Grades[i].Grades[a]
				gradeB := filteredStudent.Grades[i].Grades[b]

				isAsc := sortDir != "desc"

				switch sortBy {
				case "title":
					if isAsc {
						return strings.ToLower(gradeA.Title) < strings.ToLower(gradeB.Title)
					}
					return strings.ToLower(gradeA.Title) > strings.ToLower(gradeB.Title)
				case "grade":
					if isAsc {
						return gradeA.Value < gradeB.Value
					}
					return gradeA.Value > gradeB.Value
				case "coef":
					if isAsc {
						return gradeA.Coefficient < gradeB.Coefficient
					}
					return gradeA.Coefficient > gradeB.Coefficient
				case "date":
					if isAsc {
						return gradeA.Date.Before(gradeB.Date)
					}
					return gradeA.Date.After(gradeB.Date)
				case "type":
					if isAsc {
						return strings.ToLower(gradeA.Type) < strings.ToLower(gradeB.Type)
					}
					return strings.ToLower(gradeA.Type) > strings.ToLower(gradeB.Type)
				}
				return false
			})
		}
	}

	return c.Render("partials/grades", fiber.Map{
		"Student": filteredStudent,
		"SortBy":  sortBy,
		"SortDir": sortDir,
	}, "")
}

// HandlePrint renders the print-friendly version of the grades
func (h *GradeHandler) HandlePrint(c *fiber.Ctx) error {
	if h.currentStudent == nil {
		return c.Redirect("/")
	}

	// Get current year for the academic year display
	currentYear := time.Now().Year()
	academicYear := fmt.Sprintf("%d-%d", currentYear-1, currentYear)

	return c.Render("print", fiber.Map{
		"Student":      h.currentStudent,
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
	if h.currentStudent == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No grades data available",
		})
	}

	jsonData, err := scform.ExportToJSON(h.currentStudent)
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
	if h.currentStudent == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No grades data available",
		})
	}

	f, err := scform.ExportToExcel(h.currentStudent)
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
