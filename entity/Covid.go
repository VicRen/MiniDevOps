package entity

type Duration int

const (
	Daily Duration = iota
	Weekly
	Monthly
	Custom
)

type Covid struct {
	Confirmed   int
	Deaths      int
	Recovered   int
	Active      int
	Duration    Duration
	CustomStart int
	CustomEnd   int
}

func (c *Covid) Days() {

}
