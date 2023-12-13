package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) router() http.Handler {
  router := httprouter.New()

  router.HandlerFunc(http.MethodGet, "/tasks", app.handleTaskList)
  router.HandlerFunc(http.MethodPost, "/tasks", app.handleTaskCreate)

  return router
}
