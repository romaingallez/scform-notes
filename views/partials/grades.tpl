{{if .Student}}
<div class="space-y-6 pb-8 bg-gray-200">
    <div class="bg-gray-200 shadow-sm rounded-lg p-4">
        <h2 class="text-lg font-semibold text-gray-900">Moyenne Générale: {{printf "%.2f" .Student.TotalAverage}}/20</h2>
    </div>

    {{range .Student.Grades}}
    <div class="bg-gray-200 shadow-sm rounded-lg overflow-hidden">
        <div class="px-4 py-3 bg-gray-50 flex justify-between items-center">
            <h3 class="text-lg font-medium text-gray-900">
                {{.Name}}
            </h3>
            <span class="text-sm text-gray-500">Moyenne: {{printf "%.2f" .Average}}/20</span>
        </div>
        <div class="overflow-x-auto w-full">
            <div class="min-w-full inline-block align-middle">
                <div class="overflow-hidden">
                    <table class="min-w-full divide-y divide-gray-200">
                        <colgroup>
                            <col class="w-[25%]"> <!-- Titre -->
                            <col class="w-[10%]"> <!-- Note -->
                            <col class="w-[10%]"> <!-- Coef -->
                            <col class="w-[10%]"> <!-- Date -->
                            <col class="w-[15%]"> <!-- Type -->
                            <col class="w-[30%]"> <!-- Commentaires -->
                        </colgroup>
                        <thead class="bg-gray-50">
                            <tr>
                                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase truncate">
                                    <button class="flex items-center space-x-1 hover:text-gray-700"
                                            hx-get="/search"
                                            hx-include="#search-input"
                                            hx-target="#grades-container"
                                            hx-vals='{"sort": "title", "dir": "{{if and (eq $.SortBy "title") (eq $.SortDir "asc")}}desc{{else}}asc{{end}}"}'>
                                        <span>Titre</span>
                                        {{if eq $.SortBy "title"}}
                                            {{if eq $.SortDir "asc"}}↑{{else}}↓{{end}}
                                        {{end}}
                                    </button>
                                </th>
                                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase truncate">
                                    <button class="flex items-center space-x-1 hover:text-gray-700"
                                            hx-get="/search"
                                            hx-include="#search-input"
                                            hx-target="#grades-container"
                                            hx-vals='{"sort": "grade", "dir": "{{if and (eq $.SortBy "grade") (eq $.SortDir "asc")}}desc{{else}}asc{{end}}"}'>
                                        <span>Note</span>
                                        {{if eq $.SortBy "grade"}}
                                            {{if eq $.SortDir "asc"}}↑{{else}}↓{{end}}
                                        {{end}}
                                    </button>
                                </th>
                                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase truncate">
                                    <button class="flex items-center space-x-1 hover:text-gray-700"
                                            hx-get="/search"
                                            hx-include="#search-input"
                                            hx-target="#grades-container"
                                            hx-vals='{"sort": "coef", "dir": "{{if and (eq $.SortBy "coef") (eq $.SortDir "asc")}}desc{{else}}asc{{end}}"}'>
                                        <span>Coef</span>
                                        {{if eq $.SortBy "coef"}}
                                            {{if eq $.SortDir "asc"}}↑{{else}}↓{{end}}
                                        {{end}}
                                    </button>
                                </th>
                                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase truncate">
                                    <button class="flex items-center space-x-1 hover:text-gray-700"
                                            hx-get="/search"
                                            hx-include="#search-input"
                                            hx-target="#grades-container"
                                            hx-vals='{"sort": "date", "dir": "{{if and (eq $.SortBy "date") (eq $.SortDir "asc")}}desc{{else}}asc{{end}}"}'>
                                        <span>Date</span>
                                        {{if eq $.SortBy "date"}}
                                            {{if eq $.SortDir "asc"}}↑{{else}}↓{{end}}
                                        {{end}}
                                    </button>
                                </th>
                                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase truncate">
                                    <button class="flex items-center space-x-1 hover:text-gray-700"
                                            hx-get="/search"
                                            hx-include="#search-input"
                                            hx-target="#grades-container"
                                            hx-vals='{"sort": "type", "dir": "{{if and (eq $.SortBy "type") (eq $.SortDir "asc")}}desc{{else}}asc{{end}}"}'>
                                        <span>Type</span>
                                        {{if eq $.SortBy "type"}}
                                            {{if eq $.SortDir "asc"}}↑{{else}}↓{{end}}
                                        {{end}}
                                    </button>
                                </th>
                                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase truncate">Commentaires</th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            {{range .Grades}}
                            <tr class="hover:bg-gray-50">
                                <td class="px-3 py-2 text-sm text-gray-900 truncate" title="{{.Title}}">{{.Title}}</td>
                                <td class="px-3 py-2 text-sm text-gray-900">{{printf "%.2f" .Value}}/{{printf "%.0f" .OutOf}}</td>
                                <td class="px-3 py-2 text-sm text-gray-900">{{printf "%.2f" .Coefficient}}</td>
                                <td class="px-3 py-2 text-sm text-gray-900">{{.Date.Format "02/01/06"}}</td>
                                <td class="px-3 py-2 text-sm text-gray-900 truncate" title="{{.Type}}">{{.Type}}</td>
                                <td class="px-3 py-2 text-sm text-gray-500">
                                    {{if .Remarks}}
                                    <div class="truncate" title="Remarques: {{.Remarks}}">Remarques: {{.Remarks}}</div>
                                    {{end}}
                                    {{if .Observation}}
                                    <div class="truncate" title="Observation: {{.Observation}}">Observation: {{.Observation}}</div>
                                    {{end}}
                                </td>
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
    {{end}}
</div>
{{else}}
<div class="text-center text-red-600 bg-gray-200 p-4">
    Aucune note disponible. Veuillez vérifier vos identifiants et réessayer.
</div>
{{end}} 