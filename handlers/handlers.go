import "github.com/bwmarrin/discordgo"

var Handlers = []func(s *discordgo.Session, _ any){
	ready,
}
