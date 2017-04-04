package dbindex

type Params struct {
	value []interface{}
}

func NewParam(v interface{}) *Params {
	ret := &Params{}
	ret.value = append(ret.value, v)
	return ret
}

func (p *Params) Add(v interface{}) error {
	switch x := v.(type) {
	case string:
		if x != "" {
			p.value = append(p.value, x)
		}
	case []string:
		if x != nil {
			p.value = append(p.value, x)
		}
	case int:
		p.value = append(p.value, x)
	}
	return nil
}
