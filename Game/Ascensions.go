package Game

import (
	"MisFitsDiscord/Utils"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"sort"
	"sync"
)

const filePath string = "ascensions.json"

var ascensions []Ascension = nil
var ascensionLock sync.Mutex

type Ascension struct {
	Title         string `json:"title"`
	RequiredLevel int    `json:"required_level"`
	RoleId        string `json:"role_id"`
}

func LoadAscensions() {
	ascensionLock.Lock()
	defer ascensionLock.Unlock()

	jsonData, err := Utils.ReadFile(filePath)

	if nil != err {
		ascensions = make([]Ascension, 0)
		return
	}

	err = json.Unmarshal(jsonData, &ascensions)

	sort.Slice(ascensions, func(i, j int) bool {
		return ascensions[i].RequiredLevel > ascensions[j].RequiredLevel
	})
}

func NewAscension(requiredLevel int, currentRoleId string) (*Ascension, bool) {
	ascensionLock.Lock()
	defer ascensionLock.Unlock()

	for i, ascension := range ascensions {
		if ascension.RequiredLevel <= requiredLevel && ascension.RoleId != currentRoleId {

			return &ascensions[i], true
		}
	}

	return nil, false
}

type AscensionSet struct {
	Title         string `json:"title"`
	RequiredLevel int    `json:"requiredLevel"`
	RoleId        string `json:"roleId"`
}

func GetAscensions() []Ascension {
	return ascensions
}

func SetAscensions(sets []AscensionSet) {
	ascensionLock.Lock()

	ascensions = make([]Ascension, len(sets))
	for i, ascension := range sets {
		ascensions[i] = Ascension{
			Title:         ascension.Title,
			RequiredLevel: ascension.RequiredLevel,
			RoleId:        ascension.RoleId,
		}
	}

	sort.Slice(ascensions, func(i, j int) bool {
		return ascensions[i].RequiredLevel > ascensions[j].RequiredLevel
	})

	ascensionLock.Unlock()
	SaveAscensions()
}

func AddAscension(session *discordgo.Session, env *Utils.Env, requiredLevel int, currentRoleId string) {
	ascensionLock.Lock()
	defer ascensionLock.Unlock()

	guildId := Utils.EnvGet(env, "discord.guild.id", "1370077540160372847")

	roles, err := session.GuildRoles(guildId)

	if nil != err {
		fmt.Printf("Failed to fetch roles, error: %s\n", err.Error())
		return
	}

	for _, role := range roles {
		if role.ID == currentRoleId {
			ascensions = append(ascensions, Ascension{
				Title:         role.Name,
				RequiredLevel: requiredLevel,
				RoleId:        role.ID,
			})
		}
	}
}

func RemoveAscension(roleId string) {
	ascensionLock.Lock()
	defer ascensionLock.Unlock()

	for i, ascension := range ascensions {
		if ascension.RoleId == roleId {
			ascensions = append(ascensions[:i], ascensions[i+1:]...)
		}
	}
}

func SaveAscensions() {
	ascensionLock.Lock()
	defer ascensionLock.Unlock()

	jsonData, err := json.Marshal(ascensions)

	if nil != err {
		fmt.Printf("Failed to save ascensions: %s\n", err.Error())
		return
	}

	err = os.WriteFile(filePath, jsonData, 0644)

	if nil != err {
		fmt.Printf("Failed to save ascensions: %s\n", err.Error())
		return
	}
}
