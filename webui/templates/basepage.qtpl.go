// Code generated by qtc from "basepage.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// This is a base page template. All the other template pages implement this interface.
//

//line templates/basepage.qtpl:3
package templates

//line templates/basepage.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/basepage.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/basepage.qtpl:4
type Page interface {
//line templates/basepage.qtpl:4
	Title() string
//line templates/basepage.qtpl:4
	StreamTitle(qw422016 *qt422016.Writer)
//line templates/basepage.qtpl:4
	WriteTitle(qq422016 qtio422016.Writer)
//line templates/basepage.qtpl:4
	Body() string
//line templates/basepage.qtpl:4
	StreamBody(qw422016 *qt422016.Writer)
//line templates/basepage.qtpl:4
	WriteBody(qq422016 qtio422016.Writer)
//line templates/basepage.qtpl:4
}

// Page prints a page implementing Page interface.

//line templates/basepage.qtpl:11
func StreamPageTemplate(qw422016 *qt422016.Writer, p Page) {
//line templates/basepage.qtpl:11
	qw422016.N().S(`
<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-CuOF+2SnTUfTwSZjCXf01h7uYhfOBuxIhGKPbfEJ3+FqH/s6cIFN9bGr1HmAg4fQ" crossorigin="anonymous">

    <script src="https://cdn.jsdelivr.net/npm/luxon@1.25.0/build/global/luxon.min.js" integrity="sha256-OVk2fwTRcXYlVFxr/ECXsakqelJbOg5WCj1dXSIb+nU=" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js@3.0.0-beta.4/dist/chart.min.js" integrity="sha256-f3G7brzKP7yZvxb4b3eSpi75AN9bq91NvjpkvzM5ODw=" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/chartjs-adapter-luxon@0.2.2/dist/chartjs-adapter-luxon.min.js" integrity="sha256-bgbnCTiuk9rPHmlLrX1soTSIxQJs26agg9kSWIhdcfc=" crossorigin="anonymous"></script>

    <title>`)
//line templates/basepage.qtpl:24
	p.StreamTitle(qw422016)
//line templates/basepage.qtpl:24
	qw422016.N().S(` I guess</title>


    `)
//line templates/basepage.qtpl:30
	qw422016.N().S(`

    <style>
      .test {
          font-size: 80%;
      }

      .subtest {
          font-size: 80%;
      }
    </style>

    <script src="https://kit.fontawesome.com/d6f03257be.js" crossorigin="anonymous"></script>
  </head>
  <body>
    <nav class="navbar navbar-expand-lg navbar-light bg-light mb-3">
      <div class="container-fluid">
        <a class="navbar-brand mb-0 h1" href="/">Tester</a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav">
            <li class="nav-item">
              <a class="nav-link" href="/labels">Labels</a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="/runs">Runs</a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <main role="main" class="container-fluid">
      `)
//line templates/basepage.qtpl:65
	p.StreamBody(qw422016)
//line templates/basepage.qtpl:65
	qw422016.N().S(`
    </main>

    <script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js" integrity="sha384-9/reFTGAW83EW2RDu2S0VKaIzap3H66lZH81PoYlFhbGU+6BZp6G7niu735Sk7lN" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-alpha3/dist/js/bootstrap.min.js" integrity="sha384-t6I8D5dJmMXjCsRLhSzCltuhNZg6P10kE0m0nAncLUjH6GeYLhRU1zfLoW3QNQDF" crossorigin="anonymous"></script>
    <script>
      var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-toggle="tooltip"]'))
      var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl)
      })
      var popoverTriggerList = [].slice.call(document.querySelectorAll('[data-toggle="popover"]'))
      var popoverList = popoverTriggerList.map(function (popoverTriggerEl) {
        return new bootstrap.Popover(popoverTriggerEl)
      })
    </script>
  </body>
</html>


`)
//line templates/basepage.qtpl:84
}

//line templates/basepage.qtpl:84
func WritePageTemplate(qq422016 qtio422016.Writer, p Page) {
//line templates/basepage.qtpl:84
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/basepage.qtpl:84
	StreamPageTemplate(qw422016, p)
//line templates/basepage.qtpl:84
	qt422016.ReleaseWriter(qw422016)
//line templates/basepage.qtpl:84
}

//line templates/basepage.qtpl:84
func PageTemplate(p Page) string {
//line templates/basepage.qtpl:84
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/basepage.qtpl:84
	WritePageTemplate(qb422016, p)
//line templates/basepage.qtpl:84
	qs422016 := string(qb422016.B)
//line templates/basepage.qtpl:84
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/basepage.qtpl:84
	return qs422016
//line templates/basepage.qtpl:84
}

// Base page implementation. Other pages may inherit from it if they need
// overriding only certain Page methods

//line templates/basepage.qtpl:88
type BasePage struct{}

//line templates/basepage.qtpl:89
func (p *BasePage) StreamTitle(qw422016 *qt422016.Writer) {
//line templates/basepage.qtpl:89
	qw422016.N().S(`Tester`)
//line templates/basepage.qtpl:89
}

//line templates/basepage.qtpl:89
func (p *BasePage) WriteTitle(qq422016 qtio422016.Writer) {
//line templates/basepage.qtpl:89
	qw422016 := qt422016.AcquireWriter(qq422016)
//line templates/basepage.qtpl:89
	p.StreamTitle(qw422016)
//line templates/basepage.qtpl:89
	qt422016.ReleaseWriter(qw422016)
//line templates/basepage.qtpl:89
}

//line templates/basepage.qtpl:89
func (p *BasePage) Title() string {
//line templates/basepage.qtpl:89
	qb422016 := qt422016.AcquireByteBuffer()
//line templates/basepage.qtpl:89
	p.WriteTitle(qb422016)
//line templates/basepage.qtpl:89
	qs422016 := string(qb422016.B)
//line templates/basepage.qtpl:89
	qt422016.ReleaseByteBuffer(qb422016)
//line templates/basepage.qtpl:89
	return qs422016
//line templates/basepage.qtpl:89
}
