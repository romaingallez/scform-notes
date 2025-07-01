<!DOCTYPE html>
<html lang="en">
  <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">

      <link rel="stylesheet" href="/assets/tailwind.css">

      <link rel="stylesheet" href="/assets/style.css">

      <script src="/assets/htmx.min.js"></script>

      <script>
            htmx.defineExtension('json-enc', {
        onEvent: function (name, evt) {
            if (name === "htmx:configRequest") {
                evt.detail.headers['Content-Type'] = "application/json";
            }
        },
        
        encodeParameters : function(xhr, parameters, elt) {
            xhr.overrideMimeType('text/json');
            return (JSON.stringify(parameters));
        }
    });

      </script>


      <script defer src="/assets/alpine.js"></script>

      <script src="/assets/hyperscript.js"></script>

      <script>
        // Global Alpine.js functions
        function gradesTable() {
            return {
                courses: [],
                filteredCourses: [],
                loading: false,
                sortField: '',
                sortDirection: 'asc',
                totalGrades: 0,
                searchQuery: '',
                currentPage: 1,
                itemsPerPage: 25,

                get totalPages() {
                    return Math.ceil(this.filteredCourses.length / this.itemsPerPage);
                },

                get paginatedCourses() {
                    const start = (this.currentPage - 1) * this.itemsPerPage;
                    const end = start + this.itemsPerPage;
                    return this.filteredCourses.slice(start, end);
                },

                async loadGrades() {
                    this.loading = true;
                    try {
                        const response = await fetch('/api/grades');
                        const data = await response.json();
                        
                        if (response.ok) {
                            this.courses = data.courses || [];
                            this.filteredCourses = [...this.courses];
                            this.totalGrades = data.total || 0;
                            this.currentPage = 1; // Reset to first page
                        } else {
                            console.error('Error loading grades:', data.error);
                        }
                    } catch (error) {
                        console.error('Error loading grades:', error);
                    } finally {
                        this.loading = false;
                    }
                },

                sortBy(field) {
                    if (this.sortField === field) {
                        this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
                    } else {
                        this.sortField = field;
                        this.sortDirection = 'asc';
                    }

                    this.filteredCourses.sort((a, b) => {
                        let aVal = a[field];
                        let bVal = b[field];

                        // Handle numeric values
                        if (field === 'courseAvg' || field === 'gradeCount') {
                            aVal = parseFloat(aVal);
                            bVal = parseFloat(bVal);
                        }

                        // Handle string values
                        if (typeof aVal === 'string') {
                            aVal = aVal.toLowerCase();
                            bVal = bVal.toLowerCase();
                        }

                        if (this.sortDirection === 'asc') {
                            return aVal > bVal ? 1 : -1;
                        } else {
                            return aVal < bVal ? 1 : -1;
                        }
                    });

                    this.currentPage = 1; // Reset to first page after sorting
                },

                filterGrades(query) {
                    this.searchQuery = query.toLowerCase();
                    this.filteredCourses = this.courses.filter(course => 
                        course.course.toLowerCase().includes(this.searchQuery) ||
                        course.grades.some(grade => 
                            grade.title.toLowerCase().includes(this.searchQuery) ||
                            grade.type.toLowerCase().includes(this.searchQuery)
                        )
                    );
                    this.currentPage = 1; // Reset to first page after filtering
                },

                nextPage() {
                    if (this.currentPage < this.totalPages) {
                        this.currentPage++;
                    }
                },

                previousPage() {
                    if (this.currentPage > 1) {
                        this.currentPage--;
                    }
                }
            }
        }
      </script>



  <head>
  <body class="flex flex-col min-h-screen" x-data="{ searchQuery: '' }">
    {{template "partials/header" .}}
    
    <main class="bg-gray-200 flex-grow">
        {{embed}}
    </main>      
    
    {{template "partials/footer" .}}
  </body>

</html>