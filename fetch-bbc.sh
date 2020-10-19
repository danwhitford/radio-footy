#! /bin/bash

for i in $(seq 0 7)
do
    date=$(date -d "+$i days" "+%Y-%m-%d")
    wget "https://rms.api.bbc.co.uk/v2/experience/inline/schedules/bbc_radio_five_live/$date" -O "feeds/bbc/5live-$date.json"
    jq -s "[.[].data | .[].data | .[]]" feeds/bbc/5live*.json > feeds/bbc.json
done
