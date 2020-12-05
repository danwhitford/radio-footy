import * as pug from "pug";
import * as fs from "fs";
const assert = require("assert").strict;

const data = fs.readFileSync("data.json", "utf8");
assert.ok(data.length > 0);
const compiledFunction = pug.compileFile("template.pug");

let comps = JSON.parse(data).flatMap(day => day.matches.map(match => match.competition)).filter((x, i, a) => a.indexOf(x) == i).sort()
let teams = JSON.parse(data).flatMap(day => day.matches.flatMap(match => match.title.split(" v "))).filter((x, i, a) => a.indexOf(x) == i).sort()
const filterOptions = [...comps, ...teams]

const site = compiledFunction({ matches: JSON.parse(data), filterOptions: filterOptions });

const fd = fs.openSync("site/index.html", "w");
fs.writeFileSync(fd, site);
fs.closeSync(fd);
