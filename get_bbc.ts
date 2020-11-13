import { from } from "rxjs";
import { filter, map, mergeAll, mergeMap, pluck, tap, toArray } from "rxjs/operators";
import { normaliseCompetitionName } from "./shared";
import { add, format } from 'date-fns'
import { fetch } from 'cross-fetch'

const baseUrl = 'https://rms.api.bbc.co.uk/v2/experience/inline/schedules/bbc_radio_five_live/'

function * generateWeek() {
  const today = new Date()
  for (let i = 0; i <= 7; ++i) {
    const d = add(today, { days: i })
    yield format(d, 'Y-M-dd')
  }
}

const programs = from(generateWeek()).pipe(
  map(v => baseUrl + v),
  mergeMap(url => fetch(url)),
  mergeMap(res => res.json()),
  pluck('data'),
  mergeAll(),
  pluck('data'),
  mergeAll(),
)

const observable = from(programs).pipe(
  filter((prog: any) => {
    const { titles } = prog;
    return (
      titles.primary === "5 Live Sport" &&
      (titles.tertiary as string).includes(" v ") &&
      (titles.secondary as string).includes("Football")
    );
  }),
  map((prog) => {
    const { titles } = prog;
    return {
      station: "BBC Radio 5 Live",
      datetime: prog.start,
      title: titles.tertiary,
      competition: normaliseCompetitionName(titles.secondary),
    };
  }),
  toArray()
);

export default observable
