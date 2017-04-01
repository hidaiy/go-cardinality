package consoleTable

type Footer struct {
	Base
}

func (c *Footer) Write(width width) (int, error) {
	c.width = width
	return c.writeSeparatorLine("└", "┴", "┘")
}
