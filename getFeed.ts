import * as pug from "pug";
import * as fs from "fs";
import matchesObservable from "./observables/matchesObservable";

matchesObservable.subscribe((v) => {
  const fd = fs.openSync("./data.json", "w");
  fs.writeFileSync(fd, JSON.stringify(v));
  fs.closeSync(fd);
});
