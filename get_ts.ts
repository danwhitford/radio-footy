import * as fs from 'fs'

const fd = fs.openSync('feeds/ts.json', 'r')
const contents = fs.readFileSync(fd)
const schedule = JSON.parse(contents.toString('binary'))
fs.closeSync(fd)

const ret = []
for(let match of schedule) {
    if (match['League'] === 'Premier League' && (match['livefeed'] as any[]).some(feed => feed['feedname'] === 'talkSPORT')) {

        const d = new Date(match['Date'])
        ret.push({
            station: 'TalkSport',
            datetime: d,
            title: (match['title'] as string).split(': ').pop()
        })
    }
}

const outFd = fs.openSync('data/ts.json', 'w')
fs.writeFileSync(outFd, JSON.stringify(ret))
