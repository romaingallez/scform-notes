package scform

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

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
}

func GetStudentGrades(scformURL, username, password string) (*Student, error) {
	// use already installed chrome browser
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	l := launcher.New().Bin(chromePath).Headless(false)

	// Launch and connect to the browser
	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect().NoDefaultDevice()
	defer browser.Close()

	page := browser.MustPage(scformURL)

	page.MustWaitDOMStable()
	page.MustElement("input[id='MainContent_LoginUser_UserName']").MustInput(username)
	page.MustElement("input[id='MainContent_LoginUser_Password']").MustInput(password)

	page.MustWaitStable()

	page.MustEval(`() => {
		LoginBt();
	}`)

	page.MustWaitStable()

	page.MustEval(`() => {
		GoTo('Eleve/MesNotes.aspx');
	}`)

	page.MustWaitStable()

	// MainContent_RadioButtonAffichage_1 click
	page.MustElement("input[id='MainContent_RadioButtonAffichage_1']").MustClick()

	page.MustWaitStable()

	// Get all course tables
	courseTables, err := page.Elements("table.AfficheInfoEnMieux")
	if err != nil {
		return nil, fmt.Errorf("failed to find course tables: %v", err)
	}

	var courses []Course

	for _, table := range courseTables {
		// Extract course name
		nameElement, err := table.Element("span[id*='NomCompletLabel']")
		if err != nil {
			log.Println("Failed to find course name element, skipping table")
			continue
		}

		courseName, err := nameElement.Text()
		if err != nil {
			log.Println("Failed to get course name text, skipping table")
			continue
		}
		courseName = strings.TrimSpace(courseName)

		// Create new course
		course := Course{
			Name:   courseName,
			Grades: []Grade{},
		}

		// Find all grade divs within the table
		gradeDivs, err := table.Elements("div[id='DivNOTE']")
		if err != nil {
			log.Printf("Failed to find grade divs for course %s: %v\n", courseName, err)
			continue
		}

		for _, gradeDiv := range gradeDivs {
			grade := Grade{}

			// Extract grade value and maximum
			if valueSpan, err := gradeDiv.Element("span[id*='Label1']"); err == nil {
				if value, err := valueSpan.Text(); err == nil {
					value = strings.TrimSpace(value)
					if value != "" {
						grade.Value = parseFloat(value)
						grade.OutOf = 20 // Default out of 20
					}
				}
			}

			// Extract coefficient
			if coeffSpan, err := gradeDiv.Element("span[id*='Label3']"); err == nil {
				if coeffText, err := coeffSpan.Text(); err == nil {
					coeffText = strings.TrimSpace(coeffText)
					coeffText = strings.TrimPrefix(coeffText, "coeff. ")
					grade.Coefficient = parseFloat(coeffText)
				}
			}

			// Extract title and date
			if titleSpan, err := gradeDiv.Element("span[id*='Label7']"); err == nil {
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

			// Extract type
			if typeSpan, err := gradeDiv.Element("span[id*='Label8']"); err == nil {
				if typeText, err := typeSpan.Text(); err == nil {
					grade.Type = strings.TrimSpace(typeText)
				}
			}

			// Extract remarks if present
			if remarksSpan, err := gradeDiv.Element("span[id*='Label9']"); err == nil {
				if remarks, err := remarksSpan.Text(); err == nil {
					grade.Remarks = strings.TrimSpace(remarks)
					grade.Remarks = strings.TrimPrefix(grade.Remarks, "Remarque : ")
				}
			}

			// Extract observations if present
			if obsSpan, err := gradeDiv.Element("span[id*='Label10']"); err == nil {
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

	// Create a student and calculate averages
	student := &Student{
		Name:   username,
		Grades: courses,
	}
	student.CalculateTotalAverage()

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
