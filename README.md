# Survey Hail
This program utilizes the Bandwidth API to send / receive sms messages.
The purpose was to support SMS survey's, ie, you give out a number to an audience, they all text what they think of something, and then later you text the survey number, from a whitelisted admin phone number the command "report" and it will send you back a MMS containing a CSV of all responses.
## Catapult info
You need to get a userId, token, and secret from the catapult API. https://catapult.inetwork.com
## Config file
See example config file in repo, sqldb is where the sqlite3db will be stored, admin is a list of numbers that are allowed to request a report, surveyNumber is the number you are monitoring, reportFile is where the .csv is saved before sending via mms to the number requesting it.
## Running service
I will be running this in docker, so that is what the Makefile and instructions are built for, if there is any interest in an rpm / systemd service files I can add them.
## Docker
To use the default run.sh in this repo, you will need to have a /opt/survey_hail directory, and a /opt/survey_hail/etc directory on the host system.
Place your config.yaml in /opt/survey_hail/etc, and for ease of use, your run.sh in /opt/survey_hail.
Then simply run run.sh (assuming you have Docker installed already).
## Install
run the ```make``` command, follow instructions
