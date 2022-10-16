# gopass
gopass is a password database manager. It uses Google Sheets as a database.

## Setup
1. Create a new Google Sheet: https://sheets.new.
2. Add three columns, first column is "site", second one is "username" and last one is "password"
3. Save the ID of your Sheet from URL (it’s a random string that looks like `1SMeoyesCaGHRlYdGj9VyqD-qhXtab1jrcgHZ0irvNDs`)
4. Go to https://console.cloud.google.com/apis/credentials, Create new Service Account.
5. Save credentials of this account as a json file
6. Create `.gopass` folder in `$HOME` directory. 
7. Move `credentials.json` to `$HOME/.gopass/`
8. Create a json file named `conf.json` in same folder, and paste your Sheet ID like `{"sheetId":"YOUR_SHEET_ID"}`
9. Go to your Google Sheets, click "Share" and give this email address "Editor" access on your sheet.
10. Go to https://console.developers.google.com/apis/api/sheets.googleapis.com/overview and make sure the Google Sheets API is enabled.




## Demo
[![Demo](https://img.youtube.com/vi/gLtd0f_lEhA/0.jpg)](https://www.youtube.com/watch?v=gLtd0f_lEhA)

## Configuration of Sheets API
coming soon.
