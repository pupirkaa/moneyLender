# Money lender
💸 Simplify debt distribution among groups of friends! 

# Purpose
Money Lender is a study backend project inspired by Splitwise, designed to simplify debt distribution among groups of people. Offering 3 options of data storage - in file, in memory and in sqlite db. Users can interact with Money Lender through a browser using simple HTML templates, or leverage its REST API to connect with a client application

# Running the Program
There are few flags you should set to run programm correctly:
* -storage - use that to choice type of storage: inmem(in memory), fs(in file), sqlite (in sqlite db). 
* -users, -sessions, -txs - flags for file storage with path to files with users, sessions and transactions. 
* -sqlite - flag for sqlite storage with path to file with db
 
  Use `go run .` + flags to run the program or `go build` to build the program binary.
