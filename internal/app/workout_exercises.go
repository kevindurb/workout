package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type WorkoutsExercisesHandler struct {
	q  *db.Queries
	sm *SessionManager
	fp *formparser.FormParser
}

func (h *WorkoutsExercisesHandler) Routes(mux *http.ServeMux) {
	mux.Handle("GET /workouts/{workout_id}/exercises/new", ghttp.Adapt(h.new))
	mux.HandleFunc("POST /workouts/{workout_id}/exercises", h.create)
	mux.HandleFunc("POST /workouts_exercises/{id}/delete", h.delete)
}

func (h *WorkoutsExercisesHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("Choose Exercise")),
		Form(
			Method("POST"),
			Action("/workouts_exercises"),
			Button(Type("submit"), Text("Add")),
		),
	), nil
}

func (h *WorkoutsExercisesHandler) create(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	var data createWorkoutBody
	if err := h.fp.Parse(&data, r); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	workout, _ := h.q.CreateWorkout(r.Context(), db.CreateWorkoutParams{
		UserID: userID,
		Name:   data.Name,
	})

	http.Redirect(w, r, fmt.Sprintf("/workouts/%d", workout.ID), http.StatusFound)
}

func (h *WorkoutsExercisesHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, _ := pathInt(r, "id")
	userID := h.sm.UserID(r.Context())
	h.q.DeleteWorkoutByID(r.Context(), db.DeleteWorkoutByIDParams{
		ID:     id,
		UserID: userID,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
