package app

import (
	"fmt"
	"net/http"

	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

type createExerciseBody struct {
	Name string `form:"name,required"`
}

type updateExerciseBody struct {
	Name string `form:"name,required"`
}

type ExercisesHandler struct {
	queries *db.Queries
	sm      *SessionManager
	fp      *formparser.FormParser
}

func (h *ExercisesHandler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{id}", ghttp.Adapt(h.show))
	mux.Handle("GET /{id}/edit", ghttp.Adapt(h.edit))
	mux.Handle("GET /", ghttp.Adapt(h.list))
	mux.Handle("GET /new", ghttp.Adapt(h.new))
	mux.HandleFunc("POST /", h.create)
	mux.HandleFunc("POST /{id}", h.update)
	mux.HandleFunc("POST /{id}/delete", h.delete)

	return mux
}

func (h *ExercisesHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	id, _ := pathInt(r, "id")
	exercise, err := h.queries.GetExerciseByID(r.Context(), db.GetExerciseByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text(exercise.Name)),
		A(Href(fmt.Sprintf("/exercises/%d/edit", exercise.ID)), Text("Edit")),
	), nil
}

func (h *ExercisesHandler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	exercises, _ := h.queries.ListAllExercises(r.Context(), h.sm.UserID(r.Context()))
	return Layout(
		H1(Text("Exercises")),
		Map(exercises, func(exercise db.Exercise) Node {
			return P(Text(exercise.Name))
		}),
	), nil
}

func (h *ExercisesHandler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("New Exercise")),
		Form(
			Method("POST"),
			Action("/exercises"),
			Label(For("name"), Text("Name")),
			Input(Type("text"), ID("name"), Name("name"), Required()),
			Button(Type("submit"), Text("Create")),
		),
	), nil
}

func (h *ExercisesHandler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	userID := h.sm.UserID(r.Context())
	id, _ := pathInt(r, "id")
	exercise, err := h.queries.GetExerciseByID(r.Context(), db.GetExerciseByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, StatusCodeError{http.StatusNotFound}
	}
	return Layout(
		H1(Text("Edit "+exercise.Name)),
		Form(
			Method("POST"),
			Action(fmt.Sprintf("/exercises/%d", exercise.ID)),
			Label(For("name"), Text("Name")),
			Input(Type("text"), ID("name"), Name("name"), Value(exercise.Name), Required()),
			Button(Type("submit"), Text("Save")),
		),
	), nil
}

func (h *ExercisesHandler) create(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	var data createExerciseBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exercise, _ := h.queries.CreateExercise(r.Context(), db.CreateExerciseParams{
		Name:   data.Name,
		UserID: userID,
	})

	http.Redirect(w, r, fmt.Sprintf("/exercises/%d", exercise.ID), http.StatusFound)
}

func (h *ExercisesHandler) update(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	id, _ := pathInt(r, "id")
	var data updateExerciseBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.queries.UpdateExercise(r.Context(), db.UpdateExerciseParams{
		ID:     id,
		UserID: userID,
		Name:   data.Name,
	})

	http.Redirect(w, r, fmt.Sprintf("/exercises/%d", id), http.StatusFound)
}

func (h *ExercisesHandler) delete(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	id, _ := pathInt(r, "id")
	h.queries.DeleteExerciseByID(r.Context(), db.DeleteExerciseByIDParams{
		ID:     id,
		UserID: userID,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
