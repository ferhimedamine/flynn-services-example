package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/flynn/flynn/discoverd/client"
)

var service = "services-example"

func main() {
	port := os.Getenv("PORT")
	addr := ":" + port

	if err := discoverd.Register(service, addr); err != nil {
		log.Fatal(err)
	}

	set, err := discoverd.NewServiceSet(service)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		data := map[string]interface{}{
			"Port":     port,
			"Services": set.Services(),
		}
		html.Execute(w, data)
	})

	log.Println("Listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

var html = template.Must(template.New("html").Funcs(helpers).Parse(`
<html>
  <head>
    <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap.min.css">

    <style type="text/css">
      body {
        padding-top: 20px;
      }

      .hello {
        background-color: #eee;
        font-size: 70px;
        font-weight: 500;
        padding: 30px;
        text-align: center;
      }
    </style>
  </head>

  <body>
    <div class="container">
      <div class="hello">
        <p>Hello from port {{ .Port }}!</p>

        <p>I can see {{ len .Services }} online {{ pluralize (len .Services) "service" }}:</p>

        {{ range .Services }}
          <p>{{ .Addr }}</p>
        {{ end }}
      </div>
    </div>

    <script type="text/javascript">
      window.setTimeout(function() { document.location.reload(true) }, 2000)
    </script>
  </body>
</html>
`[1:]))

var helpers = template.FuncMap{
	"pluralize": func(count int, singular string) string {
		if count == 1 {
			return singular
		}
		return singular + "s"
	},
}
