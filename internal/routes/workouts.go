package routes

import "fmt"

type workoutsRoutes struct{}

func (workoutsRoutes) Edit(id int64) string { return fmt.Sprintf("/workouts/%d/edit", id) }
func (workoutsRoutes) Show(id int64) string { return fmt.Sprintf("/workouts/%d", id) }
func (workoutsRoutes) Create() string       { return "/workouts" }
func (workoutsRoutes) New() string          { return "/workouts/new" }
