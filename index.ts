import * as pug from "pug";
import * as fs from "fs";
import matchesObservable from "./observables/matchesObservable";

matchesObservable.subscribe((v) => {
  const compiledFunction = pug.compileFile("template.pug");

  const site = compiledFunction({ matches: v });

  const fd = fs.openSync("site/index.html", "w");
  fs.writeFileSync(fd, site);
  fs.closeSync(fd);
});
