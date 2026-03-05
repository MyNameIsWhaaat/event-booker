function pick(obj, paths, fallback = undefined) {
  for (const p of paths) {
    let cur = obj;
    let ok = true;
    for (const key of p) {
      if (cur && Object.prototype.hasOwnProperty.call(cur, key)) {
        cur = cur[key];
      } else {
        ok = false;
        break;
      }
    }
    if (ok) return cur;
  }
  return fallback;
}

function getEventObj(row) {
  return pick(row, [["Event"], ["event"]], row);
}

function getStatsObj(row) {
  return pick(row, [["stats"], ["Stats"]], {});
}

function fmt(v) {
  if (v === null || v === undefined) return "";
  return String(v);
}

async function loadEvents() {
  const res = await fetch("/events");
  const data = await res.json();
  const items = pick(data, [["items"], ["Items"]], []);

  const tbody = document.querySelector("#events tbody");
  tbody.innerHTML = "";

  items.forEach((row) => {
    const e = getEventObj(row);
    const s = getStatsObj(row);

    const title = pick(e, [["Title"], ["title"]], "");
    const starts = pick(e, [["StartsAt"], ["starts_at"], ["startsAt"]], "");
    const cap = pick(e, [["Capacity"], ["capacity"]], "");

    const pending = pick(s, [["pending"], ["Pending"]], 0);
    const confirmed = pick(s, [["confirmed"], ["Confirmed"]], 0);
    const free = pick(s, [["free_seats"], ["FreeSeats"], ["freeSeats"]], 0);

    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${fmt(title)}</td>
      <td>${fmt(starts)}</td>
      <td>${fmt(cap)}</td>
      <td>${fmt(pending)}</td>
      <td>${fmt(confirmed)}</td>
      <td><b>${fmt(free)}</b></td>
    `;
    tbody.appendChild(tr);
  });
}

document.getElementById("createForm").onsubmit = async (e) => {
  e.preventDefault();

  const msg = document.getElementById("message");
  const body = {
    title: document.getElementById("title").value,
    starts_at: document.getElementById("starts").value,
    capacity: Number(document.getElementById("capacity").value),
    booking_ttl_seconds: Number(document.getElementById("ttl").value),
    requires_payment: document.getElementById("payment").checked,
  };

  const res = await fetch("/events", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (res.ok) {
    msg.innerText = "Event created";
  } else {
    const t = await res.text();
    msg.innerText = `Error creating event: ${t}`;
  }

  await loadEvents();
};

loadEvents();
setInterval(loadEvents, 4000);