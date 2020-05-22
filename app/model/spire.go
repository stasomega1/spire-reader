package model

//json
type JsonSpireRun struct {
	PathPerFloor []*string `json:"path_per_floor"`
	//M - Monster - Монстр
	//T - Treasure - Сундук
	//R - Rest - Отдых
	//E - Mini boss - Мини босс
	//B - Boss - Босс
	//? - Random encounter - Случайное событие
	//$ - Shop - Магазин
	//null - Between floors - Между этажами
	PathPerFloorNormal []string
	//PathPerFloor w/o null (null -> "N")
	CharacterChosen string `json:"character_chosen"`
	//Ironclad - Латоносец
	//Silent
	//Defect
	//Watcher
	ItemsPurchased      []string `json:"items_purchased"`
	GoldPerFloor        []int    `json:"gold_per_floor"`
	FloorReached        int      `json:"floor_reached"`
	CampfireRested      int      `json:"campfire_rested"`
	Playtime            int      `json:"playtime"`
	CurrentHpPerFloor   []int    `json:"current_hp_per_floor"`
	ItemsPurged         []string `json:"items_purged"`
	Gold                int      `json:"gold"`
	Score               int      `json:"score"`
	PlayID              string   `json:"play_id"`
	LocalTime           string   `json:"local_time"`
	IsProd              bool     `json:"is_prod"`
	IsDaily             bool     `json:"is_daily"`
	ChoseSeed           bool     `json:"chose_seed"`
	IsAscensionMode     bool     `json:"is_ascension_mode"`
	CampfireUpgraded    int      `json:"campfire_upgraded"`
	Timestamp           int      `json:"timestamp"`
	PathTaken           []string `json:"path_taken"`
	BuildVersion        string   `json:"build_version"`
	SeedSourceTimestamp int64    `json:"seed_source_timestamp"`
	PurchasedPurges     int      `json:"purchased_purges"`
	Victory             bool     `json:"victory"`
	MasterDeck          []string `json:"master_deck"`
	MaxHpPerFloor       []int    `json:"max_hp_per_floor"`
	Relics              []string `json:"relics"`
	CardChoices         []struct {
		NotPicked []string `json:"not_picked"`
		Picked    string   `json:"picked"`
		Floor     int      `json:"floor"`
	} `json:"card_choices"`
	PlayerExperience  int   `json:"player_experience"`
	PotionsFloorUsage []int `json:"potions_floor_usage"`
	DamageTaken       []struct {
		Damage  int    `json:"damage"`
		Enemies string `json:"enemies"`
		Floor   int    `json:"floor"`
		Turns   int    `json:"turns"`
	} `json:"damage_taken"`
	EventChoices []struct {
		EventName    string `json:"event_name"`
		PlayerChoice string `json:"player_choice"`
		Floor        int    `json:"floor"`
		DamageTaken  int    `json:"damage_taken"`
	} `json:"event_choices"`
	BossRelics []struct {
		NotPicked []string `json:"not_picked"`
		Picked    string   `json:"picked"`
	} `json:"boss_relics"`
	PotionsFloorSpawned []int  `json:"potions_floor_spawned"`
	SeedPlayed          string `json:"seed_played"`
	KilledBy            string `json:"killed_by"`
	AscensionLevel      int    `json:"ascension_level"`
	IsTrial             bool   `json:"is_trial"`
}

type SpireRun struct {
	CharacterChosen *JsonNullString   `json:"characterChosen"`
	ItemsPurchased  []*JsonNullString `json:"itemsPurchased"`
	PathPerFloor    []*JsonNullString `json:"pathPerFloor"`
	FloorReached    *JsonNullInt32    `json:"floorReached"`
	Playtime        *JsonNullInt32    `json:"playtime"`
	Victory         *JsonNullBool     `json:"victory"`
}

type SpireRuns struct {
}
