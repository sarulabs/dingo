package di

// containerPlacer contains all the functions that are useful
// to put an object to a container.
type containerPlacer struct{}

func (c *containerPlacer) put(ctn *container, name string, dst interface{}) error {
	f := func(ctn Container) (interface{}, error) {
		return dst, nil
	}

	def := ctn.definitions[name]
	def.Build = f

	ctn.definitions[name] = def

	return nil
}
