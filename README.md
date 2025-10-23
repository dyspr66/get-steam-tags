# Get Steam Tags

The goal of this program is to dynamically obtain all steam games and corresponding tags. This data will be used to generate an excel file with an adjacency matrix between games and their tags.

## Guide

Currently, this guide is for Windows.

1. Download this repository.
2. Set up Go and libraries
    1. Make sure you have the [Go Language](https://go.dev/doc/install) installed.
    2. Open a terminal using PowerShell. This will be where you'll run all the following commands.
    3. Run `cd "path\to\get-steam-tags\folder\here\"` Replace the path accordingly. After running, your terminal should be within the same folder/directory this README is in.
    4. Run `go get .` This will obtain all libraries necessary for this program.
3. Set up API Key
    1. Get a [Steam API key](https://steamcommunity.com/login/home/?goto=%2Fdev%2Fapikey). This will be used to obtain a list of all games from steampowered.com.
    2. Run `New-Item .env` This will create a file named `.env`
    3. Run `notepad.exe .env` This will open up Notepad so you can edit the .env file. Add the following line in the .env, then save. (make sure to replace ABCDE12345 with your actual API key.)

```
STEAM_WEB_API_KEY=ABCDE12345
```

3. Finally, run `go run .` This will start the program. As of writing, it takes an estimated 6-9 hours to run. You should see logs of what's been done in your terminal.

## Resources

-   [xPaw's Steam API Docs](https://steamapi.xpaw.me/)
    -   https://steamapi.xpaw.me/#IStoreService/GetAppList
-   steampowered.com
    -   Obtain API Key [here](https://steamcommunity.com/login/home/?goto=%2Fdev%2Fapikey).
-   SteamSpy
    -   https://steamspy.com/api.php?request=appdetails&appid=APP_ID_HERE
