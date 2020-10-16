import * as pug from 'pug'
import * as fs from 'fs'

const files = ['data/ts.json', 'data/5live.json']
let matches = []

for(let file of files) {
    const fd = fs.openSync(file, 'r')
    const content = fs.readFileSync(fd)
    fs.closeSync(fd)
    const m = JSON.parse(content.toString('binary'))
    matches = matches.concat(m)
}

matches.sort((m, mm) => {
    return new Date(m.datetime).valueOf() - new Date(mm.datetime).valueOf()
})

const compiledFunction = pug.compileFile('template.pug')

const site = compiledFunction({matches})

const fd = fs.openSync('site/index.html', 'w')
fs.writeFileSync(fd, site)
fs.closeSync(fd)
