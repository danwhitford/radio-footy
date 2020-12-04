import * as pug from "pug";
import * as fs from "fs";
const assert = require('assert').strict;

const data = fs.readFileSync('data.json', 'utf8')
assert.ok(data.length > 0)
const compiledFunction = pug.compileFile("template.pug");

const site = compiledFunction({ matches: JSON.parse(data) });

const fd = fs.openSync("site/index.html", "w");
fs.writeFileSync(fd, site);
fs.closeSync(fd);
