package exercises

import (
	"net/http"

	"github.com/kevindurb/planner/internal/database/sqlcgen"
	formparser "github.com/kevindurb/planner/internal/form_parser"
	. "github.com/kevindurb/planner/internal/html"
	"github.com/kevindurb/planner/internal/middleware"
	"github.com/kevindurb/planner/internal/routes"
	"github.com/kevindurb/planner/internal/session"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type createExerciseBody struct {
	Name string `form:"name,required"`
}

type updateExerciseBody struct {
	Name string `form:"name,required"`
}

type Handler struct {
	q  *sqlcgen.Queries
	sm *session.Manager
	fp *formparser.FormParser
}

func NewHandler(
	q *sqlcgen.Queries,
	sm *session.Manager,
	fp *formparser.FormParser,
) *Handler {
	return &Handler{q, sm, fp}
}

func (h *Handler) show(w http.ResponseWriter, r *http.Request) (Node, error) {
	exercise := middleware.FromContext[sqlcgen.Exercise](r.Context())
	return Layout(
		H1(Text(exercise.Name)),
		A(Href(routes.Exercises.Edit(exercise.ID)), Text("Edit")),
	), nil
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) (Node, error) {
	exercises, _ := h.q.ListAllExercises(r.Context(), h.sm.UserID(r.Context()))
	return Layout(
		H1(Text("Exercises")),
		Map(exercises, func(exercise sqlcgen.Exercise) Node {
			return P(Text(exercise.Name))
		}),
	), nil
}

func (h *Handler) new(w http.ResponseWriter, r *http.Request) (Node, error) {
	return Layout(
		H1(Text("New Exercise")),
		Form(
			Method("POST"),
			Action(routes.Exercises.Create()),
			Label(For("name"), Text("Name")),
			Input(Type("text"), ID("name"), Name("name"), Required()),
			Button(Type("submit"), Text("Create")),
		),
	), nil
}

func (h *Handler) edit(w http.ResponseWriter, r *http.Request) (Node, error) {
	exercise := middleware.FromContext[sqlcgen.Exercise](r.Context())
	return Layout(
		H1(Text("Edit "+exercise.Name)),
		Form(
			Method("POST"),
			Action(routes.Exercises.Show(exercise.ID)),
			Label(For("name"), Text("Name")),
			Input(Type("text"), ID("name"), Name("name"), Value(exercise.Name), Required()),
			Button(Type("submit"), Text("Save")),
		),
	), nil
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	var data createExerciseBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exercise, _ := h.q.CreateExercise(r.Context(), sqlcgen.CreateExerciseParams{
		Name:   data.Name,
		UserID: userID,
	})

	http.Redirect(w, r, routes.Exercises.Show(exercise.ID), http.StatusFound)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	exercise := middleware.FromContext[sqlcgen.Exercise](r.Context())
	var data updateExerciseBody
	if err := h.fp.Parse(&data, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.q.UpdateExercise(r.Context(), sqlcgen.UpdateExerciseParams{
		ID:     exercise.ID,
		UserID: userID,
		Name:   data.Name,
	})

	http.Redirect(w, r, routes.Exercises.Show(exercise.ID), http.StatusFound)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	userID := h.sm.UserID(r.Context())
	exercise := middleware.FromContext[sqlcgen.Exercise](r.Context())
	h.q.DeleteExerciseByID(r.Context(), sqlcgen.DeleteExerciseByIDParams{
		ID:     exercise.ID,
		UserID: userID,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
