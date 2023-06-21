package ingest

type ScalarEvents struct {
	// These three slices are parallel in order to save memory.
	// i.e. the i-th element of each slice corresponds to the same event.
	Timestamp []float64
	Step      []int64
	Value     []float32
}

type Run struct {
	Name    string
	Scalars map[string]*ScalarEvents
}

func NewRun(name string) *Run {
	return &Run{Name: name, Scalars: make(map[string]*ScalarEvents)}
}

func (r *Run) AddScalarEvent(tag string, timestamp float64, step int64, value float32) {
	if r.Scalars == nil {
		panic("Run is not initialized")
	}

	scalar, ok := r.Scalars[tag]
	if !ok {
		scalar = &ScalarEvents{}
		r.Scalars[tag] = scalar
	}
	scalar.Timestamp = append(scalar.Timestamp, timestamp)
	scalar.Step = append(scalar.Step, step)
	scalar.Value = append(scalar.Value, value)
}
