package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/russross/blackfriday"
	"github.com/ttacon/go-tigertonic"
	"github.com/yext/glog"
)

var (
	port = flag.String("p", "19000", "port to run merved on")
)

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	mux := tigertonic.NewTrieServeMux()
	mux.HandleFunc("GET", "/gh/gist/{gid}", renderGHGist)
	mux.HandleFunc("POST", "/gh/gist/local", renderLocalGist)
	mux.HandleFunc("GET", "/", index)

	glog.Info("listening on :" + *port + "...")
	go http.ListenAndServe(":"+*port, mux)

	gopath := os.Getenv("GOPATH")
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(
			http.Dir(filepath.Join(gopath, "src/github.com/ttacon/merved/static")))))
	http.ListenAndServe(":19001", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
<head>
  <title>Markdown editor</title>
  <link rel="stylesheet" href="//localhost:19001/static/css/bootstrap.min.css">
  <script src="//localhost:19001/static/js/seli.js" type="text/javascript"></script>
  <script src="//localhost:19001/static/js/req.js" type="text/javascript"></script>
</head>
<body>
  <div>
    <input style="width: 400px;" id="gid" type="text">
    <button id="load-gist" class="btn btn-primary">Retrieve</button>
  </div>
  <textarea id="gist-loc-text">
  </textarea>
  <div id="gist-loc">
  </div>
  <script type="text/javascript">
    $s.ready(function(){
      var ret = $s.find('#load-gist');
      ret.on('click', function() {
        var gid = $s.find('#gid').el.value;
        $req.get('//localhost:19000/gh/gist/' + gid).then(function(resp) {
          var data = JSON.parse(resp.response);
          $s.find('#gist-loc').html(data.html);
          $s.find('#gist-loc-text').text(data.md);
        }, function(code, resp) { console.log(code); console.log(resp); });
      });

      // on change, reparse
      var txt = $s.find('#gist-loc-text');
      txt.on('input', function() {
        $req.post('//localhost:19000/gh/gist/local', JSON.stringify({
          'md': txt.el.value
        })).then(function(resp) {
          $s.find('#gist-loc').html(resp.response);
        }, function(resp) { console.log('yolo'); });
      });
      txt.on('propertychange', function() {
        console.log(txt.el.value);
      });
    });
  </script>
</body>
`))
}

func renderGHGist(w http.ResponseWriter, r *http.Request) {
	gistID := r.URL.Query().Get("gid")

	cli := github.NewClient(nil)
	gist, resp, err := cli.Gists.Get(gistID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Body.Close()

	// for now we ignore files other than the first (will add tabs later)
	var content string
	for _, file := range gist.Files {
		content = *file.Content
		break
	}

	unsafe := blackfriday.MarkdownCommon([]byte(content))
	data, err := json.Marshal(map[string]string{
		"md":   string(content),
		"html": string(unsafe),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(data)
}

func renderLocalGist(w http.ResponseWriter, r *http.Request) {
	var d map[string]string
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unsafe := blackfriday.MarkdownCommon([]byte(d["md"]))
	w.Write(unsafe)
}
