{{if .Student}}
<div class="space-y-6 pb-8 bg-gray-200 min-h-screen" x-data="gradesTable()" x-init="loadGrades()">
    <!-- Student Info Card -->
    <div class="bg-gray-200 shadow-sm rounded-lg p-4">
        <h2 class="text-lg font-semibold text-gray-900">Moyenne Générale: {{printf "%.2f" .Student.TotalAverage}}/20</h2>
        <div class="flex items-center justify-between mt-2">
            <span class="text-sm text-gray-600">Étudiant: {{.Student.Name}}</span>
            <div class="flex items-center space-x-2">
                <span class="text-sm text-gray-600">Total: <span x-text="totalGrades"></span> notes</span>
                <button @click="loadGrades()" class="btn btn-sm btn-outline">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                </button>
            </div>
        </div>
    </div>

    <!-- Excel-like Table -->
    <div class="bg-white shadow-lg rounded-lg overflow-hidden flex-1">
        <!-- Table Controls -->
        <div class="px-4 py-3 bg-gray-50 border-b border-gray-200 flex items-center justify-between">
            <div class="flex items-center space-x-4">
                <span class="text-sm text-gray-600">
                    Affichage de <span x-text="(currentPage - 1) * itemsPerPage + 1"></span> à 
                    <span x-text="Math.min(currentPage * itemsPerPage, filteredCourses.length)"></span> 
                    sur <span x-text="filteredCourses.length"></span> matières
                </span>
            </div>
            <div class="flex items-center space-x-2">
                <label class="text-sm text-gray-600">Par page:</label>
                <select x-model="itemsPerPage" @change="updatePagination()" class="select select-sm select-bordered">
                    <option value="10">10</option>
                    <option value="25" selected>25</option>
                    <option value="50">50</option>
                    <option value="100">100</option>
                </select>
            </div>
        </div>

        <div class="overflow-x-auto w-full">
            <table class="min-w-full table-auto divide-y divide-gray-200">
                <thead class="bg-gray-50 sticky top-0 z-10">
                    <tr>
                        <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors w-1/4"
                            @click="sortBy('course')">
                            <div class="flex items-center space-x-1">
                                <span>Matière</span>
                                <div class="flex flex-col">
                                    <svg class="w-3 h-3" :class="{'text-blue-600': sortField === 'course' && sortDirection === 'asc'}" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
                                    </svg>
                                    <svg class="w-3 h-3 -mt-1" :class="{'text-blue-600': sortField === 'course' && sortDirection === 'desc'}" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                                    </svg>
                                </div>
                            </div>
                        </th>
                        <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors w-1/6"
                            @click="sortBy('courseAvg')">
                            <div class="flex items-center space-x-1">
                                <span>Moyenne</span>
                                <div class="flex flex-col">
                                    <svg class="w-3 h-3" :class="{'text-blue-600': sortField === 'courseAvg' && sortDirection === 'asc'}" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
                                    </svg>
                                    <svg class="w-3 h-3 -mt-1" :class="{'text-blue-600': sortField === 'courseAvg' && sortDirection === 'desc'}" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                                    </svg>
                                </div>
                            </div>
                        </th>
                        <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors w-1/12"
                            @click="sortBy('gradeCount')">
                            <div class="flex items-center space-x-1">
                                <span>Notes</span>
                                <div class="flex flex-col">
                                    <svg class="w-3 h-3" :class="{'text-blue-600': sortField === 'gradeCount' && sortDirection === 'asc'}" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
                                    </svg>
                                    <svg class="w-3 h-3 -mt-1" :class="{'text-blue-600': sortField === 'gradeCount' && sortDirection === 'desc'}" fill="currentColor" viewBox="0 0 20 20">
                                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                                    </svg>
                                </div>
                            </div>
                        </th>
                        <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-1/2">
                            <span>Détails des Notes</span>
                        </th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    <!-- Use a single template for the entire course block -->
                    <template x-for="(course, courseIndex) in paginatedCourses" :key="course.course">
                        <tr class="course-group">
                            <td colspan="4" class="p-0">
                                <table class="w-full">
                                    <!-- Course Header Row -->
                                    <tr class="bg-blue-50 border-b-2 border-blue-200">
                                        <td class="px-4 py-4 text-sm font-bold text-blue-900 w-1/4" x-text="course.course"></td>
                                        <td class="px-4 py-4 text-sm font-bold text-blue-900 w-1/6">
                                            <span x-text="course.courseAvg + '/20'"></span>
                                        </td>
                                        <td class="px-4 py-4 text-sm font-bold text-blue-900 w-1/12" x-text="course.gradeCount"></td>
                                        <td class="px-4 py-4 text-sm text-blue-900 w-1/2">
                                            <span class="font-medium">Détails des évaluations</span>
                                        </td>
                                    </tr>
                                    
                                    <!-- Grade Sub-rows -->
                                    <template x-for="(grade, gradeIndex) in course.grades" :key="course.course + grade.title + grade.date">
                                        <tr class="transition-colors hover:bg-gray-50" :class="{
                                            'bg-white': gradeIndex % 2 === 0,
                                            'bg-gray-50': gradeIndex % 2 === 1
                                        }">
                                            <td class="px-8 py-3 text-sm text-gray-600 w-1/4">
                                                <div class="flex items-center">
                                                    <div class="w-2 h-2 bg-gray-400 rounded-full mr-3"></div>
                                                    <span x-text="grade.title"></span>
                                                </div>
                                            </td>
                                            <td class="px-4 py-3 text-sm text-gray-900 w-1/6">
                                                <span class="font-medium" x-text="grade.value + '/' + grade.outOf"></span>
                                            </td>
                                            <td class="px-4 py-3 text-sm text-gray-900 w-1/12" x-text="grade.coefficient"></td>
                                            <td class="px-4 py-3 text-sm text-gray-600 w-1/2">
                                                <div class="grid grid-cols-2 gap-2 text-xs">
                                                    <div>
                                                        <span class="font-medium">Date:</span> <span x-text="grade.dateFormatted"></span>
                                                    </div>
                                                    <div>
                                                        <span class="font-medium">Type:</span> <span x-text="grade.type"></span>
                                                    </div>
                                                    <div x-show="grade.remarks" class="col-span-2">
                                                        <span class="font-medium">Remarques:</span> <span x-text="grade.remarks"></span>
                                                    </div>
                                                    <div x-show="grade.observation" class="col-span-2">
                                                        <span class="font-medium">Observation:</span> <span x-text="grade.observation"></span>
                                                    </div>
                                                </div>
                                            </td>
                                        </tr>
                                    </template>
                                </table>
                            </td>
                        </tr>
                    </template>
                </tbody>
            </table>
        </div>
        
        <!-- Pagination -->
        <div x-show="totalPages > 1" class="px-4 py-3 bg-gray-50 border-t border-gray-200 flex items-center justify-between">
            <div class="flex items-center space-x-2">
                <button @click="previousPage()" :disabled="currentPage === 1" 
                        class="btn btn-sm btn-outline" :class="{'opacity-50 cursor-not-allowed': currentPage === 1}">
                    Précédent
                </button>
                <span class="text-sm text-gray-600">
                    Page <span x-text="currentPage"></span> sur <span x-text="totalPages"></span>
                </span>
                <button @click="nextPage()" :disabled="currentPage === totalPages" 
                        class="btn btn-sm btn-outline" :class="{'opacity-50 cursor-not-allowed': currentPage === totalPages}">
                    Suivant
                </button>
            </div>
            <div class="flex items-center space-x-2">
                <span class="text-sm text-gray-600">Aller à:</span>
                <input type="number" x-model.number="currentPage" min="1" :max="totalPages" 
                       class="input input-sm input-bordered w-16 text-center">
            </div>
        </div>
        
        <!-- Loading State -->
        <div x-show="loading" class="flex justify-center items-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <span class="ml-2 text-gray-600">Chargement...</span>
        </div>
        
        <!-- Empty State -->
        <div x-show="!loading && filteredCourses.length === 0" class="text-center py-8 text-gray-500">
            Aucune matière trouvée
        </div>
    </div>
</div>


{{else}}
<div class="text-center text-red-600 bg-gray-200 p-4">
    Aucune note disponible. Veuillez vérifier vos identifiants et réessayer.
</div>
{{end}} 