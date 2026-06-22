package html

import (
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

func Layout(children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    "Planner",
		Language: "en",
		Head: []g.Node{
			h.Script(
				h.Src("https://unpkg.com/htmx.org@2.0.4"),
				h.Integrity("sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+"),
				h.CrossOrigin("anonymous"),
			),
			h.Script(
				h.Src("https://unpkg.com/htmx-ext-sse@2.2.2"),
				h.Integrity("sha384-fw+eTlCc7suMV/1w/7fr2/PmwElUIt5i82bi+qTiLXvjRXZ2/FkiTNA/w0MhXnGI"),
				h.CrossOrigin("anonymous"),
			),
			h.Script(
				h.Src("https://unpkg.com/idiomorph@0.7.3/dist/idiomorph-ext.min.js"),
				h.Integrity("sha384-szktAZju9fwY15dZ6D2FKFN4eZoltuXiHStNDJWK9+FARrxJtquql828JzikODob"),
				h.CrossOrigin("anonymous"),
			),
			h.Link(
				h.Rel("stylesheet"),
				h.Href("https://fonts.googleapis.com/icon?family=Material+Icons"),
			),
			h.Link(
				h.Rel("stylesheet"),
				h.Href("https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css"),
			),
			h.Script(
				h.Type("module"),
				h.Src("/static/js/main.js"),
			),
			h.Link(
				h.Rel("stylesheet"),
				h.Href("/static/css/main.css"),
			),
		},
		Body: []g.Node{
			hx.Boost("true"),
			hx.Ext("sse,morph"),
			h.Main(
				h.Class("main"),
				g.Group(children),
			),
		},
	})
}
