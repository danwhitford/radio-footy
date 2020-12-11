import * as pug from "pug";
import * as fs from "fs";
const assert = require("assert").strict;

const data = fs.readFileSync("data.json", "utf8");
assert.ok(data.length > 0);

function urlEncode(s: string) {
  return encodeURI(s.toLocaleLowerCase().replace(/[ ]+/g, '-'))
}

let competitions = JSON.parse(data)
  .flatMap((day) => day.matches.map((match) => match.competition))
  .filter((x, i, a) => a.indexOf(x) == i)
  .sort();
let teams = JSON.parse(data)
  .flatMap((day) => day.matches.flatMap((match) => match.title.split(" v ")))
  .filter((x, i, a) => a.indexOf(x) == i)
  .sort();
let stations = JSON.parse(data)
  .flatMap((day) => day.matches.map((match) => match.station))
  .filter((x, i, a) => a.indexOf(x) == i)
  .sort();

const compiledFunction = pug.compileFile("templates/listings.pug");
const site = compiledFunction({
  matches: JSON.parse(data),
  competitions,
  stations,
  teams,
  urlEncode,
});

let fd = fs.openSync("site/index.html", "w");
fs.writeFileSync(fd, site);
fs.closeSync(fd);

const dirTemplate = pug.compileFile('templates/dir-index.pug')
let dirIndex = dirTemplate({
  directory: 'stations',
  pages: stations,
  competitions,
  stations,
  teams,
  urlEncode,
})
fd = fs.openSync("site/stations/index.html", "w");
fs.writeFileSync(fd, dirIndex);
fs.closeSync(fd);

dirIndex = dirTemplate({
  directory: 'competitions',
  pages: competitions,
  competitions,
  stations,
  teams,
  urlEncode,
})
fd = fs.openSync("site/competitions/index.html", "w");
fs.writeFileSync(fd, dirIndex);
fs.closeSync(fd);

dirIndex = dirTemplate({
  directory: 'teams',
  pages: teams,
  competitions,
  stations,
  teams,
  urlEncode,
})
fd = fs.openSync("site/teams/index.html", "w");
fs.writeFileSync(fd, dirIndex);
fs.closeSync(fd);

const siteData = JSON.parse(data)
for (const team of teams) {
  const filtered = siteData
    .map(day => {
      return {
        ...day,
        matches: day.matches.filter(match => match.title.split(" v ").includes(team))
      }})
    .filter(day => day.matches.length > 0)

    const site = compiledFunction({
      matches: filtered,
      competitions,
      stations,
      teams,
  urlEncode,
    });
    
    const addy = urlEncode(team)
    let fd = fs.openSync(`site/teams/${addy}.html`, "w");
    fs.writeFileSync(fd, site);
    fs.closeSync(fd);
}

for (const comp of competitions) {
  const filtered = siteData
    .map(day => {
      return {
        ...day,
        matches: day.matches.filter(match => match.competition === comp)
      }})
    .filter(day => day.matches.length > 0)

    const site = compiledFunction({
      matches: filtered,
      competitions,
      stations,
      teams,
  urlEncode,
    });
    
    const addy = urlEncode(comp)
    let fd = fs.openSync(`site/competitions/${addy}.html`, "w");
    fs.writeFileSync(fd, site);
    fs.closeSync(fd);
}

for (const station of stations) {
  const filtered = siteData
    .map(day => {
      return {
        ...day,
        matches: day.matches.filter(match => match.station === station)
      }})
    .filter(day => day.matches.length > 0)

    const site = compiledFunction({
      matches: filtered,
      competitions,
      stations,
      teams,
  urlEncode,
    });
    
    const addy = urlEncode(station)
    let fd = fs.openSync(`site/stations/${addy}.html`, "w");
    fs.writeFileSync(fd, site);
    fs.closeSync(fd);
}
