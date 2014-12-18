package main

import (
  "net/http"
  "fmt"
  "html/template"
  "chitchat/data"
)

// GET /err?msg=
// shows the error message page
func err(writer http.ResponseWriter, request *http.Request) {
  vals := request.URL.Query()
  _, err := session(writer, request)
  var t *template.Template
  if err != nil {
    t = parseTemplateFiles("layout", "public.navbar", "error")
  } else {
    t = parseTemplateFiles("layout", "private.navbar", "error")
  }  
  t.Execute(writer, vals.Get("msg")) 
}

// GET /
// index page
func index(writer http.ResponseWriter, request *http.Request) {
  sess, err := session(writer, request)
  var t *template.Template
  if err != nil {
    t = parseTemplateFiles("layout", "public.navbar", "index")
  } else {
    fmt.Println(sess)
    t = parseTemplateFiles("layout", "private.navbar", "index")
  }  

  threads, err := data.Threads(); if err != nil {
    error_message(writer, request, "Cannot get threads")
  } else {
    t.Execute(writer, threads)
  }  
}