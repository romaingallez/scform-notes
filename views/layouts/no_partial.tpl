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



  <head>
  <body class="flex flex-col h-screen">
    
    <main class="bg-gray-200 flex-grow">
        {{embed}}
    </main>      
    
  </body>

</html>