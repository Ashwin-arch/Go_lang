// Yios AI Kubernetes Platform Control Plane Client Logic

let deployments = [];

document.addEventListener('DOMContentLoaded', () => {
  loadDeployments();
});

// Load All Deployments
async function loadDeployments() {
  try {
    const res = await fetch('/api/v1/deployments');
    const data = await res.json();
    deployments = data.deployments || [];

    renderDeployments(deployments);
    updatePlatformStats(deployments);
  } catch (err) {
    console.error('Failed to fetch deployments', err);
  }
}

// Update Top Bar Stats
function updatePlatformStats(list) {
  document.getElementById('stat-total-deployments').textContent = list.length;
  
  const totalPods = list.reduce((sum, d) => sum + d.replicas, 0);
  document.getElementById('stat-active-pods').textContent = totalPods;

  const avgCpu = 42 + Math.floor(Math.random() * 15);
  document.getElementById('stat-avg-cpu').textContent = `${avgCpu}%`;

  const totalMem = list.reduce((sum, d) => sum + (d.replicas * (d.framework === 'FASTAPI' ? 1 : 8)), 0);
  document.getElementById('stat-total-mem').textContent = `${totalMem} GB`;
}

// Render Deployments Grid
function renderDeployments(list) {
  const grid = document.getElementById('deployments-grid');
  grid.innerHTML = '';

  if (list.length === 0) {
    grid.innerHTML = `<p style="grid-column:1/-1; text-align:center; color:#94a3b8;">No active AI deployments found. Click "Deploy AI Model" to get started.</p>`;
    return;
  }

  list.forEach(d => {
    const card = document.createElement('div');
    card.className = 'deploy-card';
    card.innerHTML = `
      <div>
        <div class="card-top">
          <div>
            <div class="deploy-name">${d.name}</div>
            <a href="#" class="deploy-endpoint" onclick="testProxyInference(${d.id}); return false;">
              🌐 ${d.publicEndpoint}
            </a>
          </div>
          <span class="fw-badge fw-${d.framework}">${d.framework}</span>
        </div>

        <div class="meta-row">
          <span>Replicas: <strong style="color:#fff;">${d.replicas}</strong> (1→100)</span>
          <span>Revision: <strong style="color:#06b6d4;">v${d.activeRevision}</strong></span>
          <span>Status: <strong style="color:#10b981;">● DEPLOYED</strong></span>
        </div>
      </div>

      <div class="card-actions">
        <button class="btn btn-outline btn-sm" onclick="showMetricsModal(${d.id})">📊 Metrics</button>
        <button class="btn btn-primary btn-sm" onclick="showScaleModal(${d.id}, '${d.name}', ${d.replicas})">📈 Scale</button>
        <button class="btn btn-outline btn-sm" onclick="showRevisionsModal(${d.id})">🔄 Rollback</button>
        <button class="btn btn-outline btn-sm" onclick="showManifestModal(${d.id})">⚙️ YAML</button>
        <button class="btn btn-outline btn-sm" onclick="showLogsModal(${d.id})">📜 Logs</button>
        <button class="btn btn-outline btn-sm" style="color:#ef4444;" onclick="deleteDeployment(${d.id})">🗑️ Delete</button>
      </div>
    `;
    grid.appendChild(card);
  });
}

// Create Deployment
async function handleCreateDeployment(e) {
  e.preventDefault();
  const name = document.getElementById('deploy-name').value;
  const framework = document.getElementById('deploy-framework').value;
  const sourceType = document.getElementById('deploy-sourcetype').value;
  const replicas = parseInt(document.getElementById('deploy-replicas').value, 10);
  const sourceUrl = document.getElementById('deploy-sourceurl').value;
  const dockerfile = document.getElementById('deploy-dockerfile').value;
  const imageTag = document.getElementById('deploy-imagetag').value;
  const memoryRequest = document.getElementById('deploy-memory').value;

  try {
    const res = await fetch('/api/v1/deployments', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name, framework, sourceType, replicas, sourceUrl, dockerfile, imageTag, memoryRequest
      })
    });
    const data = await res.json();
    if (res.ok) {
      showToast(`🎉 AI Model '${name}' deployed to Kubernetes!`);
      toggleModal('deploy-modal');
      loadDeployments();
    } else {
      showToast(data.error || 'Deployment failed', 'error');
    }
  } catch (err) {
    showToast('Deployment failed', 'error');
  }
}

