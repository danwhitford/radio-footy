import * as fs from "fs";
import { from } from "rxjs";
import { filter, map, tap, toArray } from "rxjs/operators";
import { normaliseCompetitionName } from "./shared";

const fd = fs.openSync("feeds/bbc.json", "r");
const contents = fs.readFileSync(fd);
const schedule = JSON.parse(contents.toString("binary"));
fs.closeSync(fd);

const ret = [];

const observable = from(schedule).pipe(
  filter((prog: any) => {
    const { titles } = prog;
    return (
      titles.primary === "5 Live Sport" &&
      (titles.tertiary as string).includes(" v ") &&
      (titles.secondary as string).includes("Football")
    );
  }),
  map((prog) => {
    const { titles } = prog;
    return {
      station: "BBC Radio 5 Live",
      datetime: prog.start,
      title: titles.tertiary,
      competition: normaliseCompetitionName(titles.secondary),
    };
  }),
  toArray()
);

observable.subscribe((v) => {
  const outFd = fs.openSync("data/5live.json", "w");
  fs.writeFileSync(outFd, JSON.stringify(v));
});
