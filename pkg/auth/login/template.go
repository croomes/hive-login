package login

import (
	"html/template"
	"log"
	"net/http"
)

var indexTmpl = template.Must(template.New("index.html").Parse(`<html>
  <head>
    <style>
form  { display: table;      }
p     { display: table-row;  }
label { display: table-cell; }
input { display: table-cell; }
    </style>
  </head>
  <body>
    <form action="/login" method="post">
      <p>
        <label> Authenticate for: </label>
        <input type="text" name="cross_client" placeholder="list of client-ids">
      </p>
      <p>
        <label>Extra scopes: </label>
        <input type="text" name="extra_scopes" placeholder="list of scopes">
      </p>
      <p>
        <label>Connector ID: </label>
        <input type="text" name="connector_id" placeholder="connector id">
      </p>
      <p>
        <label>Request offline access: </label>
        <input type="checkbox" name="offline_access" value="yes" checked>
      </p>
      <p>
	    <input type="submit" value="Login">
      </p>
    </form>
  </body>
</html>`))

func renderIndex(w http.ResponseWriter) {
	err := indexTmpl.Execute(w, nil)
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *template.Error:
		// An ExecError guarantees that Execute has not written to the underlying reader.
		log.Printf("Error rendering template %s: %s", indexTmpl.Name(), err)

		// TODO(ericchiang): replace with better internal server error.
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	default:
		// An error with the underlying write, such as the connection being
		// dropped. Ignore for now.
	}
}
