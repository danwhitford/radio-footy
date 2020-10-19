import * as pug from "pug";
import * as fs from "fs";

const files = ["data/ts.json", "data/5live.json"];
let matches = [];

for (let file of files) {
  const fd = fs.openSync(file, "r");
  const content = fs.readFileSync(fd);
  fs.closeSync(fd);
  const m = JSON.parse(content.toString("binary"));
  matches = matches.concat(m);
}

matches.sort((m, mm) => {
  return new Date(m.datetime).valueOf() - new Date(mm.datetime).valueOf();
});

const rolledMatches = {};
for (let match of matches) {
  const d = new Date(Date.parse(match.datetime));

  const date = new Date(match.datetime).toLocaleDateString("en-GB", {
    timeZone: "Europe/London",
    weekday: "long",
    day: "numeric",
    month: "long",
  });
  match.time = new Date(match.datetime).toLocaleTimeString("en-GB", {
    timeZone: "Europe/London",
    hour12: false,
    hour: "numeric",
    minute: "numeric",
  });

  const d2 = new Date(d);
  d2.setHours(d2.getHours() + 2);

  const from = d
    .toISOString()
    .replace(/-/g, "")
    .replace(/:/g, "")
    .replace(/\./g, "");
  const to = d2
    .toISOString()
    .replace(/-/g, "")
    .replace(/:/g, "")
    .replace(/\./g, "");

  match.calString = `http://www.google.com/calendar/event?action=TEMPLATE&dates=${from}%2F${to}&text=${match.title}&location=${match.station}&details=${match.title}`;
  if (date in rolledMatches) {
    rolledMatches[date].push(match);
  } else {
    rolledMatches[date] = [match];
  }
}

const rolledMatchesArray = Object.keys(rolledMatches).map((ob) => {
  return {
    date: [ob],
    matches: rolledMatches[ob],
  };
});

const compiledFunction = pug.compileFile("template.pug");

const site = compiledFunction({ matches: rolledMatchesArray });

const fd = fs.openSync("site/index.html", "w");
fs.writeFileSync(fd, site);
fs.closeSync(fd);
