import * as fs from "fs";
import { title } from "process";

const fd = fs.openSync("feeds/bbc.json", "r");
const contents = fs.readFileSync(fd);
const schedule = JSON.parse(contents.toString("binary"));
fs.closeSync(fd);

const ret = [];

for (let program of schedule) {
  const { titles } = program;
  if (
    titles.primary === "5 Live Sport" &&
    (titles.tertiary as string).includes(" v ") &&
    [
      "Premier League Football 2020-21",
      "Champions League Football 2020-21",
    ].includes(titles.secondary)
  ) {
    ret.push({
      station: "BBC Radio 5 Live",
      datetime: program.start,
      title: titles.tertiary,
      competition: titles.secondary,
    });
  }
}

const outFd = fs.openSync("data/5live.json", "w");
fs.writeFileSync(outFd, JSON.stringify(ret));
