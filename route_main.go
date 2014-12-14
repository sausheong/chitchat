package main

import (
  "net/http"
  "chitchat/data"
)

// GET /err?msg=
// shows the error message page
func err(writer http.ResponseWriter, request *http.Request) {
  vals := request.URL.Query()
  t := parseTemplateFiles("layout", "public.navbar", "error")
  t.Execute(writer, vals.Get("msg")) 
}

// GET /
// index page
func index(writer http.ResponseWriter, request *http.Request) {
  loggedin(writer, request)
  t := parseTemplateFiles("layout", "public.navbar", "index")
  conversations, err := data.Conversations(); if err != nil {
    error_message(writer, request, "Cannot get conversations")
  } else {
    t.Execute(writer, conversations)
  }  
}