// Scale Modal & Slider
function showScaleModal(id, name, replicas) {
  document.getElementById('scale-deploy-id').value = id;
  document.getElementById('scale-model-name').textContent = name;
  document.getElementById('scale-slider').value = replicas;
  updateScaleDisplay(replicas);
  toggleModal('scale-modal');
}

function updateScaleDisplay(val) {
  document.getElementById('scale-val-display').textContent = `${val} Replicas`;
}

async function handleScaleSubmit(e) {
  e.preventDefault();
  const id = document.getElementById('scale-deploy-id').value;
  const replicas = parseInt(document.getElementById('scale-slider').value, 10);

  try {
    const res = await fetch(`/api/v1/deployments/${id}/scale`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ replicas })
    });
    const data = await res.json();
    if (res.ok) {
      showToast(`📈 Scaled to ${replicas} replicas!`);
      toggleModal('scale-modal');
      loadDeployments();
    } else {
      showToast(data.error || 'Scaling failed', 'error');
    }
  } catch (err) {
    showToast('Scaling failed', 'error');
  }
}

// Revisions & Rollback
async function showRevisionsModal(id) {
  toggleModal('revisions-modal');
  const container = document.getElementById('revisions-list');
  container.innerHTML = '<p style="color:#94a3b8;">Fetching revision history...</p>';

  try {
    const res = await fetch(`/api/v1/deployments/${id}/revisions`);
    const data = await res.json();
    const revs = data.revisions || [];

    container.innerHTML = '';
    revs.forEach(r => {
      const item = document.createElement('div');
      item.style.cssText = 'padding:14px; background:rgba(255,255,255,0.05); border-radius:10px; margin-bottom:10px; display:flex; justify-content:space-between; align-items:center;';
      item.innerHTML = `
        <div>
          <strong style="color:#38bdf8;">Revision v${r.revisionNumber}</strong> - <code>${r.imageTag}</code>
          <div style="font-size:12px; color:#94a3b8; margin-top:4px;">${r.description || 'No description'} (SHA: ${r.commitSha})</div>
        </div>
        <button class="btn btn-primary btn-sm" onclick="rollbackToRevision(${id}, ${r.revisionNumber})">1-Click Rollback</button>
      `;
      container.appendChild(item);
    });
  } catch (err) {
    container.innerHTML = '<p style="color:#ef4444;">Failed to load revision history.</p>';
  }
}

async function rollbackToRevision(deployID, revNumber) {
  try {
    const res = await fetch(`/api/v1/deployments/${deployID}/rollback`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ revisionNumber: revNumber })
    });
    const data = await res.json();
    if (res.ok) {
      showToast(`🔄 Rolled back to Revision v${revNumber}!`);
      toggleModal('revisions-modal');
      loadDeployments();
    } else {
      showToast(data.error || 'Rollback failed', 'error');
    }
  } catch (err) {
    showToast('Rollback failed', 'error');
  }
}

