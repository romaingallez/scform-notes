<!DOCTYPE html>
<html>
<head>
    <title>Bulletin de Notes - {{.Student.Name}}</title>
    <style>
        @media print {
            @page {
                size: A4;
                margin: 10mm;
            }
            .no-print { display: none !important; }
            body { margin: 0; padding: 0; }
        }
    </style>
</head>
<body class="bg-gray-50 min-h-screen">
    <a href="#" 
       class="fixed top-5 right-5 bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded font-bold shadow-lg z-50 no-print" 
       onclick="window.close(); return false;">Fermer</a>
    
    <div class="max-w-[210mm] mx-auto bg-white shadow-lg rounded-lg my-8 print:my-0 print:shadow-none">
        <div class="p-8 print:p-4">
            <!-- Header -->
            <div class="text-center border-b-2 border-gray-300 pb-6 mb-8 print:pb-4 print:mb-6">
                <h1 class="text-3xl font-bold text-gray-800 mb-2 print:text-2xl">Bulletin de Notes</h1>
                <h2 class="text-lg text-gray-600 print:text-base">Année Académique {{.AcademicYear}}</h2>
            </div>

            <!-- Student Info -->
            <div class="text-center mb-8 print:mb-6">
                <h3 class="text-xl font-bold text-gray-800 bg-gray-100 border-2 border-gray-300 py-3 px-6 rounded-lg print:text-lg print:py-2 print:px-4">{{.Student.Name}}</h3>
            </div>

            <!-- Course Sections -->
            <div id="grades">
            {{range .Student.Grades}}
            <div class="mb-6 print:mb-4 print:break-inside-avoid">
                <table class="w-full border-collapse border border-gray-300 mb-4 print:mb-2">
                    <tr class="bg-gray-700 text-white font-bold text-center">
                        <td colspan="5" class="p-3 border border-gray-600 print:p-2 print:text-sm">{{.Name}} - Moyenne: {{printf "%.2f" .Average}}</td>
                    </tr>
                    <tr class="bg-gray-100">
                        <th class="p-2 border border-gray-300 text-center text-sm font-bold text-gray-700 print:p-1 print:text-xs">Évaluation</th>
                        <th class="p-2 border border-gray-300 text-center text-sm font-bold text-gray-700 print:p-1 print:text-xs">Type</th>
                        <th class="p-2 border border-gray-300 text-center text-sm font-bold text-gray-700 print:p-1 print:text-xs">Date</th>
                        <th class="p-2 border border-gray-300 text-center text-sm font-bold text-gray-700 print:p-1 print:text-xs">Coefficient</th>
                        <th class="p-2 border border-gray-300 text-center text-sm font-bold text-gray-700 print:p-1 print:text-xs">Note</th>
                    </tr>
                    {{range .Grades}}
                    <tr class="hover:bg-blue-50 print:hover:bg-transparent">
                        <td class="p-2 border border-gray-300 text-center text-sm print:p-1 print:text-xs">{{.Title}}</td>
                        <td class="p-2 border border-gray-300 text-center text-sm print:p-1 print:text-xs">{{.Type}}</td>
                        <td class="p-2 border border-gray-300 text-center text-sm print:p-1 print:text-xs">{{.Date.Format "02/01/2006"}}</td>
                        <td class="p-2 border border-gray-300 text-center text-sm print:p-1 print:text-xs">{{.Coefficient}}</td>
                        <td class="p-2 border border-gray-300 text-center text-sm font-medium print:p-1 print:text-xs">{{printf "%.2f" .Value}}</td>
                    </tr>
                    {{end}}
                </table>
            </div>
            {{end}}
            </div>
            <!-- Total Average -->
            <div class="text-center mt-8 print:mt-6">
                <div class="bg-gray-700 text-white font-bold py-4 px-6 rounded-lg text-lg print:py-3 print:px-4 print:text-base">
                    Moyenne Générale: {{printf "%.2f" .Student.TotalAverage}}
                </div>
            </div>
        </div>
    </div>
</body>
</html> 