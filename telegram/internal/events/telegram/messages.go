package telegram

const msgHelp = `I can save your habits and track them
Please write a habit name and short description`

const msgHello = "Hello! \n\n" + msgHelp

const (
	msgUnknownCommand     = "Unknown command 🤔"
	msgNoHabitCreated     = "Could not create a habit 😕"
	msgCreated            = "Habit created! 😄"
	msgHabitAlreadyExists = "This habit already exists 😬"
)

/*

setting menu commands:
start - Start the bot
help - What can this bot do?
new habit - Create a new habit
delete habit - Delete a habit
*/