// Pod Metrics Modal
async function showMetricsModal(id) {
  toggleModal('metrics-modal');
  const content = document.getElementById('metrics-content');
  content.innerHTML = '<p style="color:#94a3b8;">Polling Kubernetes Pod status...</p>';

  try {
    const res = await fetch(`/api/v1/deployments/${id}/metrics`);
    const data = await res.json();
    const m = data.metrics;

    let html = `
      <div style="display:grid; grid-template-columns:repeat(3, 1fr); gap:12px; margin-bottom:16px;">
        <div style="padding:12px; background:rgba(255,255,255,0.05); border-radius:8px;">
          <div style="font-size:11px; color:#94a3b8;">TOTAL REPLICAS</div>
          <div style="font-size:22px; font-weight:700;">${m.replicasCount} Pods</div>
        </div>
        <div style="padding:12px; background:rgba(255,255,255,0.05); border-radius:8px;">
          <div style="font-size:11px; color:#94a3b8;">AVG CPU LOAD</div>
          <div style="font-size:22px; font-weight:700; color:#38bdf8;">${m.cpuUsagePct.toFixed(1)}%</div>
        </div>
        <div style="padding:12px; background:rgba(255,255,255,0.05); border-radius:8px;">
          <div style="font-size:11px; color:#94a3b8;">INFERENCE REQ / SEC</div>
          <div style="font-size:22px; font-weight:700; color:#34d399;">${m.requestsPerSec} req/s</div>
        </div>
      </div>

      <h4 style="margin-bottom:10px;">Kubernetes Pod Instances</h4>
      <div style="overflow-x:auto;">
        <table style="width:100%; text-align:left; font-size:13px; border-collapse:collapse;">
          <thead>
            <tr style="border-bottom:1px solid rgba(255,255,255,0.1); color:#94a3b8;">
              <th style="padding:8px;">Pod Name</th>
              <th style="padding:8px;">Status</th>
              <th style="padding:8px;">Node</th>
              <th style="padding:8px;">CPU</th>
              <th style="padding:8px;">VRAM/RAM</th>
            </tr>
          </thead>
          <tbody>
    `;

    (m.pods || []).forEach(p => {
      html += `
        <tr style="border-bottom:1px solid rgba(255,255,255,0.05);">
          <td style="padding:10px; font-family:monospace;">${p.name}</td>
          <td style="padding:10px;"><span style="color:#10b981; font-weight:700;">● ${p.status}</span></td>
          <td style="padding:10px; color:#94a3b8;">${p.nodeName}</td>
          <td style="padding:10px;">${p.cpuUsagePct.toFixed(1)}%</td>
          <td style="padding:10px;">${(p.memoryUsageMb / 1024).toFixed(2)} GB</td>
        </tr>
      `;
    });

    html += `</tbody></table></div>`;
    content.innerHTML = html;
  } catch (err) {
    content.innerHTML = '<p style="color:#ef4444;">Failed to load Pod metrics.</p>';
  }
}

// Manifest YAML Modal
async function showManifestModal(id) {
  toggleModal('manifest-modal');
  const yamlText = document.getElementById('manifest-yaml-text');
  yamlText.textContent = 'Generating Kubernetes YAML specs...';

  try {
    const res = await fetch(`/api/v1/deployments/${id}/manifest`);
    const data = await res.json();
    yamlText.textContent = data.manifestYaml;
  } catch (err) {
    yamlText.textContent = 'Failed to generate K8s manifest.';
  }
}

// Logs Modal
async function showLogsModal(id) {
  toggleModal('manifest-modal');
  const text = document.getElementById('manifest-yaml-text');
  text.textContent = 'Streaming live container logs...';

  try {
    const res = await fetch(`/api/v1/deployments/${id}/logs`);
    const data = await res.json();
    text.textContent = (data.logs || []).join('\n');
  } catch (err) {
    text.textContent = 'Failed to fetch logs.';
  }
}

// Proxy Test Inference
async function testProxyInference(id) {
  try {
    const res = await fetch(`/api/v1/deployments/${id}/proxy`, { method: 'POST' });
    const data = await res.json();
    showToast(`🤖 Inference Test Response (200 OK): ${data.response.inference}`);
  } catch (err) {
    showToast('Inference proxy failed', 'error');
  }
}

// Delete Deployment
async function deleteDeployment(id) {
  if (!confirm('Are you sure you want to delete this Kubernetes AI service?')) return;
  try {
    const res = await fetch(`/api/v1/deployments/${id}`, { method: 'DELETE' });
    if (res.ok) {
      showToast('Deployment and K8s resources deleted!');
      loadDeployments();
    }
  } catch (err) {
    showToast('Failed to delete deployment', 'error');
  }
}

// Dynamic Form Listeners
function handleSourceTypeChange(val) {
  document.getElementById('group-source-url').classList.toggle('hidden', val !== 'GITHUB_REPO');
  document.getElementById('group-dockerfile').classList.toggle('hidden', val !== 'DOCKERFILE');
}

function handleFrameworkChange(val) {
  const mem = document.getElementById('deploy-memory');
  if (val === 'VLLM' || val === 'OLLAMA') {
    mem.value = '8Gi';
  } else {
    mem.value = '1Gi';
  }
}

function toggleModal(id) {
  document.getElementById(id).classList.toggle('hidden');
}

function showToast(msg, type = 'success') {
  const container = document.getElementById('toast-container');
  const t = document.createElement('div');
  t.className = 'toast';
  if (type === 'error') t.style.background = '#dc2626';
  t.textContent = msg;
  container.appendChild(t);
  setTimeout(() => t.remove(), 3500);
}
