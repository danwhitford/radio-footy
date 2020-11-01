const competitions = [
  "Premier League",
  "Champions League",
  "Europa League",
  "Championship",
  "FA Cup",
  "English Football League Trophy",
];

export function normaliseCompetitionName(name: string) {
  for (let comp of competitions) {
    if (name.startsWith(comp)) {
      return comp;
    }
  }
  throw new Error("Competition not recognised " + name);
}
