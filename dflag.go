package dflag

type Option struct {
	Name   string
	Value  string
	IsFlag bool
	Raw    []string
}

type OptionSet struct {
	Parsed  bool
	current string
	remain  []string
}

func (s *OptionSet) next() bool {
	if len(s.remain) == 0 {
		return false
	}
	s.current = s.remain[0]
	s.remain = s.remain[1:]
	return true
}

func NewOptionSet(args []string) *OptionSet {
	return &OptionSet{
		Parsed:  false,
		current: "",
		remain:  args,
	}
}

func (s *OptionSet) Parse() *Option {
	if !s.next() {
		s.Parsed = true
		return nil
	}
	if s.current[0] == '-' {
		if s.current != "-" && s.current != "--" {
			return s.ParseFlag()
		}
	}
	return s.ParseValue()
}

func (s *OptionSet) ParseValue() *Option {
	return &Option{
		Name:   "",
		Value:  s.current,
		IsFlag: false,
		Raw:    []string{s.current},
	}
}

func (s *OptionSet) ParseFlag() *Option {
	c := s.current
	if c[1] == '-' {
		for i := 1; i < len(c); i++ {
			if c[i] == '=' {
				return &Option{
					Name:   c[:i],
					Value:  (c[i+1:]),
					IsFlag: true,
					Raw:    []string{c},
				}
			}
		}
	} else if len(c) > 2 {
		return &Option{
			Name:   c[:2],
			Value:  c[2:],
			IsFlag: true,
			Raw:    []string{c},
		}
	}
	if o := s.Parse(); o != nil {
		if o.IsFlag {
			s.remain = append([]string{s.current}, s.remain...)
			s.current = c
		} else {
			return &Option{
				Name:   c,
				Value:  o.Value,
				IsFlag: true,
				Raw:    append([]string{c}, o.Raw...),
			}
		}
	}
	return &Option{
		Name:   c,
		Value:  "",
		IsFlag: true,
		Raw:    []string{c},
	}
}
