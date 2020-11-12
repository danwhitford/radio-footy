import * as fs from "fs";
import { zonedTimeToUtc } from "date-fns-tz";
import { normaliseCompetitionName } from "./shared";
import { from } from "rxjs";
import { filter, map, toArray } from "rxjs/operators";

const fd = fs.openSync("feeds/ts.json", "r");
const contents = fs.readFileSync(fd);
const schedule = JSON.parse(contents.toString("binary"));
fs.closeSync(fd);

const observable = from(schedule).pipe(
  filter(match => {
    return (match["livefeed"] as any[]).some(
      (feed) =>
        feed["feedname"] === "talkSPORT" || feed["feedname"] == "talkSPORT2"
    )
  }),
  map(match => {
    const d = new Date(match["Date"]);
    const utc = zonedTimeToUtc(d, "Europe/London");
    const channel = match["livefeed"]
      .map((feed) => feed["feedname"])
      .filter(
        (feedname) => feedname === "talkSPORT" || feedname == "talkSPORT2"
      )
      .pop();
    return {
      station: channel,
      datetime: utc,
      title: (match["title"] as string).split(": ").pop(),
      competition: normaliseCompetitionName(match["League"]),
    };
  }),
  toArray(),
)

observable.subscribe(v => {
  const outFd = fs.openSync("data/ts.json", "w");
  fs.writeFileSync(outFd, JSON.stringify(v));
})
