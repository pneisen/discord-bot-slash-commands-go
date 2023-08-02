package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	BotToken string `split_words:"true" required:"true"`
	GuildID  string `split_words:"true" required:"true"`
}

type discordHandler struct {
	config Config
}

func main() {
	// Config
	var config Config
	err := envconfig.Process("bot", &config)
	if err != nil {
		fmt.Printf("Error with configuration: %s\n", err.Error())
		os.Exit(1)
	}

	// Create discord session
	discord, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		fmt.Printf("Error creating Discord session: %s\n", err.Error())
		os.Exit(1)
	}

	// Register callbacks
	dh := &discordHandler{}
	discord.AddHandler(dh.ready)
	discord.AddHandler(dh.command)

	// Open the websocket
	err = discord.Open()
	if err != nil {
		fmt.Printf("Error opening Discord connection: %s\n", err.Error())
		os.Exit(1)
	}

	// Register command
	command := &discordgo.ApplicationCommand{
		Name:        "test",
		Description: "A test command with subcommand-groups and subcommands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "test-a",
				Description: "Test-a sub-command group",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "test-a-a",
						Description: "Test-a-a sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "test-a-b",
						Description: "Test-a-b sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
			},
			{
				Name:        "test-b",
				Description: "Test-b sub-command group",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "test-b-a",
						Description: "Test-b-a sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "test-b-b",
						Description: "Test-b-b sub-command",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
			},
			{
				Name:        "test-c",
				Description: "Test-c sub-command",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}

	cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.GuildID, command)
	if err != nil {
		fmt.Printf("Error adding command: %s\n", err.Error())
	}

	// Block until we get ctrl-c
	fmt.Println("Bot running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	fmt.Println("Bot exiting")
	err = discord.ApplicationCommandDelete(discord.State.User.ID, config.GuildID, cmd.ID)
	if err != nil {
		fmt.Printf("Error removing command: %s\n", err.Error())
	}
	discord.Close()
}

func (dh *discordHandler) ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateListeningStatus("/test")
}

func (dh *discordHandler) command(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {

	case "test":
		options := i.ApplicationCommandData().Options
		response := ""

		switch options[0].Name {

		case "test-c":
			response = "Test C Command"

		case "test-a":
			options := options[0].Options
			switch options[0].Name {

			case "test-a-a":
				response = "Test A A Command"

			case "test-a-b":
				response = "Test A B Command"

			default:
				response = "Error!"
			}

		case "test-b":
			options := options[0].Options
			switch options[0].Name {

			case "test-b-a":
				response = "Test B A Command"

			case "test-b-b":
				response = "Test B B Command"

			default:
				response = "Error!"
			}

		default:
			response = "Error!"
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	}
}
