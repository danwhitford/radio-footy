import * as fs from 'fs'
import bbcObservable from './get_bbc'
import tsObservable from './get_ts'

bbcObservable.subscribe((v) => {
    const outFd = fs.openSync("data/5live.json", "w");
    fs.writeFileSync(outFd, JSON.stringify(v));
});

tsObservable.subscribe((v) => {
    const outFd = fs.openSync("data/ts.json", "w");
    fs.writeFileSync(outFd, JSON.stringify(v));
});
