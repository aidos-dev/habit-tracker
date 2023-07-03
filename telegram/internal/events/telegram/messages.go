package telegram

const msgHelp = `I can save your habits and track them. You can manage your habits with comands from menu. 
The habit creation proces can be aborted with the comand /cancel from the menu`

const msgHello = "Hello! \n\n" + msgHelp

const (
	msgUnknownCommand     = "Unknown command 🤔"
	msgNoHabitCreated     = "Could not create a habit 😕"
	msgCreated            = "Habit created! 😄"
	msgHabitAlreadyExists = "This habit already exists 😬"
	msgHabitTitle         = "Please write the habit name"
	msgHabitDescription   = "Write short description for the habit"
	msgUnitOfMessure      = "What is the unit of messure for your habit?"
	msgGoal               = "How long do you want to follow this habit?"
	msgFrequency          = "Write the NUMBER of how many times a day you want to make your habit actions?"
	msgStartDate          = "Write the starting date for your habit in the format dd/mm/yyyy"
	msgEndDate            = "Write the end date for you habit in the format dd/mm/yyyy"
	timeFormat            = "02/01/2006"
)

/*

setting menu commands:
start - Start the bot
help - What can this bot do?
new_habit - Create a new habit
delete_habit - Delete a habit
cancel - Cancel the habit creation
*/
