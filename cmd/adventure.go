/*
Copyright © 2023 Kevin.Jayne@iCloud.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/seemywingz/goai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Character struct {
	Name        string
	Description string
	HP          float64
	MP          float64
	Level       int64
	Strength    float64
	Defense     float64
	Dexterity   float64
	Intellect   float64
	Hunger      float64
}

var player Character

var adventureSystemMessage = `
You are the narrator and won't accept any answers that are not relevant to the current story or anything that wasn't mentioned yet.
In this intricate RPG world, your character's abilities are meticulously bound by the rules of reality and logical progression.
Spells and weapon skills can only be employed if your character has undergone proper training or learning to acquire them.
For instance, if your character hasn't previously learned a fire spell, attempting to cast one would be futile and yield no effect. 
Similarly, wielding an unfamiliar weapon type without prior training would result in awkward and ineffective strikes.
You cannot take any action that isn't being asked of you. 
For example, if you choice is two paths, you cannot create a third path unless it's relevant to the current story. 
You must choose from the options presented to you.

You are a young adventurer who is just starting out on your journey. 
You have no money, no weapons, and no armor. 
You are wearing a simple tunic with defense 1 and trousers with defense 1. 
You have a small pouch of coins, but not enough to buy anything useful.

MP represents your magical power and recovers over time.
HP represents your health and can be restored by drinking potions or resting. if your HP reaches 0, you will die.
Strength represents your physical strength and affects your ability to wield weapons and armor it also increases your HP by 2 per level.
Defense represents your ability to defend yourself from attacks and is affected by your armor.
Dexterity represents your ability to dodge attacks and perform acrobatic feats and use ranged weapons.
Intellect represents your ability to cast spells and use magic and increases your MP by .5 per level. 
Hunger represents how hungry you are and affects your ability to perform strenuous tasks and is increased by strenuous tasks. 
	hunger increases by 0.01 every action and decreases only when you eat food. if your hunger reaches 100, you will die.

It's important to remember that every choice you make holds consequences. 
Your decisions will directly shape the flow of your adventure, affecting both your immediate challenges and the unveiling of hidden secrets.
Proceed wisely, for your path is filled with challenges and secrets yet to be unveiled. 
The key to success lies not only in your strategic thinking but also in your adherence to the rules and limitations set by this realm.
May your journey be both thrilling and strategic as you navigate this richly detailed realm!

`
var adventureMessages = []goai.Message{}

// adventureCmd represents the adventure command
var adventureCmd = &cobra.Command{
	Use:   "adventure",
	Short: "lets you dive into a captivating text adventure",
	Long:  `immerses you in a dynamic virtual story. Through text prompts, you'll make choices that lead your character through a series of challenges and decisions. Each choice you make affects the storyline's development, creating a unique and interactive narrative experience. Get ready to explore, solve puzzles, and shape the adventure's outcome entirely through your imagination and decisions.`,
	Run: func(cmd *cobra.Command, args []string) {
		startAdventure()
	},
}

func init() {
	rootCmd.AddCommand(adventureCmd)
	adventureCmd.Flags().BoolVarP(&sayText, "say", "s", false, "Say text out loud (MacOS only)")
}

func adventureChat(prompt string) string {
	adventureMessages = append(adventureMessages, goai.Message{
		Role:    "user",
		Content: prompt,
	})
	oaiResponse, err := openai.ChatCompletion(adventureMessages)
	catchErr(err)
	adventureMessages = append(adventureMessages, goai.Message{
		Role:    "assistant",
		Content: oaiResponse.Choices[0].Message.Content,
	})
	return oaiResponse.Choices[0].Message.Content
}

func adventureImage(prompt, imageFile string) {
	fmt.Println("🖼  Creating Image...")
	res := openai.ImageGen(prompt, "", 1)

	url := res.Data[0].URL
	// fmt.Println("🌐 Image URL: " + url)

	promptFormatted := formatPrompt(prompt)
	filePath := viper.GetString("openAI_image_downloadPath")
	currentUser, err := user.Current()
	homeDir := currentUser.HomeDir
	catchErr(err)
	if filePath == `~` || strings.HasPrefix(filePath, "~") { // Replace ~ with home directory
		filePath = strings.Replace(filePath, "~", homeDir, 1)
	}

	fileName := promptFormatted[0:9] + strconv.Itoa(0) + ".jpg"
	fullFilePath := filepath.Join(filePath, fileName)
	// Create the directory (if it doesn't exist)
	err = os.MkdirAll(filePath, os.ModePerm)
	catchErr(err)
	// fmt.Printf("💾 Downloading Image:")
	url = httpDownloadFile(url, fullFilePath)
	// fmt.Printf(" \"%s\"\n", url)
	// fmt.Println("💻 Opening Image...")
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform for opening files: %s", runtime.GOOS)
	}
	if err != nil {
		trace()
		fmt.Println(err)
	}
}

func narratorSay(text string) {
	fmt.Println("\n🗣️ : ", text)
	if sayText {
		say(text)
	}
}

func getPlayerInput(player *Character) string {
	fmt.Print("\n🗡️ " + player.Name + "🛡️ : ")
	playerInput, err := getUserInput()
	catchErr(err)
	return playerInput
}

func startAdventure() {
	narratorSay("What is your name?")
	playerName, err := getUserInput()
	catchErr(err)

	player = Character{
		Name:        playerName,
		Description: "",
		HP:          100,
		MP:          100,
		Level:       1,
		Strength:    1,
		Defense:     1,
		Dexterity:   1,
		Intellect:   1,
	}

	narratorSay("Welcome " + player.Name + ", to the world of adventure! Describe your character, be as detailed as you like.")
	playerDescription := getPlayerInput(&player)
	player.Description = playerDescription

	playerString, err := json.Marshal(player)
	catchErr(err)

	adventureSystemMessage = adventureSystemMessage + "\n YOUR STARTING CHARACTER STATS:\n" + string(playerString)

	adventureMessages = []goai.Message{{
		Role:    "system",
		Content: adventureSystemMessage,
	}}

	startMessage := adventureChat("My name is " + player.Name + " start adventure")
	narratorSay(startMessage)
	// adventureImage(startMessage, startMessage)

	for {
		playerInput := getPlayerInput(&player)
		adventureResponse := adventureChat(playerInput)
		narratorSay(adventureResponse)
		// adventureImage(adventureResponse, adventureResponse)
	}
}
