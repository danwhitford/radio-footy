import * as pug from "pug";
import * as fs from "fs";
const assert = require("assert").strict;

function urlEncode(s: string) {
  return encodeURI(s.toLocaleLowerCase().replace(/[ ]+/g, "-"));
}

function build(buildFn, data, out) {
  const page = buildFn(data);
  const fd = fs.openSync(out, "w");
  fs.writeFileSync(fd, page);
  fs.closeSync(fd);
}

function buildListingsPage(template: string, data: any, out: string) {
  const fn = pug.compileFile(template);
  build(
    fn,
    {
      matches: data,
      competitions,
      stations,
      teams,
      urlEncode,
    },
    out
  );
}

function buildDirIndexPage(
  template: string,
  directory: string,
  data: any,
  out: string
) {
  const fn = pug.compileFile(template);
  build(
    fn,
    {
      directory: directory,
      pages: data,
      competitions,
      stations,
      teams,
      urlEncode,
    },
    out
  );
}

function buildFilteredListingsPage(
  template: string,
  data: any,
  outDir: string,
  els: any[],
  filterFunc: any
) {
  for (const el of els) {
    const filtered = data
      .map((day) => {
        return {
          ...day,
          matches: day.matches.filter(filterFunc(el)),
        };
      })
      .filter((day) => day.matches.length > 0);
    const addy = urlEncode(el);
    buildListingsPage(template, filtered, `${outDir}/${addy}.html`);
  }
}

const data = fs.readFileSync("data.json", "utf8");
assert.ok(data.length > 0);
const siteData = JSON.parse(data);

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

buildListingsPage("templates/listings.pug", siteData, "site/index.html");
buildFilteredListingsPage(
  "templates/listings.pug",
  siteData,
  "site/teams",
  teams,
  (el) => (match) => match.title.split(" v ").includes(el)
);
buildFilteredListingsPage(
  "templates/listings.pug",
  siteData,
  "site/stations",
  teams,
  (el) => (match) => match.station === el
);
buildFilteredListingsPage(
  "templates/listings.pug",
  siteData,
  "site/competitions",
  teams,
  (el) => (match) => match.competition === el
);

buildDirIndexPage(
  "templates/dir-index.pug",
  "stations",
  stations,
  "site/stations/index.html"
);
buildDirIndexPage(
  "templates/dir-index.pug",
  "teams",
  teams,
  "site/teams/index.html"
);
buildDirIndexPage(
  "templates/dir-index.pug",
  "competitions",
  competitions,
  "site/competitions/index.html"
);
