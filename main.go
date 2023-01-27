package main

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type yamlMap map[string]string

func main() {
	var inputPaths []string
	addr := ":" + os.Args[1]
	if len(os.Args) >= 3 {
		inputPaths = os.Args[2:]
		for i, path := range inputPaths {
			inputPaths[i] = strings.TrimLeft(path, `.\`)
			inputPaths[i] = strings.TrimLeft(inputPaths[i], "./")
		}
	} else {
		inputPaths = []string{"docs"}
	}

	var yamlpaths = make(yamlMap, 8)
	for _, inputPath := range inputPaths {
		file, err := os.Open(inputPath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if fileInfo.IsDir() {
			err = readYamlsFromDir(&yamlpaths, inputPath, file)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		} else {
			yamlpaths[inputPath] = getUriFromPath(inputPath)
		}
	}

	loader := openapi3.NewLoader()

	for ypath, filename := range yamlpaths {
		doc, _ := loader.LoadFromFile(ypath)
		http.HandleFunc("/"+filename+".json", serveJSON(doc))
		http.HandleFunc("/"+filename, servePage(filename))
	}
	http.HandleFunc("/", serveIndex(yamlpaths))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func serveIndex(yamlpaths yamlMap) http.HandlerFunc {

	t, err := template.ParseFiles("template/index.tmpl")
	if err != nil {
		log.Print(err)
		return serverError
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
		w.WriteHeader(200)
		t.Execute(w, yamlpaths)

	}
}

func serverError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func readYamlsFromDir(yamlpaths *yamlMap, basepath string, file *os.File) error {
	files, err := file.ReadDir(0)
	if err != nil {
		return err
	}
	basepath += "/"
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if path.Ext(file.Name()) == ".yaml" {
			ypath := basepath + file.Name()
			(*yamlpaths)[ypath] = getUriFromPath(ypath)
		}
	}
	return nil
}

func serveJSON(doc *openapi3.T) http.HandlerFunc {
	json, _ := doc.MarshalJSON()
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(json)
	}
}

func servePage(filename string) http.HandlerFunc {
	t, err := template.ParseFiles("template/openapi.tmpl")
	if err != nil {
		log.Print(err)
		return serverError
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		t.Execute(w, filename)
		//		w.Write([]byte(`
		//<!DOCTYPE html>
		//<html lang="en">
		//<head>
		//  <meta charset="utf-8" />
		//  <meta name="viewport" content="width=device-width, initial-scale=1" />
		//  <meta
		//    name="description"
		//    content="SwaggerUI"
		//  />
		//  <title>SwaggerUI</title>
		//  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
		//</head>
		//<body>
		//<div id="swagger-ui"></div>
		//<script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js" crossorigin></script>
		//<script>
		//  window.onload = () => {
		//    window.ui = SwaggerUIBundle({
		//      url: 'http://localhost:8080/` + filename + `.json',
		//			  dom_id: '#swagger-ui',
		//			});
		//		  };
		//		</script>
		//		</body>
		//		</html>
		//		`))
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Openapi Not Found</title>
</head>
<body>
<h3><a href="/">Go to Openapi List</a></h3>
		</body>
		</html>`)
	}
}

func getUriFromPath(path string) (uri string) {
	path = strings.TrimRight(path, ".yaml")
	uri = strings.Replace(path, "\\", "/", -1)
	return
}
