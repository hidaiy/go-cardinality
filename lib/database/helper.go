package database

// Params
type Params struct {
	Values []interface{}
}

// NewParams returns Params as pointer
func NewParams(v interface{}) *Params {
	ret := &Params{}
	ret.Values = append(ret.Values, v)
	return ret
}

func (p *Params) Add(v interface{}) error {
	switch x := v.(type) {
	case string:
		if x != "" {
			p.Values = append(p.Values, x)
		}
	case []string:
		if x != nil {
			p.Values = append(p.Values, x)
		}
	case int:
		p.Values = append(p.Values, x)
	}
	return nil
}
