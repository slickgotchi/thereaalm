package jobs

type ESPJobPeaks struct {
	Ecto  int
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
			Ecto:  540, // Moderate luck for survival
			Spark: 940, // High strength for combat
			Pulse: 620, // Moderate endurance for battles
		},
	}
	Warden = Job{
		Type: "warden",
		Peak: ESPJobPeaks{
			Ecto:  520, // Moderate focus for vigilance
			Spark: 830, // High strength to defend
			Pulse: 750, // High endurance for long watches
		},
	}
	Thief = Job{
		Type: "thief",
		Peak: ESPJobPeaks{
			Ecto:  870, // High luck to avoid capture
			Spark: 670, // Moderate strength for stealth
			Pulse: 560, // Low endurance for quick getaways
		},
	}
	Beastmaster = Job{
		Type: "beastmaster",
		Peak: ESPJobPeaks{
			Ecto:  560, // Moderate connection with animals
			Spark: 890, // High strength to handle beasts
			Pulse: 650, // Moderate endurance for trekking
		},
	}
	Medic = Job{
		Type: "medic",
		Peak: ESPJobPeaks{
			Ecto:  790, // High focus for healing
			Spark: 700, // Moderate strength for supplies
			Pulse: 610, // Moderate endurance for shifts
		},
	}
	Merchant = Job{
		Type: "merchant",
		Peak: ESPJobPeaks{
			Ecto:  880, // High connection for trade
			Spark: 740, // Moderate drive for negotiation
			Pulse: 480, // Low endurance for travel
		},
	}
	Crafter = Job{
		Type: "crafter",
		Peak: ESPJobPeaks{
			Ecto:  740, // High focus for crafting
			Spark: 720, // High strength for crafting
			Pulse: 640, // Moderate endurance
		},
	}
	Farmer = Job{
		Type: "farmer",
		Peak: ESPJobPeaks{
			Ecto:  500, // Low luck for harvests
			Spark: 660, // Moderate strength for labor
			Pulse: 940, // High endurance for farming
		},
	}
	MinerJack = Job{
		Type: "minerjack",
		Peak: ESPJobPeaks{
			Ecto:  480, // Low luck for ore
			Spark: 770, // High strength for mining
			Pulse: 850, // High endurance for slow mining
		},
	}
	Builder = Job{
		Type: "builder",
		Peak: ESPJobPeaks{
			Ecto:  450, // Low focus for building
			Spark: 750, // High strength for construction
			Pulse: 900, // High endurance for projects
		},
	}
	Alchemist = Job{
		Type: "alchemist",
		Peak: ESPJobPeaks{
			Ecto:  750, // High focus for experiments
			Spark: 590, // Moderate strength for gathering
			Pulse: 760, // Moderate endurance
		},
	}
	Explorer = Job{
		Type: "explorer",
		Peak: ESPJobPeaks{
			Ecto:  760, // High luck for discoveries
			Spark: 650, // Moderate strength for travel
			Pulse: 690, // High endurance for exploration
		},
	}
	Scholar = Job{
		Type: "scholar",
		Peak: ESPJobPeaks{
			Ecto:  950, // High focus for study
			Spark: 490, // Low strength
			Pulse: 660, // Moderate endurance for study
		},
	}
	Engineer = Job{
		Type: "engineer",
		Peak: ESPJobPeaks{
			Ecto:  810, // High focus for engineering
			Spark: 620, // Moderate strength for building machines
			Pulse: 670, // Moderate endurance
		},
	}
	Diplomat = Job{
		Type: "diplomat",
		Peak: ESPJobPeaks{
			Ecto:  900, // High connection for diplomacy
			Spark: 480, // Low drive for negotiation
			Pulse: 720, // Moderate endurance
		},
	}
)