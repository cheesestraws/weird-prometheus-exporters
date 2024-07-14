package wiltshirebins

type Collections struct {
	HouseholdWaste int `prometheus:"household_waste"`
	Recycling      int `prometheus:"recycling"`

	Errors map[string]int `prometheus_map:"errors" prometheus_map_key:"error"`
}

func (c *Collections) errorFill(errorType string) {
	if c.Errors == nil {
		c.Errors = make(map[string]int)
	}

	c.Errors[errorType]++
}

type Calendar [31]Collections

func (c *Calendar) errorFill(errorType string) {
	for i := range *c {
		coll := &((*c)[i])
		coll.errorFill(errorType)
	}
}

// Checks to see whether the data looks approximately sane
func (c *Calendar) sanityCheck() {
	var totals Collections

	for _, coll := range *c {
		totals.HouseholdWaste = totals.HouseholdWaste + coll.HouseholdWaste
		totals.Recycling = totals.Recycling + coll.Recycling
	}

	hwRatio := float64(totals.HouseholdWaste) / float64(31)
	reRatio := float64(totals.HouseholdWaste) / float64(31)

	if hwRatio < 0.04 {
		c.errorFill("dubious_data/household_too_low")
	}

	if hwRatio > 0.21 {
		c.errorFill("dubious_data/household_too_high")
	}

	if reRatio < 0.04 {
		c.errorFill("dubious_data/recycle_too_low")
	}

	if reRatio > 0.21 {
		c.errorFill("dubious_data/recycle_too_high")
	}
}
