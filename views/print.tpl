<!DOCTYPE html>
<html>
<head>
    <title>Bulletin de Notes - {{.Student.Name}}</title>
    <style>
        @media print {
            @page {
                size: A4;
                margin: 2cm;
            }
            
            body {
                font-family: Arial, sans-serif;
                line-height: 1.6;
                margin: 0;
                padding: 0;
            }

            .no-print {
                display: none !important;
            }
        }

        .header {
            text-align: center;
            margin-bottom: 2rem;
            padding-bottom: 1rem;
            border-bottom: 2px solid #333;
        }

        .student-info {
            margin-bottom: 2rem;
        }

        .grades-table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 2rem;
        }

        .grades-table th,
        .grades-table td {
            border: 1px solid #333;
            padding: 0.5rem;
            text-align: left;
        }

        .grades-table th {
            background-color: #f0f0f0;
        }

        .course-header {
            background-color: #e0e0e0;
            font-weight: bold;
        }

        .course-average {
            font-weight: bold;
        }

        .total-average {
            text-align: right;
            font-weight: bold;
            font-size: 1.2rem;
            margin-top: 2rem;
            padding-top: 1rem;
            border-top: 2px solid #333;
        }

        .back-button {
            position: fixed;
            top: 1rem;
            right: 1rem;
            padding: 0.5rem 1rem;
            background-color: #4a5568;
            color: white;
            text-decoration: none;
            border-radius: 0.25rem;
        }

        @media screen {
            body {
                max-width: 210mm;
                margin: 2rem auto;
                padding: 2rem;
                background-color: #f0f0f0;
            }

            .print-container {
                background-color: white;
                padding: 2cm;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
        }
    </style>
</head>
<body>
    <a href="/" class="back-button no-print">Retour</a>
    <div class="print-container">
        <div class="header">
            <h1>Bulletin de Notes</h1>
            <h2>Année Académique {{.AcademicYear}}</h2>
        </div>

        <div class="student-info">
            <h3>{{.Student.Name}}</h3>
        </div>

        {{range .Student.Grades}}
        <div class="course-section">
            <table class="grades-table">
                <tr class="course-header">
                    <td colspan="5">{{.Name}} - Moyenne: {{printf "%.2f" .Average}}</td>
                </tr>
                <tr>
                    <th>Évaluation</th>
                    <th>Type</th>
                    <th>Date</th>
                    <th>Coefficient</th>
                    <th>Note</th>
                </tr>
                {{range .Grades}}
                <tr>
                    <td>{{.Title}}</td>
                    <td>{{.Type}}</td>
                    <td>{{.Date.Format "02/01/2006"}}</td>
                    <td>{{.Coefficient}}</td>
                    <td>{{printf "%.2f" .Value}}</td>
                </tr>
                {{end}}
            </table>
        </div>
        {{end}}

        <div class="total-average">
            Moyenne Générale: {{printf "%.2f" .Student.TotalAverage}}
        </div>
    </div>
</body>
</html> 