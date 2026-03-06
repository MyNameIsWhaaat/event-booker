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

function formatDate(value) {
  if (!value) return "";
  return String(value);
}

const openedBookings = new Set();

async function loadEvents() {
  const res = await fetch("/events");
  const data = await res.json();
  const items = pick(data, [["items"], ["Items"]], []);

  const tbody = document.querySelector("#events tbody");
  tbody.innerHTML = "";

  for (const row of items) {
    const e = getEventObj(row);
    const s = getStatsObj(row);

    const id = pick(e, [["ID"], ["Id"], ["id"]], "");
    const title = pick(e, [["Title"], ["title"]], "");
    const starts = pick(e, [["StartsAt"], ["starts_at"], ["startsAt"]], "");
    const capacity = pick(e, [["Capacity"], ["capacity"]], "");

    const pending = pick(s, [["pending"], ["Pending"]], 0);
    const confirmed = pick(s, [["confirmed"], ["Confirmed"]], 0);
    const free = pick(s, [["free_seats"], ["FreeSeats"], ["freeSeats"]], 0);

    const eventRow = document.createElement("tr");
    eventRow.innerHTML = `
      <td>${fmt(title)}</td>
      <td>${fmt(starts)}</td>
      <td>${fmt(capacity)}</td>
      <td>${fmt(pending)}</td>
      <td>${fmt(confirmed)}</td>
      <td><b>${fmt(free)}</b></td>
      <td>
        <button class="show-bookings-btn" data-event-id="${id}">
          ${openedBookings.has(id) ? "Hide bookings" : "Show bookings"}
        </button>
      </td>
    `;
    tbody.appendChild(eventRow);

    const bookingsRow = document.createElement("tr");
    bookingsRow.id = `bookings-row-${id}`;
    bookingsRow.style.display = openedBookings.has(id) ? "table-row" : "none";
    bookingsRow.innerHTML = `
      <td colspan="7">
        <div id="bookings-container-${id}">
          ${openedBookings.has(id) ? "Loading..." : ""}
        </div>
      </td>
    `;
    tbody.appendChild(bookingsRow);
  }

  for (const eventID of openedBookings) {
    const row = document.getElementById(`bookings-row-${eventID}`);
    if (row) {
      await loadBookings(eventID);
    }
  }
}

async function loadBookings(eventID) {
  const container = document.getElementById(`bookings-container-${eventID}`);
  if (!container) return;

  container.innerHTML = "Loading...";

  const res = await fetch(`/events/${eventID}/bookings`);
  if (!res.ok) {
    container.innerHTML = "Failed to load bookings";
    return;
  }

  const data = await res.json();
  const items = pick(data, [["items"], ["Items"]], []);

  if (!items.length) {
    container.innerHTML = "<p>No bookings yet</p>";
    return;
  }

  let html = `
    <table class="nested-table">
      <thead>
        <tr>
          <th>Booking ID</th>
          <th>User Email</th>
          <th>Status</th>
          <th>Created</th>
          <th>Expires</th>
          <th>Confirmed</th>
          <th>Cancelled</th>
        </tr>
      </thead>
      <tbody>
  `;

  items.forEach((b) => {
    const id = pick(b, [["ID"], ["Id"], ["id"]], "");
    const email = pick(b, [["UserEmail"], ["user_email"], ["userEmail"]], "");
    const status = pick(b, [["Status"], ["status"]], "");
    const createdAt = pick(b, [["CreatedAt"], ["created_at"], ["createdAt"]], "");
    const expiresAt = pick(b, [["ExpiresAt"], ["expires_at"], ["expiresAt"]], "");
    const confirmedAt = pick(b, [["ConfirmedAt"], ["confirmed_at"], ["confirmedAt"]], "");
    const cancelledAt = pick(b, [["CancelledAt"], ["cancelled_at"], ["cancelledAt"]], "");

    html += `
      <tr>
        <td>${fmt(id)}</td>
        <td>${fmt(email)}</td>
        <td>${fmt(status)}</td>
        <td>${formatDate(createdAt)}</td>
        <td>${formatDate(expiresAt)}</td>
        <td>${formatDate(confirmedAt)}</td>
        <td>${formatDate(cancelledAt)}</td>
      </tr>
    `;
  });

  html += `
      </tbody>
    </table>
  `;

  container.innerHTML = html;
}

async function toggleBookings(eventID) {
  const row = document.getElementById(`bookings-row-${eventID}`);
  if (!row) return;

  const isHidden = row.style.display === "none";

  if (isHidden) {
    row.style.display = "table-row";
    openedBookings.add(eventID);
    await loadBookings(eventID);
  } else {
    row.style.display = "none";
    openedBookings.delete(eventID);
  }
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
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  if (res.ok) {
    msg.innerText = "Event created";
    document.getElementById("createForm").reset();
  } else {
    const text = await res.text();
    msg.innerText = `Error creating event: ${text}`;
  }

  await loadEvents();
};

document.addEventListener("click", async (e) => {
  const btn = e.target.closest(".show-bookings-btn");
  if (!btn) return;

  const eventID = btn.dataset.eventId;
  await toggleBookings(eventID);

  btn.innerText = openedBookings.has(eventID) ? "Hide bookings" : "Show bookings";
});

loadEvents();
setInterval(loadEvents, 5000);