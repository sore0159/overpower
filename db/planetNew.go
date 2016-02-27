package db

type Planet struct {
	*MiniPlanet
}

func NewPlanet() *Planet {
	return &Planet{
		MiniPlanet: &MiniPlanet{},
	}
}

func (p *Planet) SQLTable() string {
	return "planets"
}

func (group *PlanetGroup) SQLTable() string {
	return "planets"
}

func (group *PlanetGroup) SelectCols() []string {
	return []string{
		"gid",
		"locx",
		"locy",
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

func (group *PlanetGroup) UpdateCols() []string {
	return []string{
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

func (group *PlanetGroup) PKCols() []string {
	return []string{"gid", "locx", "locy"}
}

func (group *PlanetGroup) InsertCols() []string {
	return []string{
		"gid",
		"locx",
		"locy",
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

func (group *PlanetGroup) InsertScanCols() []string {
	return nil
}
