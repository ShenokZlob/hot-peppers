let offset = 0;
let isLoading = false;
let allData = [];
const API_BASE = window.API_BASE || "https://hot-peppers.onrender.com";

function fetchPeppers() {
  if (isLoading) return;
  isLoading = true;

  fetch(`${API_BASE}/scroll?offset=${offset}`)
    .then(res => res.json())
    .then(data => {
      allData.push(...data);
      renderTable();
      offset += data.length;
    })
    .catch(err => console.error('Ошибка загрузки:', err))
    .finally(() => isLoading = false);
}

function renderTable() {
  const tbody = document.querySelector('#pepperTable tbody');
  if (!tbody) return;

  tbody.innerHTML = ''; // очистка перед вставкой
  const query = document.getElementById('searchInput')?.value.toLowerCase() || '';
  const minMid = parseInt(document.getElementById('shuMidFilter')?.value) || 0;
  const maxUp = parseInt(document.getElementById('shuUpFilter')?.value) || Infinity;

  allData.forEach(item => {
  const name = item.name.toLowerCase();
  if (name.includes(query) && item.shu_mid >= minMid && item.shu_up <= maxUp) {
    const row = document.createElement('tr');
    row.innerHTML = `
      <td>
        <a 
          href="${API_BASE}/pepper/${encodeURIComponent(item.name)}" 
          class="text-decoration-underline link-primary pepper-link" 
          target="_blank"
        >
          ${item.name}
        </a>
      </td>
      <td>${item.shu_low}</td>
      <td>${item.shu_mid}</td>
      <td>${item.shu_up}</td>
    `;
    tbody.appendChild(row);
  }
});
}

// document.querySelectorAll('.pepper-link').forEach(btn => {
//   btn.addEventListener('click', function () {
//     const name = this.getAttribute('data-name');
//     fetch(`${API_BASE}/pepper/${name}`)
//       .then(response => response.text())
//       .then(html => {
//         const content = document.getElementById('content');
//         content.innerHTML = `
//           <div>
//             <button id="back-btn" class="btn btn-outline-light mb-3">← Назад</button>
//             <div id="pepper-detail">${html}</div>
//           </div>
//         `;

//         document.getElementById('back-btn').addEventListener('click', function () {
//           location.reload(); // или вручную восстанови исходный вид
//         });
//       })
//       .catch(() => alert('Ошибка загрузки перца'));
//   });
// });

function initTable() {
  offset = 0;
  isLoading = false;
  allData = [];

  fetchPeppers();

  const container = document.getElementById('tableContainer');
  container?.addEventListener('scroll', () => {
    if (container.scrollTop + container.clientHeight >= container.scrollHeight - 20) {
      fetchPeppers();
    }
  });

  ['searchInput', 'shuMidFilter', 'shuUpFilter'].forEach(id => {
    document.getElementById(id)?.addEventListener('input', renderTable);
  });
}

window.addEventListener('DOMContentLoaded', initTable);
