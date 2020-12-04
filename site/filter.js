window.onload = (function () {
    document.getElementById('filter-input').oninput = (function (e) {
        const filter = e.target.value
        const rows = document.getElementsByTagName("tr")
        for (row of rows) {
            const comp = row.getAttribute('data-competition')
            const teams = row.getAttribute('data-teams').split(' ').filter(t => !['v', 'vs'].includes(t) )
            if (comp.startsWith(filter) || teams.some(t => t.startsWith(filter))) {
                row.hidden = false
            } else {
                row.hidden = true
            }
        }

        const tables = document.getElementsByTagName('table')
        for (table of tables) {
            const rows = [...table.getElementsByTagName('tr')]
            hide = true
            for(row of rows) {
                if (!row.hidden) {
                    hide = false
                }
            }
            console.log(hide)
            table.hidden = hide
            const header = document.querySelector(`h2[data-date="${table.getAttribute('data-date')}"]`)
            header.hidden = hide
        }
    })
})
