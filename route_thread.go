package main

import (
  "fmt"
  "net/http"
  "html/template"
  "chitchat/data"
)


// GET /threads/new
// Show the new thread form page
func newThread(writer http.ResponseWriter, request *http.Request) {
  _, err := session(writer, request)
  if err != nil {
    http.Redirect(writer, request, "/login", 302)
  } else {
    t := parseTemplateFiles("layout", "private.navbar", "new.thread")
    t.Execute(writer, nil)      
  }
}

// POST /signup
// Create the user account
func createThread(writer http.ResponseWriter, request *http.Request) {
  sess, err := session(writer, request)
  if err != nil {
    http.Redirect(writer, request, "/login", 302)
  } else {    
    err = request.ParseForm()
    if err != nil {
      danger(err, "Cannot parse form")
    }  
    user, err := sess.User(); if err != nil {
      danger(err, "Cannot get user from session")
    }
    topic := request.PostFormValue("topic")
    if _, err := user.CreateThread(topic); err != nil {
      danger(err, "Cannot create thread")
    }
    http.Redirect(writer, request, "/", 302)        
  }
}

// GET /thread/read
// Show the details of the thread, including the posts and the form to write a post
func readThread(writer http.ResponseWriter, request *http.Request) {
  vals := request.URL.Query()
  uuid := vals.Get("id")
  _, err := session(writer, request)
  var t *template.Template
  if err != nil {
    t = parseTemplateFiles("layout", "public.navbar", "thread")
  } else {
    t = parseTemplateFiles("layout", "private.navbar", "thread")    
  }
  thread, err := data.ThreadByUUID(uuid); if err != nil {
    error_message(writer, request, "Cannot read thread")
  } else {
    t.Execute(writer, &thread)  
  }
}

// POST /thread/post
// Create the post
func postThread(writer http.ResponseWriter, request *http.Request) {
  sess, err := session(writer, request)
  if err != nil {
    http.Redirect(writer, request, "/login", 302)
  } else {
    err = request.ParseForm()
    if err != nil {
      danger(err, "Cannot parse form")
    }  
    user, err := sess.User(); if err != nil {
      danger(err, "Cannot get user from session")
    }
    body := request.PostFormValue("body")
    uuid := request.PostFormValue("uuid")
    thread, err := data.ThreadByUUID(uuid); if err != nil {
      error_message(writer, request, "Cannot read thread")
    }
    if _, err := user.CreatePost(thread, body); err != nil {
      danger(err, "Cannot create post")
    }
    url := fmt.Sprint("/thread/read?id=", uuid)
    http.Redirect(writer, request, url , 302)        
  }
}