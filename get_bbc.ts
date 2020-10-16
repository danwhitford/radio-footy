import * as fs from 'fs'
import * as cheerio from 'cheerio'

const fd = fs.openSync('feeds/5live.html', 'r')
const contents = fs.readFileSync(fd)
const $ = cheerio.load(contents)
fs.closeSync(fd)

const ret = []

const days = $('.schedule-col')
days.each((_i, day) => {
    const programs = $('.schedule-program-row', day)
    programs.each((i, program) => {
        const description = $('.schedule-program-descr', program)
        description.each((i, d) => {
            const prem = (d.children.filter(c => c?.data?.startsWith('Premier League Football')))
            if (prem.length > 0) {
                const d = new Date(parseInt(program.attribs['data-time']) * 1000)
                const title = prem[0].data.split(',')[1].trim()
                ret.push({
                    station: 'BBC Radio 5 Live',
                    datetime: d,
                    title: title,
                })
            }
        })
    })
})

const outFd = fs.openSync('data/5live.json', 'w')
fs.writeFileSync(outFd, JSON.stringify(ret))
