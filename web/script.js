// Get DOM elements
const numConnPerSec = document.getElementById("numConnPerSec");
const numActiveConn = document.getElementById("numActiveConn");
const numIpsPerSec = document.getElementById("numIpsPerSec");
const inboundMBps = document.getElementById("inboundMBps");
const chartCanvas = document.getElementById("chart");

const limit = 120;
const chartCtx = chartCanvas.getContext("2d");

const chart = new Chart(chartCtx, {
  type: "line",
  data: {
    labels: [],
    datasets: [
      {
        label: "Connections Per Second",
        data: [],
        backgroundColor: "rgba(255, 99, 132, 0.2)",
        borderColor: "rgba(255, 99, 132, 1)",
        borderWidth: 4,
        fill: true,
        pointRadius: 0,
        pointHoverRadius: 3,
      },
      {
        label: "Active Connections",
        data: [],
        backgroundColor: "rgba(54, 162, 235, 0.2)",
        borderColor: "rgba(54, 162, 235, 1)",
        borderWidth: 4,
        fill: true,
        pointRadius: 0,
        pointHoverRadius: 3,
      },
      {
        label: "IPs Per Sec",
        data: [],
        backgroundColor: "rgba(137, 235, 137, 0.2)",
        borderColor: "rgba(137, 235, 137, 1)",
        borderWidth: 4,
        fill: true,
        pointRadius: 0,
        pointHoverRadius: 3,
      },
      {
        label: "Inbound MBps",
        data: [],
        backgroundColor: "rgba(255, 206, 86, 0.2)",
        borderColor: "rgba(255, 206, 86, 1)",
        borderWidth: 4,
        fill: true,
        pointRadius: 0,
        pointHoverRadius: 3,
      },
    ],
  },
  options: {
    responsive: true,
    scales: {
      xAxes: [{ display: false }],
      yAxes: [{ display: true, ticks: { beginAtZero: true } }],
    },
  },
});

const ws = new WebSocket(`ws://${window.location.hostname}:8080`);

// Handle WebSocket messages
ws.addEventListener("message", (event) => {
  const data = JSON.parse(event.data);
  console.log(data);
  // Update text content of elements
  numConnPerSec.textContent = data.numConnPerSec;
  numActiveConn.textContent = data.numActiveConn;
  numIpsPerSec.textContent = data.numIpsPerSec;
  inboundMBps.textContent = data.inboundMBps;

  const timestamp = new Date().toLocaleTimeString();
  chart.data.labels.push(timestamp);
  chart.data.datasets[0].data.push(data.numConnPerSec);
  chart.data.datasets[1].data.push(data.numActiveConn);
  chart.data.datasets[2].data.push(data.numIpsPerSec);
  chart.data.datasets[3].data.push(data.inboundMBps);

  if (chart.data.labels.length > limit) {
    chart.data.labels.shift();
    chart.data.datasets[0].data.shift();
    chart.data.datasets[1].data.shift();
    chart.data.datasets[2].data.shift();
    chart.data.datasets[3].data.shift();
  }

  chart.update();
});

