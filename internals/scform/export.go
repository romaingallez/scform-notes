package scform

import (
	"encoding/json"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// ExportToJSON converts a Student struct to a JSON byte array
func ExportToJSON(student *Student) ([]byte, error) {
	return json.MarshalIndent(student, "", "  ")
}

// ExportToExcel exports student grades to Excel creating a new workbook
func ExportToExcel(student *Student) (*excelize.File, error) {
	// Create a new Excel file
	f := excelize.NewFile()

	// Get the default sheet name
	sheetName := f.GetSheetName(0)

	// Set column widths for better readability
	f.SetColWidth(sheetName, "A", "A", 30) // Student Name
	f.SetColWidth(sheetName, "B", "B", 30) // Subject Name
	f.SetColWidth(sheetName, "C", "C", 30) // Exam Title
	f.SetColWidth(sheetName, "D", "D", 15) // Date
	f.SetColWidth(sheetName, "E", "E", 15) // Type
	f.SetColWidth(sheetName, "F", "F", 10) // Value
	f.SetColWidth(sheetName, "G", "G", 10) // OutOf
	f.SetColWidth(sheetName, "H", "H", 12) // Coefficient
	f.SetColWidth(sheetName, "I", "I", 30) // Remarks
	f.SetColWidth(sheetName, "J", "J", 30) // Observation
	f.SetColWidth(sheetName, "K", "K", 15) // Module Average

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create header style: %v", err)
	}

	// Create data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create data style: %v", err)
	}

	// Create number style (for grades and coefficients)
	numberStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
		NumFmt: 2, // Built-in number format for 2 decimal places
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create number style: %v", err)
	}

	// Set headers
	headers := []string{
		"Student Name",
		"Subject Name",
		"Exam Title",
		"Date",
		"Type",
		"Value",
		"Out Of",
		"Coefficient",
		"Remarks",
		"Observation",
		"Module Average",
	}

	// Write headers
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Start from row 2 for data
	currentRow := 2

	// Write data
	for _, course := range student.Grades {
		for _, grade := range course.Grades {
			// Format date as string
			dateStr := ""
			if !grade.Date.IsZero() {
				dateStr = grade.Date.Format("02/01/2006")
			}

			// Prepare row data
			rowData := []interface{}{
				student.Name,
				course.Name,
				grade.Title,
				dateStr,
				grade.Type,
				grade.Value,
				grade.OutOf,
				grade.Coefficient,
				grade.Remarks,
				grade.Observation,
				course.Average,
			}

			// Write row data with appropriate styles
			for col, value := range rowData {
				cell, _ := excelize.CoordinatesToCellName(col+1, currentRow)
				f.SetCellValue(sheetName, cell, value)

				// Apply number style to numeric columns (Value, OutOf, Coefficient, Module Average)
				if col == 5 || col == 6 || col == 7 || col == 10 {
					f.SetCellStyle(sheetName, cell, cell, numberStyle)
				} else {
					f.SetCellStyle(sheetName, cell, cell, dataStyle)
				}
			}
			currentRow++
		}
	}

	// Add a summary section
	summaryRow := currentRow + 1
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", summaryRow), "Total Average")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", summaryRow), student.TotalAverage)

	// Style the summary
	summaryStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E2EFDA"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
		NumFmt: 2, // Built-in number format for 2 decimal places
	})

	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", summaryRow), fmt.Sprintf("A%d", summaryRow), summaryStyle)
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", summaryRow), fmt.Sprintf("B%d", summaryRow), summaryStyle)

	return f, nil
}
