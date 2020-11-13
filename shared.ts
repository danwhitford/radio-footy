const competitions = [
  "Premier League",
  "Champions League",
  "Europa League",
  "Championship",
  "FA Cup",
  "English Football League Trophy",
  "International Football",
];

export function normaliseCompetitionName(name: string) {
  for (let comp of competitions) {
    if (name.startsWith(comp)) {
      return comp;
    }
  }
  console.log("No conversion for", name);
  return name;
}

export function fromEntries<T>(entries: [keyof T, T[keyof T]][]): T {
  return entries.reduce(
    (acc, [key, value]) => ({ ...acc, [key]: value }),
    <T>{}
  );
}
