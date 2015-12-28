Overpower is a multiplayer simultanious-action turn-based hex-grid game of conquering a galaxy.

The tech used is postgres database for persistance, simple html for UI served over http, and cookie based login/logouts.  It is all terribly insecure.

Overpower is for 2-8 players, played out over several days.  When all players have submitted their orders for the turn, or some time limit is met, all orders are executed and the next turn begins.  The game starts with each player in control of one planet on the edge of a large galaxy, and on each turn the players may construct space ships to travel to and colonize more planets.

How To Use This Code:
    PostGreSQL must be installed, and the overpower/db package needs a file with the database name, user name, and password you wish to use.  Running go test in the overpower/db directory with makeTables_test.go variable UPDATETABLES set to true sets up the proper tables AND DROPS ALL EXISTING DATA IN OLDER TABLES OF THOSE NAMES.
    The overpower/server package builds an executable that requires the TEMPLATES, DATA, and STATIC directories in the overpower/server directory.  When run, it starts a http server that allows browsers to connect, login, and start/play games of Overpower.

ATTRIBUTIONS:

Golang standard library (It's Amazing!):
    https://golang.org/pkg/
Database: PostGreSQL
    http://www.postgresql.org/
Drivers: PostGreSQL Go Drivers
    https://github.com/lib/pq
Graphics library: Draw2D
    https://github.com/llgcode/draw2d/
Font: Droid Sans Mono 
    http://www.fontsquirrel.com/fonts/droid-sans-mono
