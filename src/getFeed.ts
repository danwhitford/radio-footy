import * as fs from "fs";
import matchesObservable from "./observables/matchesObservable";

matchesObservable.subscribe((v) => {
  const fd = fs.openSync("./docs/data.json", "w");
  fs.writeFileSync(fd, JSON.stringify(v));
  fs.closeSync(fd);
});
