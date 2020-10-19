import * as fs from 'fs'
import * as cheerio from 'cheerio'

const fd = fs.openSync('feeds/bbc.json', 'r')
const contents = fs.readFileSync(fd)
const schedule = JSON.parse(contents.toString('binary'))
fs.closeSync(fd)

const ret = []

for(let program of schedule) {
    const {titles} = program
    if (titles.primary === "5 Live Sport" && titles.secondary === "Premier League Football 2020-21") {
        ret.push({
            station: 'BBC Radio 5 Live',
            datetime: program.start,
            title: titles.tertiary,
        })
    }
}

const outFd = fs.openSync('data/5live.json', 'w')
fs.writeFileSync(outFd, JSON.stringify(ret))
