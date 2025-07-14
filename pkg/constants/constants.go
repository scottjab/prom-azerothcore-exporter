package constants

import "fmt"

// WoW race to faction mapping
var RaceToFaction = map[int]string{
	1:  "Alliance", // Human
	2:  "Horde",    // Orc
	3:  "Alliance", // Dwarf
	4:  "Alliance", // Night Elf
	5:  "Horde",    // Undead
	6:  "Horde",    // Tauren
	7:  "Alliance", // Gnome
	8:  "Horde",    // Troll
	9:  "Horde",    // Goblin
	10: "Horde",    // Blood Elf
	11: "Alliance", // Draenei
	22: "Horde",    // Worgen
	24: "Alliance", // Pandaren (Neutral)
	25: "Alliance", // Pandaren (Alliance)
	26: "Horde",    // Pandaren (Horde)
}

// WoW race names
var RaceNames = map[int]string{
	1:  "Human",
	2:  "Orc",
	3:  "Dwarf",
	4:  "Night Elf",
	5:  "Undead",
	6:  "Tauren",
	7:  "Gnome",
	8:  "Troll",
	9:  "Goblin",
	10: "Blood Elf",
	11: "Draenei",
	22: "Worgen",
	24: "Pandaren (Neutral)",
	25: "Pandaren (Alliance)",
	26: "Pandaren (Horde)",
}

// WoW class names
var ClassNames = map[int]string{
	1:  "Warrior",
	2:  "Paladin",
	3:  "Hunter",
	4:  "Rogue",
	5:  "Priest",
	6:  "Death Knight",
	7:  "Shaman",
	8:  "Mage",
	9:  "Warlock",
	10: "Monk",
	11: "Druid",
	12: "Demon Hunter",
}

// Helper functions for readable names
func GetDifficultyName(difficulty int) string {
	difficultyNames := map[int]string{
		0: "Normal",
		1: "Heroic",
		2: "10_Player",
		3: "25_Player",
		4: "10_Player_Heroic",
		5: "25_Player_Heroic",
	}
	if name, exists := difficultyNames[difficulty]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", difficulty)
}

func GetLFGStateName(state int) string {
	stateNames := map[int]string{
		0: "None",
		1: "RoleCheck",
		2: "Queued",
		3: "Proposal",
		4: "Boot",
		5: "Dungeon",
		6: "FinishedDungeon",
	}
	if name, exists := stateNames[state]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", state)
}

func GetIPActionTypeName(actionType int) string {
	typeNames := map[int]string{
		0: "Login",
		1: "Failed_Login",
		2: "Logout",
		3: "Character_Create",
		4: "Character_Delete",
		5: "Character_Login",
		6: "Character_Logout",
		7: "Password_Change",
		8: "Account_Create",
		9: "Account_Delete",
	}
	if name, exists := typeNames[actionType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", actionType)
}

func GetLagTypeName(lagType int) string {
	lagTypeNames := map[int]string{
		0: "World",
		1: "Instance",
		2: "Battleground",
		3: "Arena",
		4: "Raid",
	}
	if name, exists := lagTypeNames[lagType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", lagType)
}

func GetBattlegroundTypeName(bgType int) string {
	bgTypeNames := map[int]string{
		1:  "Alterac Valley",
		2:  "Warsong Gulch",
		3:  "Arathi Basin",
		4:  "Eye of the Storm",
		5:  "Strand of the Ancients",
		6:  "Isle of Conquest",
		7:  "Twin Peaks",
		8:  "Battle for Gilneas",
		9:  "Temple of Kotmogu",
		10: "Silvershard Mines",
		11: "Deepwind Gorge",
	}
	if name, exists := bgTypeNames[bgType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", bgType)
}

func GetDesertionTypeName(desertionType int) string {
	desertionTypeNames := map[int]string{
		0: "Leave",
		1: "Offline",
		2: "Desert",
		3: "Finish",
	}
	if name, exists := desertionTypeNames[desertionType]; exists {
		return name
	}
	return fmt.Sprintf("Unknown_%d", desertionType)
}
