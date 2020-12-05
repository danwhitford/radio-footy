import { concat, zip, of } from "rxjs";
import { mergeAll, toArray, map, groupBy, mergeMap } from "rxjs/operators";
import bbcObservable from "./bbc";
import talksportObservable from "./talksport";

function sortByDate(ar) {
  return ar.sort((m, mm) => {
    return new Date(m.datetime).valueOf() - new Date(mm.datetime).valueOf();
  });
}

function addDateString(m) {
  return {
    date: new Date(m.datetime).toLocaleDateString("en-GB", {
      timeZone: "Europe/London",
      weekday: "long",
      day: "numeric",
      month: "long",
    }),
    ...m,
  };
}

function addCalString(m) {
  const d = new Date(m.datetime);

  const d2 = new Date(d);
  d2.setHours(d2.getHours() + 2);

  const from = d
    .toISOString()
    .replace(/-/g, "")
    .replace(/:/g, "")
    .replace(/\./g, "");
  const to = d2
    .toISOString()
    .replace(/-/g, "")
    .replace(/:/g, "")
    .replace(/\./g, "");

  return {
    calString: `http://www.google.com/calendar/event?action=TEMPLATE&dates=${from}%2F${to}&text=${m.title}&location=${m.station}&details=${m.title}`,
    ...m,
  };
}

function addTimeString(m) {
  return {
    time: new Date(m.datetime).toLocaleTimeString("en-GB", {
      timeZone: "Europe/London",
      hour12: false,
      hour: "numeric",
      minute: "numeric",
    }),
    ...m,
  };
}

function prepareForPug(pair) {
  return {
    date: pair[0],
    matches: pair[1],
  };
}

function normaliseTitle(m) {
  const proper = m.title.replace(/\bvs\b/g, "v");
  return {
    ...m,
    title: proper,
  };
}

const matchObservable = concat(bbcObservable, talksportObservable).pipe(
  mergeAll(),
  toArray(),
  map(sortByDate),
  mergeAll(),
  map(addDateString),
  map(addCalString),
  map(addTimeString),
  map(normaliseTitle),
  groupBy((m) => m.date),
  mergeMap((group) => zip(of(group.key), group.pipe(toArray()))),
  map(prepareForPug),
  toArray()
);

export default matchObservable;
