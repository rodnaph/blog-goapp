
package hello

import (
    "fmt"
    "http"
    "appengine"
    "appengine/user"
    "appengine/datastore"
    "template"
    "time"
)

type Picture struct {
    Author string
    Image appengine.BlobKey
    Date datastore.Time
}

func init() {
    // http.HandleFunc( "/", handler )
    http.HandleFunc( "/", root )
    http.HandleFunc( "/post", post )
    http.HandleFunc( "/pic", pic )
}

func root( w http.ResponseWriter, r *http.Request ) {

    c := appengine.NewContext( r )
    q := datastore.NewQuery( "Picture" ).Order( "-Date" ).Limit ( 10 )

    pictures := make( []Picture, 0, 10 )

    if _, err := q.GetAll( c, &pictures); err != nil {
        http.Error( w, err.String(), http.StatusInternalServerError )
        return
    }

    postFormTemplate.Execute( w, pictures )

}

var postFormTemplate = template.MustParse( postFormHTML, nil )

const postFormHTML = `
<html>
<body>

<h1>Pictures</h1>

<ul>
    {.repeated section @}
    <li>
        ITEM: <img src="/pic" />
    </li>
    {.end}
</ul>

<h2>Add New Picture</h2>

<form action="/post" method="post">
  <input type="file" name="picture" />
  <input type="submit" value="Add picture" />
</form>
</body>
</html>
`

func pic( w http.ResponseWriter, r *http.Request ) {

    c := appengine.NewContext( r )
    q := datastore.NewQuery( "Picture" ).Limit( 1 )

    pictures := make( []Picture, 0, 1 )

    q.GetAll( c, &pictures )

    w.Header().Set( "Content-Type", "image/png" )
    
    fmt.Fprint( pictures[0].Image )

}

func post( w http.ResponseWriter, r *http.Request ) {

    c := appengine.NewContext( r )
    g := Picture {
        Image: appengine.BlobKey(r.FormValue( "picture" )),
        Date: datastore.SecondsToTime( time.Seconds()),
    }

    _, err := datastore.Put( c, datastore.NewIncompleteKey("Picture"), &g )
    if err != nil {
        http.Error( w, err.String(), http.StatusInternalServerError )
    }

    http.Redirect( w, r, "/", http.StatusFound )

}

var postTemplate = template.MustParse( postTemplateHTML, nil )

const postTemplateHTML = `
<html>
<body>
  <p>You uploaded:</p>
</body>
</html>
`

func handler( w http.ResponseWriter, r *http.Request ) {

    c := appengine.NewContext( r )
    u := user.Current( c )

    if u == nil {
        redirectToLogin( c, r, w )
        return
    }

    fmt.Fprintf( w, "Hello, %v!", u )

}

func redirectToLogin( c appengine.Context, r *http.Request, w http.ResponseWriter ) {

    url, err := user.LoginURL( c, r.URL.String() )
    if err != nil {
        http.Error( w, err.String(), http.StatusInternalServerError )
        return
    }

    w.Header().Set( "Location", url )
    w.WriteHeader( http.StatusFound )

}
