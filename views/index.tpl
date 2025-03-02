<div class="container mx-auto bg-gray-200 px-4 py-8">
    <div class="card bg-gray-200 shadow-xl max-w-md mx-auto">
        <div class="card-body">
            <h1 class="card-title text-2xl font-bold text-center text-primary mb-6">Visualiseur de Notes SCForm</h1>
            
            <form hx-post="/grades" 
                  hx-target="#grades-container" 
                  hx-indicator="#spinner"
                  hx-on::after-request="initWebSocket()"
                  class="space-y-4">
                
                <div class="form-control w-full">
                    <label for="url" class="label">
                        <span class="label-text">URL SCForm</span>
                    </label>
                    <input type="text" 
                           id="url" 
                           name="url" 
                           placeholder="https://scform.example.com"
                           value="{{if .DefaultURL}}{{.DefaultURL}}{{end}}"
                           class="input input-bordered w-full">
                </div>

                <div class="form-control w-full">
                    <label for="username" class="label">
                        <span class="label-text">Nom d'utilisateur</span>
                    </label>
                    <input type="text" 
                           id="username" 
                           name="username" 
                           placeholder="Votre nom d'utilisateur"
                           value="{{if .DefaultUsername}}{{.DefaultUsername}}{{end}}"
                           class="input input-bordered w-full">
                </div>

                <div class="form-control w-full">
                    <label for="password" class="label">
                        <span class="label-text">Mot de passe</span>
                    </label>
                    <input type="password" 
                           id="password" 
                           name="password" 
                           placeholder="Votre mot de passe"
                           value="{{if .DefaultPassword}}{{.DefaultPassword}}{{end}}"
                           class="input input-bordered w-full">
                </div>

                <button type="submit" 
                        class="btn btn-primary w-full">
                    Obtenir les Notes
                </button>
            </form>

            <div id="progress-container" class="hidden mt-4">
                <div class="w-full bg-gray-300 rounded-full h-2.5">
                    <div id="progress-bar" class="bg-primary h-2.5 rounded-full" style="width: 0%"></div>
                </div>
                <p id="progress-message" class="text-sm text-gray-600 mt-2 text-center"></p>
            </div>

            <div id="spinner" class="htmx-indicator">
                <div class="flex justify-center items-center mt-4">
                    <span class="loading loading-spinner loading-lg text-primary"></span>
                </div>
            </div>
        </div>
    </div>

    <div class="container mx-auto mt-8 mb-8 overflow-x-auto">
        <div class="flex gap-4 mb-4">
            <div class="card bg-gray-200 shadow-xl hidden flex-1" id="search-container">
                <div class="card-body">
                    <div class="relative">
                        <input type="text"
                               id="search-input"
                               name="q"
                               placeholder="Rechercher un cours..."
                               class="input input-bordered w-full pl-10"
                               hx-get="/search"
                               hx-trigger="input changed delay:300ms, search"
                               hx-target="#grades-container"
                               hx-include="this"
                               hx-indicator="#search-spinner">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <svg class="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                            </svg>
                        </div>
                        <div id="search-spinner" class="htmx-indicator absolute inset-y-0 right-0 pr-3 flex items-center">
                            <span class="loading loading-spinner loading-sm text-primary"></span>
                        </div>
                    </div>
                </div>
            </div>
            <button id="print-button"
                    class="btn btn-secondary hidden"
                    onclick="openPrintPopup()">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                </svg>
                Version Imprimable
            </button>
            <button id="download-button"
                    class="btn btn-accent hidden"
                    hx-get="/export"
                    hx-target="this"
                    hx-swap="none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                Télécharger JSON
            </button>
            <button id="excel-download-button"
                    class="btn btn-success hidden"
                    onclick="downloadExcel()">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                Télécharger Excel
            </button>
        </div>
        <div id="grades-container" class="overflow-x-auto"></div>
    </div>
</div>

<style>
    .htmx-indicator{
        display:none;
    }
    .htmx-request .htmx-indicator{
        display:block;
    }
    .htmx-request.htmx-indicator{
        display:block;
    }
</style>

<script>
    let ws = null;

    function initWebSocket() {
        if (ws !== null) {
            ws.close();
        }

        const progressContainer = document.getElementById('progress-container');
        const progressBar = document.getElementById('progress-bar');
        const progressMessage = document.getElementById('progress-message');
        progressContainer.classList.remove('hidden');

        // if we are on http, use ws, if we are on https, use wss   
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';

        const web_socket_url = `${protocol}//${window.location.host}/ws`;

        console.log(web_socket_url);

        // Create WebSocket connection
        ws = new WebSocket(web_socket_url);
        
        ws.onmessage = function(event) {
            const data = JSON.parse(event.data);
            
            // Update progress bar and message
            progressBar.style.width = `${data.progress * 100}%`;
            progressMessage.textContent = data.message;
            
            if (data.status === 'complete') {
                // Reload grades container
                htmx.ajax('GET', '/search', '#grades-container');
                // Hide progress after a delay
                setTimeout(() => {
                    progressContainer.classList.add('hidden');
                }, 1000);
                ws.close();
            }
        };
        
        ws.onclose = function() {
            console.log('WebSocket connection closed');
            ws = null;
        };
        
        ws.onerror = function(error) {
            console.error('WebSocket error:', error);
            progressMessage.textContent = 'Error: Could not connect to WebSocket';
            progressMessage.classList.add('text-red-500');
        };
    }

    document.body.addEventListener('htmx:afterRequest', function(evt) {
        if (evt.detail.target.id === 'grades-container' && evt.detail.successful) {
            document.getElementById('search-container').classList.remove('hidden');
            document.getElementById('print-button').classList.remove('hidden');
            document.getElementById('download-button').classList.remove('hidden');
            document.getElementById('excel-download-button').classList.remove('hidden');
        }
    });

    document.body.addEventListener('htmx:afterRequest', function(evt) {
        if (evt.detail.target.id === 'download-button' && evt.detail.successful) {
            const blob = new Blob([evt.detail.xhr.response], { type: 'application/json' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'grades.json';
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        }
    });

    function downloadExcel() {
        fetch('/export/excel')
            .then(response => {
                const contentDisposition = response.headers.get('Content-Disposition');
                const filename = contentDisposition ? contentDisposition.split('filename=')[1].replace(/["']/g, '') : 'grades.xlsx';
                return response.blob().then(blob => ({ blob, filename }));
            })
            .then(({ blob, filename }) => {
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.style.display = 'none';
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                document.body.removeChild(a);
            })
            .catch(error => console.error('Error downloading Excel file:', error));
    }

    function openPrintPopup() {
        const width = 900;
        const height = 800;
        const left = (window.innerWidth - width) / 2;
        const top = (window.innerHeight - height) / 2;
        
        const popup = window.open(
            '/print',
            'PrintWindow',
            `width=${width},height=${height},left=${left},top=${top},menubar=no,toolbar=no,location=no,status=no`
        );
        
        if (popup) {
            popup.focus();
        } else {
            alert('Please allow popups for this website to use the print feature.');
        }
    }

    // Update HTMX redirect handler to only handle necessary cases
    document.body.addEventListener('htmx:beforeRedirect', function(evt) {
        // Remove the window.close() behavior as it's not needed
    });
</script>