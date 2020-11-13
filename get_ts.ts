import * as fs from "fs";
import { zonedTimeToUtc } from "date-fns-tz";
import { normaliseCompetitionName } from "./shared";
import { from } from "rxjs";
import { filter, map, tap, toArray, mergeAll } from "rxjs/operators";
import { fetch } from 'cross-fetch'

const res = fetch('https://talksport.com/wp-json/talksport/v2/talksport-live/commentary')
  .then(r => r.json())

<<<<<<< HEAD
const observable = from(res).pipe(
  mergeAll(),
=======
const observable = from(schedule).pipe(
>>>>>>> 56b45867d5d5171de75741415b0a4c55ddfabfdc
  filter((match) => {
    return (match["livefeed"] as any[]).some(
      (feed) =>
        feed["feedname"] === "talkSPORT" || feed["feedname"] == "talkSPORT2"
    );
  }),
  map((match) => {
    const d = new Date(match["Date"]);
    const utc = zonedTimeToUtc(d, "Europe/London");
    const channel = match["livefeed"]
      .map((feed: { [x: string]: any }) => feed["feedname"])
      .filter(
        (feedname: string) =>
          feedname === "talkSPORT" || feedname == "talkSPORT2"
      )
      .pop();
    return {
      station: channel,
      datetime: utc,
      title: (match["title"] as string).split(": ").pop(),
      competition: normaliseCompetitionName(match["League"]),
    };
  }),
  toArray()
);

export default observable;
