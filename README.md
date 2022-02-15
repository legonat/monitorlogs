# monitorlogs

Golang/React project for log files parsing.

IMPORTANT! Build with --tags "fts5" For Full Text Search Support

Flags list:
* -f initUsersDb -p "./data/db" (Initializes Sqlite database for keeping registered users) 
* -f initLogsDb -p "./data/db" (Initializes Sqlite database for keeping parsed logs)
* -f read -fn "./data/debug.log" (Read logfile with exact name)
* -f readFolder -p "./data" (Read folder with several logfiles)
* -f server (Starts server with port that is defined in config file)

Usage:
1. Create logs.conf file with fields described in configPattern file.
2. Init User and Logs databases.
3. Read Logfile(s) (Uses space (" ") as separator. Example log string pattern: 2020-07-16T12:30:51Z GUI: requestInitialize)
4. Run server. Server port should be specified in Config file
