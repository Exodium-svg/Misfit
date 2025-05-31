package Game

import (
	"MisFitsDiscord/Utils"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"sync"
	"time"
)

const toaistFilename string = "taoists.json"

var taoistMap map[string]*Taoist
var taoistLock sync.Mutex

type Taoist struct {
	ZodiacSign int        `json:"zodiac_sign"`
	UserID     string     `json:"user_id"`
	Ascension  *Ascension `json:"ascension:omitempty"`
	Level      int        `json:"level"`
	CurrentXp  int        `json:"current_xp"`
	LastSpoke  time.Time  `json:"last_spoke"`
}

func SaveTaoists() {
	taoistLock.Lock()
	defer taoistLock.Unlock()

	jsonData, err := json.Marshal(taoistMap)

	if nil != err {
		fmt.Printf("Failed to save taoist data error: %s\n", err.Error())
		return
	}

	err = os.WriteFile(toaistFilename, jsonData, 0644)

	if nil != err {
		fmt.Printf("Failed to save taoist data error: %s\n", err.Error())
		return
	}
}

func LoadTaoists() {
	taoistLock.Lock()
	defer taoistLock.Unlock()

	if !Utils.FileExists(toaistFilename) {
		taoistMap = make(map[string]*Taoist)
		return
	}

	LoadAscensions()

	jsonData, err := Utils.ReadFile(toaistFilename)

	if nil != err {
		fmt.Printf("Failed to load taoist data error: %s\n", err.Error())
		taoistMap = make(map[string]*Taoist)
		return
	}

	if len(jsonData) == 0 {
		taoistMap = make(map[string]*Taoist)
		return
	}

	err = json.Unmarshal(jsonData, &taoistMap)

	if nil != err {
		fmt.Printf("Failed to load taoist data error: %s\n", err.Error())
		taoistMap = make(map[string]*Taoist)
	}
}

func getTaoist(userId string) *Taoist {
	taoist, exists := taoistMap[userId]

	if false == exists {
		taoist = &Taoist{
			ZodiacSign: int(time.Now().Month()),
			UserID:     userId,
			Level:      1,
			Ascension:  nil,
			CurrentXp:  0,
			LastSpoke:  time.Now().Add(-100 * time.Hour),
		}

		taoistMap[userId] = taoist
	}

	return taoist
}

func GetTaoists() []Taoist {
	taoistLock.Lock()
	defer taoistLock.Unlock()

	taoists := taoistMap

	taoistSlice := make([]Taoist, len(taoists))

	i := 0
	for _, taoist := range taoists {
		taoistSlice[i] = *taoist
		i++
	}

	return taoistSlice
}

func GetTaoist(userID string) *Taoist {
	taoistLock.Lock()
	defer taoistLock.Unlock()

	taoist, exist := taoistMap[userID]
	if !exist {
		return nil
	}

	return taoist
}

func taoistAscensionCheck(env *Utils.Env, taoist *Taoist, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	var roleID string

	if nil == taoist.Ascension {
		roleID = "NULL"
	} else {
		roleID = taoist.Ascension.RoleId
	}

	newAscension, exists := NewAscension(taoist.Level, roleID)

	if false == exists {
		return false
	}

	// validation on discord

	if "not_loaded" == newAscension.Title {
		roles, err := session.GuildRoles(message.GuildID)

		if nil != err {
			fmt.Printf("Failed to get guild roles: %s\n", err.Error())
			return false
		}

		for _, role := range roles {
			if role.ID == newAscension.RoleId {
				newAscension.Title = role.Name
			}
		}
	}

	taoist.Ascension = newAscension

	if roleID != "NULL" {
		err := session.GuildMemberRoleRemove(message.GuildID, message.Author.ID, roleID)

		if nil != err {
			fmt.Printf("Failed to remove role error: %s\n", err.Error())
			return false
		}
	}

	err := session.GuildMemberRoleAdd(message.GuildID, message.Author.ID, newAscension.RoleId)

	if nil != err {
		fmt.Printf("Failed to add role error: %s\n", err.Error())
		return false
	}

	return true
}

func RequiredXpForLevelUp(taoist *Taoist, levelUpCurve int) int {

	if taoist.Level == 1 {
		return 10
	}

	pow := 1
	for i := 0; i < levelUpCurve; i++ {
		pow *= taoist.Level
	}

	return pow
}

func OnConversation(env *Utils.Env, session *discordgo.Session, event *discordgo.MessageCreate) {
	taoistLock.Lock()
	defer taoistLock.Unlock()

	minutesDelay := Utils.EnvGet(env, "immortal.conversation.time", 1)
	currentTime := time.Now()

	taoist := getTaoist(event.Author.ID)

	if taoist.LastSpoke.After(currentTime) {
		return
	}

	taoist.LastSpoke = currentTime.Add(time.Minute * time.Duration(minutesDelay))

	levelUpCurve := Utils.EnvGet(env, "immortal.level_up.curve", 100)

	xpForLevelUp := RequiredXpForLevelUp(taoist, levelUpCurve)

	xpModifier := 1.0

	if taoist.ZodiacSign == GetZodiacSign() {
		xpModifier = Utils.EnvGet(env, "immortal.zodiac.multiplier", 1.1)
	}

	experience := 1 + rand.Intn(3)

	taoist.CurrentXp += int(float64(experience) * xpModifier)

	if taoist.CurrentXp >= xpForLevelUp {
		taoist.CurrentXp = 0
		taoist.Level++

		_, err := session.ChannelMessageSendComplex(event.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf(Utils.GetString("taoist.level_up", "In stillness and flow, your essence refines. You have reached Level %d."), taoist.Level),
			Reference: &discordgo.MessageReference{
				MessageID: event.Message.ID,
				ChannelID: event.ChannelID,
				GuildID:   event.GuildID,
			},
		})

		if nil != err {
			fmt.Printf("Failed to send message %s\n", err.Error())
		}

		newAscension := taoistAscensionCheck(env, taoist, session, event)

		if true == newAscension {
			_, err = session.ChannelMessageSendComplex(event.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf(Utils.GetString("taoist.ascensions_up", "Through silence and stars, your spirit dissolves into the eternal flow. You ascend as: %s.."), taoist.Ascension.Title),
				Reference: &discordgo.MessageReference{
					MessageID: event.Message.ID,
					ChannelID: event.ChannelID,
				},
			})

			if nil != err {
				fmt.Printf("Failed to send message %s\n", err.Error())
			}

		}
	}

}
