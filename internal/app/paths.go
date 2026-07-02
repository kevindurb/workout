package app

import (
	"fmt"
)

type PathBuilder struct {
	name string
}

func (p *PathBuilder) List() string {
	return fmt.Sprintf("/%s", p.name)
}

func (p *PathBuilder) Show(id int64) string {
	return fmt.Sprintf("/%s/%d", p.name, id)
}

func (p *PathBuilder) New() string {
	return fmt.Sprintf("/%s/new", p.name)
}

func (p *PathBuilder) Edit(id int64) string {
	return fmt.Sprintf("/%s/%d/edit", id, p.name)
}
