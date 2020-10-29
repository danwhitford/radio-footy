
const competitions = ["Premier League", "Champions League", "Europa League", "Championship"]

export function competitionNormaliser(name: string) {
    for(let comp of competitions) {
        if (name.startsWith(comp)) {
            return comp
        }
    }
    throw new Error("Competition not recognised")
}
