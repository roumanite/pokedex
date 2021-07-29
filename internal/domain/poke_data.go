package domain

type EncounterDetail struct {
	Chance int
	MaxLevel int
	MethodName string
}

type VersionDetail struct {
	EncounterDetails []EncounterDetail
	MaxChance int
	VersionName string
}

type LocationAreaEncounter struct {
	Name string
	VersionDetails []VersionDetail
}

type StatItem struct {
	BaseStat string
	Effort int
	StatName string
}

type PokeData struct {
	Id int
	Name string
	Types []string
	Stats []StatItem
}