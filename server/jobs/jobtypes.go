package jobs

type ESPJobPeaks struct {
	Ecto int
	Spark int
	Pulse int
}

type Job struct {
	Type string
	Peak ESPJobPeaks
}

var (
	Mercenary = Job{
		Type: "mercenary",
		Peak: ESPJobPeaks{
			Ecto: 450,
			Spark: 850,
			Pulse: 750,
		},
	}
	Farmer = Job{
		Type: "farmer",
		Peak: ESPJobPeaks{
			Ecto: 650,
			Spark: 750,
			Pulse: 450,
		},
	}
	Explorer = Job{
		Type: "explorer",
		Peak: ESPJobPeaks{
			Ecto: 850,	// very eager to increase "connectedness"
			Spark: 750, // decent energy for travelling
			Pulse: 450, // fairly lightweight with not much "stability"
		},
	}
	Warden = Job{
		Type: "warden",
		Peak: ESPJobPeaks{
			Ecto: 650,
			Spark: 750,
			Pulse: 550,
		},
	}
)