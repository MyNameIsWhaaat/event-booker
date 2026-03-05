function snapshotInputs() {
  const map = new Map();
  document.querySelectorAll('#events input').forEach((inp) => {
    if (inp.id) map.set(inp.id, inp.value);
  });
  return map;
}

function restoreInputs(map) {
  document.querySelectorAll('#events input').forEach((inp) => {
    if (!inp.id) return;
    if (map.has(inp.id)) inp.value = map.get(inp.id);
  });
}

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
const saved = snapshotInputs();   // ✅ сохранить введённое

  const res = await fetch("/events");
  const data = await res.json();
  const items = pick(data, [["items"], ["Items"]], []);

  const tbody = document.querySelector("#events tbody");
  tbody.innerHTML = "";

  items.forEach((row) => {
    const e = getEventObj(row);
    const s = getStatsObj(row);

    const id = pick(e, [["ID"], ["Id"], ["id"]], "");
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

      <td>
        <input placeholder="email" id="email-${id}">
        <button data-event="${id}" class="book-btn">Book</button>
      </td>

      <td>
        <input placeholder="booking_id" id="bid-${id}" size="40">
        <button data-event="${id}" class="confirm-btn">Confirm</button>
      </td>
    `;
    tbody.appendChild(tr);
    restoreInputs(saved); 
  });
}

async function book(eventID) {
  const email = document.getElementById(`email-${eventID}`)?.value || "";
  const msg = document.getElementById("message");

  const res = await fetch(`/events/${eventID}/book`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_email: email }),
  });

  if (res.ok) {
    const data = await res.json();
    msg.innerText = `Booked! booking_id=${data.booking_id || data.BookingID || ""}`;
  } else {
    const t = await res.text();
    msg.innerText = `Booking failed: ${t}`;
  }

  await loadEvents();
}

async function confirmBooking(eventID) {
  const bookingID = document.getElementById(`bid-${eventID}`)?.value || "";
  const msg = document.getElementById("message");

  const res = await fetch(`/events/${eventID}/confirm`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ booking_id: bookingID }),
  });

  if (res.ok) {
    msg.innerText = "Confirmed";
  } else {
    const t = await res.text();
    msg.innerText = `Confirm failed: ${t}`;
  }

  await loadEvents();
}

document.addEventListener("click", (e) => {
  const bookBtn = e.target.closest(".book-btn");
  if (bookBtn) return book(bookBtn.dataset.event);

  const confBtn = e.target.closest(".confirm-btn");
  if (confBtn) return confirmBooking(confBtn.dataset.event);
});

loadEvents();
setInterval(loadEvents, 4000);