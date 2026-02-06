const ws = new WebSocket("ws://" + window.location.host + "/ws");

ws.onmessage = function (event) {
  const data = JSON.parse(event.data);

  if (data.type === "stats_update") {
    updateStatsDisplay(data.payload);
  } else if (data.type === "ticket_update") {
    showNotification(
      "Ticket " + data.payload.ticket_number + " - " + data.payload.status,
    );
  } else if (data.type === "counter_update") {
    showNotification(
      "Counter " + data.payload.name + " - " + data.payload.status,
    );
  }
};

function updateStatsDisplay(stats) {
  updateTodayStats(stats);
  updateOverallStats(stats);
}

function updateTodayStats(stats) {
  const elements = {
    totalTickets: document.querySelector("[data-total-tickets]"),
    currentlyServing: document.querySelector("[data-currently-serving]"),
    waitingTickets: document.querySelector("[data-waiting-tickets]"),
    activeCounters: document.querySelector("[data-active-counters]"),
    avgWaitTime: document.querySelector("[data-avg-wait-time]"),
    avgServiceTime: document.querySelector("[data-avg-service-time]"),
  };

  animateValue(elements.totalTickets, stats.TotalTicketsToday);
  animateValue(elements.currentlyServing, stats.CurrentlyServing);
  animateValue(elements.waitingTickets, stats.WaitingTickets);
  animateValue(elements.activeCounters, stats.ActiveCounters);
}

function updateOverallStats(stats) {
  updateQueueByCategory(stats.QueueLengthByCategory);
  updateHourlyDistribution(stats.HourlyDistribution);
}

function updateQueueByCategory(queueData) {
}

function updateHourlyDistribution(hourlyData) {
}

function animateValue(element, newValue) {
  if (!element) return;

  const startValue = parseInt(element.textContent) || 0;
  const duration = 1000;
  const startTime = performance.now();

  function update(currentTime) {
    const elapsed = currentTime - startTime;
    const progress = Math.min(elapsed / duration, 1);
    const currentValue = Math.floor(
      startValue + (newValue - startValue) * progress,
    );
    element.textContent = currentValue;

    if (progress < 1) {
      requestAnimationFrame(update);
    }
  }

  requestAnimationFrame(update);
}

function showNotification(message) {
  const toast = document.createElement("div");
  toast.className =
    "fixed bottom-4 right-4 bg-blue-600 text-white px-6 py-3 rounded-lg shadow-lg z-50 animate-pulse";
  toast.textContent = message;
  document.body.appendChild(toast);

  setTimeout(() => {
    toast.remove();
  }, 3000);
}

function refreshStats() {
  fetch("/admin/api/stats")
    .then((response) => response.json())
    .then((data) => updateStatsDisplay(data))
    .catch((error) => console.error("Failed to refresh stats:", error));
}

document.addEventListener("DOMContentLoaded", function () {
  const elements = document.querySelectorAll('[class*="text-3xl"]');
  elements.forEach((el) => {
    if (
      el.textContent.includes("Total Tiket") ||
      el.textContent.includes("Total Tickets")
    ) {
      el.setAttribute("data-total-tickets", "true");
    } else if (
      el.textContent.includes("Sedang Melayani") ||
      el.textContent.includes("Currently Serving")
    ) {
      el.setAttribute("data-currently-serving", "true");
    } else if (
      el.textContent.includes("Antrian Menunggu") ||
      el.textContent.includes("Waiting Tickets")
    ) {
      el.setAttribute("data-waiting-tickets", "true");
    } else if (
      el.textContent.includes("Loket Aktif") ||
      el.textContent.includes("Active Counters")
    ) {
      el.setAttribute("data-active-counters", "true");
    }
  });
});
