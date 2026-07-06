package routes

import "fmt"

type exercisesRoutes struct{}

func (exercisesRoutes) Edit(id int64) string { return fmt.Sprintf("/exercises/%d/edit", id) }
func (exercisesRoutes) Show(id int64) string { return fmt.Sprintf("/exercises/%d", id) }
func (exercisesRoutes) Create() string       { return "/exercises" }
func (exercisesRoutes) New() string          { return "/exercises/new" }
