package mocks

import "time"

type Clock struct {
	t time.Time
}

func (c *Clock) Now() time.Time {
	return c.t
}

func (c *Clock) SetTime(t time.Time) {
	c.t = t
}

func (c *Clock) Move(delta time.Duration) {
	c.t.Add(delta)
}
