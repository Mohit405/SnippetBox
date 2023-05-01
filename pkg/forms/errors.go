package forms

type errors map[string][]string

// method to add error messages
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Implement a Get() method to retrieve the first error message for a given
func (e errors) Get(filed string) string {
	es := e[filed]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
