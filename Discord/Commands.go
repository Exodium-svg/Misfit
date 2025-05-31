package Discord

import (
	"MisFitsDiscord/Game"
	"MisFitsDiscord/Utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sort"
	"time"
)

func interactionRespond(content string, session *discordgo.Session, interaction *discordgo.Interaction) {
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if nil != err {
		fmt.Printf("Error responding to interaction: %s\n", err.Error())
	}
	return
}

func OnRequestCurrentZodiacSign(session *discordgo.Session, event *discordgo.InteractionCreate) {
	currentZodiac := Game.GetZodiacName(Game.GetZodiacSign())
	err := session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(Utils.GetString("command.taoist.current_zodiac", "This cycle's celestial blessing falls upon: %s"), currentZodiac),
		},
	})

	if nil != err {
		fmt.Printf("Error responding to interaction: %s\n", err.Error())
		return
	}
}

func OnRequestTopTen(session *discordgo.Session, event *discordgo.InteractionCreate) {
	taoists := Game.GetTaoists()

	sort.Slice(taoists, func(i, j int) bool {
		return taoists[i].Level > taoists[j].Level
	})

	var count int

	if len(taoists) < 10 {
		count = len(taoists)
	} else {
		count = 10
	}

	if len(taoists) == 0 {
		interactionRespond(Utils.GetString("command.taoist.top.empty", "No cultivators have emerged yet—be the first to ascend!"), session, event.Interaction)
		return
	}

	fields := make([]*discordgo.MessageEmbedField, 0, count)
	for i := 0; i < count; i++ {
		taoist := taoists[i]

		var title string

		if taoist.Ascension != nil {
			title = fmt.Sprintf("%s (%d)\n", taoist.Ascension.Title, taoist.Level)
		} else {
			title = fmt.Sprintf("%s (%d)\n", "Unknown", taoist.Level)
		}

		user, err := session.User(taoist.UserID)

		if nil != err {
			fmt.Printf("Error fetching user %s\n", err.Error())
			continue
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  user.DisplayName(),
			Value: title,
		})
	}

	embed := discordgo.MessageEmbed{
		Title:     fmt.Sprintf(Utils.GetString("command.taoist.top", "Top %d Masters of the Realm"), count),
		Color:     0x00ffcc,
		Fields:    fields,
		Footer:    &discordgo.MessageEmbedFooter{},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err := session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				&embed,
			},
		},
	})

	if nil != err {
		fmt.Printf("Failed to send result of interaction error: %s\n", err.Error())
		return
	}
}

func OnRequestTaoistInfo(session *discordgo.Session, event *discordgo.InteractionCreate) {
	taoist := Game.GetTaoist(event.Member.User.ID)

	if nil == taoist {
		interactionRespond("You do not yet walk the Way.", session, event.Interaction)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🌿 Taoist Insight",
		Description: "Your place on the path is known.",
		Color:       0x88cc88, // soft green
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Level",
				Value:  fmt.Sprintf("%d", taoist.Level),
				Inline: true,
			},
			{
				Name:   "Experience",
				Value:  fmt.Sprintf("%d / %d XP", taoist.CurrentXp, Game.RequiredXpForLevelUp(taoist, Utils.EnvGet(env, "immortal.level_up.curve", 100))),
				Inline: true,
			},
			{
				Name:   "Zodiac Sign",
				Value:  Game.GetZodiacName(taoist.Level),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err := session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	if nil != err {
		fmt.Printf("Failed to send result of interaction error: %s\n", err.Error())
		return
	}
}
