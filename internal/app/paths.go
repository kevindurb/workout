package app

import (
	"fmt"
)

type Paths struct {
	name string
}

func (p *Paths) List() string {
	return fmt.Sprintf("/%s", p.name)
}

func (p *Paths) Show(id int64) string {
	return fmt.Sprintf("/%s/%d", p.name, id)
}

func (p *Paths) New() string {
	return fmt.Sprintf("/%s/new", p.name)
}

func (p *Paths) Edit(id int64) string {
	return fmt.Sprintf("/%s/%d/edit", p.name, id)
}
