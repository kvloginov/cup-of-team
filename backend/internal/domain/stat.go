package domain

type Stats struct {
	HP            int     `json:"hp"`
	Damage        int     `json:"damage"`
	CritChance    float32 `json:"crit_chance"`
	EvasionChance float32 `json:"evasion_chance"`
}
