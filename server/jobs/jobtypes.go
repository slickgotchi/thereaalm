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
	Warden = Job{
		Type: "warden",
		Peak: ESPJobPeaks{
			Ecto: 650,
			Spark: 750,
			Pulse: 550,
		},
	}
)