package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kevindurb/planner/internal/db"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/middleware"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

var exercisesPaths = Paths{"exercises"}

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

func (h *ExercisesHandler) Route(r chi.Router) {
	r.Get("/", ghttp.Adapt(h.list))
	r.Get("/new", ghttp.Adapt(h.new))
	r.Post("/", h.create)

	r.Route("/{exercise_id}", func(r chi.Router) {
		r.Use(h.exerciseCtx())
		r.Get("/", ghttp.Adapt(h.show))
		r.Get("/edit", ghttp.Adapt(h.edit))
		r.Post("/", h.update)
		r.Post("/delete", h.delete)
	})
}

func (h *ExercisesHandler) exerciseCtx() func(http.Handler) http.Handler {
	return middleware.EntityCtx(func(r *http.Request) (db.Exercise, error) {
		id, _ := pathInt(r, "exercise_id")
		userID := h.sm.UserID(r.Context())
		return h.queries.GetExerciseByID(r.Context(), db.GetExerciseByIDParams{
			ID:     id,
			UserID: userID,
		})
	})
}

func (h *ExercisesHandler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	exercise := middleware.FromContext[db.Exercise](r.Context())
	return Layout(
		H1(Text(exercise.Name)),
		A(Href(exercisesPaths.Edit(exercise.ID)), Text("Edit")),
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
	exercise := middleware.FromContext[db.Exercise](r.Context())
	return Layout(
		H1(Text("Edit "+exercise.Name)),
		Form(
			Method("POST"),
			Action(exercisesPaths.Show(exercise.ID)),
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

	http.Redirect(w, r, exercisesPaths.Show(exercise.ID), http.StatusFound)
}

func (h *ExercisesHandler) update(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	exercise := middleware.FromContext[db.Exercise](r.Context())
	var data updateExerciseBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.queries.UpdateExercise(r.Context(), db.UpdateExerciseParams{
		ID:     exercise.ID,
		UserID: userID,
		Name:   data.Name,
	})

	http.Redirect(w, r, exercisesPaths.Show(exercise.ID), http.StatusFound)
}

func (h *ExercisesHandler) delete(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	exercise := middleware.FromContext[db.Exercise](r.Context())
	h.queries.DeleteExerciseByID(r.Context(), db.DeleteExerciseByIDParams{
		ID:     exercise.ID,
		UserID: userID,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
