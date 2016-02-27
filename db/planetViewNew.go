package db

type PlanetView struct {
	*MiniPlanet
}

func NewPlanetView() *PlanetView {
	return &PlanetView{
		&MiniPlanet{},
	}
}

func (item *PlanetView) SQLTable() string {
	return "planetviews"
}

func (group *PlanetViewGroup) SQLTable() string {
	return "planetviews"
}

func (group *PlanetViewGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"turn",
		"name",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
	}
}

func (group *PlanetViewGroup) UpdateCols() []string {
	return []string{
		"turn",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
	}
}

func (group *PlanetViewGroup) PKCols() []string {
	return []string{"gid", "fid", "locx", "locy"}
}

func (group *PlanetViewGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"turn",
		"name",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
	}
}

func (group *PlanetViewGroup) InsertScanCols() []string {
	return nil
}
