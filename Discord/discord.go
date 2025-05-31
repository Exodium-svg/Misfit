package Discord

import (
	"MisFitsDiscord/Game"
	"MisFitsDiscord/Utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

var env *Utils.Env

func Start(_env *Utils.Env) *discordgo.Session {
	env = _env

	token := Utils.EnvGet(env, "discord.token", "not_found")

	if "not_found" == token {
		fmt.Print("Token wasn't found")
		return nil
	}

	session, err := discordgo.New("Bot " + token)

	if nil != err {
		fmt.Printf("error creating Discord session,%s\n", err.Error())
		return nil
	}

	session.AddHandler(onReady)
	session.AddHandler(onInteraction)
	session.AddHandler(onMessage)

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	err = session.Open()

	if nil != err {
		fmt.Printf("error opening connection,%s\n", err.Error())
		return nil
	}

	registerCommands(session)

	return session
}

func registerCommands(session *discordgo.Session) {
	cmd := []discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Pong!",
		},
		{
			Name:        "top",
			Description: "Unveils the top 10 legendary cultivators.",
		},
		{
			Name:        "zodiac",
			Description: "Reveals the sign in harmony with the Way—blessed by the heavens in this moment.",
		},
		{
			Name:        "taoist",
			Description: "Reveals your current stage on the path of cultivation and alignment with the celestial flow.",
		},
	}

	for _, command := range cmd {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, Utils.EnvGet(env, "discord.guild.id", "1370077540160372847"), &command)
		if nil != err {
			fmt.Printf("error creating ApplicationCommand %s\n", err.Error())
			os.Exit(-1)
		}
	}

}

func onReady(session *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Ready on user: %s\n", session.State.User.Username)
}

func onInteraction(session *discordgo.Session, event *discordgo.InteractionCreate) {
	if discordgo.InteractionApplicationCommand == event.Type {
		onCommand(session, event)
		return
	}
}

func onMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.Bot || event.Author.ID == session.State.User.ID || len(event.Attachments) > 0 {
		return
	}

	Game.OnConversation(env, session, event)
}

func onCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	data := event.ApplicationCommandData()
	var err error = nil

	switch data.Name {
	case "ping":
		err = session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		})
	case "top":
		OnRequestTopTen(session, event)
		break
	case "zodiac":
		OnRequestCurrentZodiacSign(session, event)
		break
	case "taoist":
		OnRequestTaoistInfo(session, event)
		break
	}

	if nil != err {
		fmt.Printf("error interaction handling %s\n", err.Error())
	}
}
