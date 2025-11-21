#### NOTE
Все ответа бота на русском языке

All bot responses in Russian

# Telegram BOT for automation the registration applications

Bot for automating the registration of applications for connection to a cellular operator and viewing tariff plans in chat.
The bot presented here is a training project written for a course assignment, and is also the first serious pet project.

## Stack of project
- Go (Golang) - version 1.24.4
- SQLite - version 3
- SMTP

### Packages
- tgbotapi v5.5.1 - Library for working with Telegram bot and Telegram API.
- yaml v2.4.0 - Library for reading yaml file
- SQLite driver from mattn v1.14.32 - A sqlite3 driver that conforms to the built-in database/sql interface.

## Bot functionality

All functionality is listed in the following format: Command name - brief description

### User 
- Все тарифы - Returns a list of all tariffs plans
- О тарифе - Return detail about selected tariff plans
- Создать заявку - Create application 
- Моя заявка - Return information from created application
- Изменить заявку - Edit information from created application

### Admin
- Все спрятанные - Return a list of all hidden tariffs plans
- Новый тариф - Create and add to database new tariffs plan
- Удалить тариф - Delete selected tariffs plan
- Изменить тариф - Edit tariffs plan
- Спрятать тариф - Hiding or Open tariff plan
- Изменить статус - Change status of user appeal
- Удалить заявку - Delete user appeal/contract

## Configuration 
Example configuration in the example.yaml file

# Known issues
1) FSM does not save state when restarting the bot
2) SQLite may become blocked under high loads.