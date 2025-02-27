package main

import (
	"log"
	"os"
	"scrapping/internals/scform"
	"scrapping/internals/utils"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {
	// setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func main() {

	utils.InitAssets()

	// Create a new engine
	engine := html.New("./views", ".html")

	// Create new Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files
	app.Static("/static", "./static")

	// Store student data globally (in production, use a proper state management)
	var currentStudent *scform.Student

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		// Get default credentials from env
		username := os.Getenv("SCFORM_USERNAME")
		password := os.Getenv("SCFORM_PASSWORD")
		scformURL := os.Getenv("SCFORM_URL")

		return c.Render("index", fiber.Map{
			"Title":           "Visualiseur de Notes SCForm",
			"DefaultUsername": username,
			"DefaultPassword": password,
			"DefaultURL":      scformURL,
		})
	})

	app.Post("/grades", func(c *fiber.Ctx) error {
		// Get form data
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

		// Get grades
		student, err := scform.GetStudentGrades(scformURL, username, password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Store student data
		currentStudent = student

		// Return partial HTML for the grades table
		return c.Render("partials/grades", fiber.Map{
			"Student": student,
			"SortBy":  "",
			"SortDir": "",
		}, "")
	})

	// Search and sort endpoint
	app.Get("/search", func(c *fiber.Ctx) error {
		if currentStudent == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "No grades data available",
			})
		}

		query := strings.ToLower(c.Query("q"))
		sortBy := c.Query("sort")
		sortDir := c.Query("dir")

		// Create a copy of the student data
		filteredStudent := &scform.Student{
			Name:         currentStudent.Name,
			TotalAverage: currentStudent.TotalAverage,
			Grades:       []scform.Course{},
		}

		// Filter courses
		for _, course := range currentStudent.Grades {
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
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// scformURL := os.Getenv("SCFORM_URL")
// username := os.Getenv("SCFORM_USERNAME")
// password := os.Getenv("SCFORM_PASSWORD")

// student, err := GetStudentGrades(scformURL, username, password)
// if err != nil {
// 	log.Fatal("Failed to get student grades:", err)
// }

// // Create the JSON output
// jsonData, err := json.MarshalIndent(student, "", "  ")
// if err != nil {
// 	log.Fatal("Failed to marshal student data to JSON:", err)
// }

// err = os.WriteFile("grades.json", jsonData, 0644)
// if err != nil {
// 	log.Fatal("Failed to write grades to file:", err)
// }
