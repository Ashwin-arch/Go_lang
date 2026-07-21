// Echo Banking Web Dashboard Application Logic

let accounts = [];
let activeAccountID = 1;

document.addEventListener('DOMContentLoaded', () => {
  loadAccounts();
});

// Load Accounts
async function loadAccounts() {
  try {
    const res = await fetch('/api/v1/accounts');
    const data = await res.json();
    accounts = data.accounts || [];

    renderAccountSelector();
    renderActiveAccount();
    loadAccountTransactions();
  } catch (err) {
    console.error('Failed to load accounts', err);
  }
}

// Render Account Selector
function renderAccountSelector() {
  const sel = document.getElementById('account-selector');
  const transferSel = document.getElementById('transfer-to-account');

  sel.innerHTML = '';
  transferSel.innerHTML = '';

  accounts.forEach(acc => {
    const opt = document.createElement('option');
    opt.value = acc.id;
    opt.textContent = `${acc.accountNumber} (${acc.accountHolder})`;
    if (acc.id === activeAccountID) opt.selected = true;
    sel.appendChild(opt);

    if (acc.id !== activeAccountID) {
      const opt2 = document.createElement('option');
      opt2.value = acc.id;
      opt2.textContent = `${acc.accountNumber} - ${acc.accountHolder} ($${acc.balance.toFixed(2)})`;
      transferSel.appendChild(opt2);
    }
  });
}

// Render Active Account Balance & Details
function renderActiveAccount() {
  const acc = accounts.find(a => a.id === activeAccountID);
  if (!acc) return;

  document.getElementById('account-balance').textContent = `$${acc.balance.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
  document.getElementById('account-num').textContent = acc.accountNumber;
  document.getElementById('account-holder').textContent = acc.accountHolder;
}

// Switch Active Account Context
function handleAccountSwitch(id) {
  activeAccountID = parseInt(id, 10);
  renderAccountSelector();
  renderActiveAccount();
  loadAccountTransactions();
  showToast(`Switched active context to Account #${activeAccountID}`);
}

// Load Transactions
async function loadAccountTransactions() {
  try {
    const res = await fetch(`/api/v1/accounts/${activeAccountID}/transactions`);
    const data = await res.json();
    const rows = document.getElementById('transaction-rows');
    rows.innerHTML = '';

    const txs = data.transactions || [];
    if (txs.length === 0) {
      rows.innerHTML = `<tr><td colspan="5" style="text-align:center; color:#94a3b8;">No transaction history found.</td></tr>`;
      return;
    }

    txs.forEach(tx => {
      const tr = document.createElement('tr');
      const date = new Date(tx.createdAt).toLocaleString();
      const isCredit = tx.toAccountId === activeAccountID || tx.type === 'DEPOSIT';
      const amtClass = isCredit ? 'style="color:#34d399; font-weight:700;"' : 'style="color:#ef4444; font-weight:700;"';
      const sign = isCredit ? '+' : '-';

      tr.innerHTML = `
        <td><code>${tx.referenceNo}</code></td>
        <td><span class="tx-badge tx-${tx.type}">${tx.type}</span></td>
        <td>${tx.description || 'N/A'}</td>
        <td style="color:#94a3b8; font-size:13px;">${date}</td>
        <td ${amtClass}>${sign}$${tx.amount.toFixed(2)}</td>
      `;
      rows.appendChild(tr);
    });
  } catch (err) {
    console.error(err);
  }
}

// Deposit
async function handleDepositSubmit(e) {
  e.preventDefault();
  const amt = parseFloat(document.getElementById('deposit-amount').value);
  const desc = document.getElementById('deposit-desc').value;

  try {
    const res = await fetch(`/api/v1/accounts/${activeAccountID}/deposit`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ amount: amt, description: desc })
    });
    const data = await res.json();

    if (res.ok) {
      showToast('Deposit successful!');
      toggleModal('deposit-modal');
      loadAccounts();
    } else {
      showToast(data.error || 'Deposit failed', 'error');
    }
  } catch (err) {
    showToast('Deposit failed', 'error');
  }
}

// Withdraw
async function handleWithdrawSubmit(e) {
  e.preventDefault();
  const amt = parseFloat(document.getElementById('withdraw-amount').value);
  const desc = document.getElementById('withdraw-desc').value;

  try {
    const res = await fetch(`/api/v1/accounts/${activeAccountID}/withdraw`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ amount: amt, description: desc })
    });
    const data = await res.json();

    if (res.ok) {
      showToast('Withdrawal successful!');
      toggleModal('withdraw-modal');
      loadAccounts();
    } else {
      showToast(data.error || 'Withdrawal failed', 'error');
    }
  } catch (err) {
    showToast('Withdrawal failed', 'error');
  }
}

// Transfer
async function handleTransferSubmit(e) {
  e.preventDefault();
  const toAccID = parseInt(document.getElementById('transfer-to-account').value, 10);
  const amt = parseFloat(document.getElementById('transfer-amount').value);
  const desc = document.getElementById('transfer-desc').value;

  try {
    const res = await fetch('/api/v1/transfer', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        fromAccountId: activeAccountID,
        toAccountId: toAccID,
        amount: amt,
        description: desc
      })
    });
    const data = await res.json();

    if (res.ok) {
      showToast('🎉 ACID Fund Transfer Completed!');
      toggleModal('transfer-modal');
      loadAccounts();
    } else {
      showToast(data.error || 'Transfer failed', 'error');
    }
  } catch (err) {
    showToast('Transfer failed', 'error');
  }
}

// Modal Toggle
function toggleModal(modalID) {
  document.getElementById(modalID).classList.toggle('hidden');
}

// Toast
function showToast(msg, type = 'success') {
  const container = document.getElementById('toast-container');
  const t = document.createElement('div');
  t.className = 'toast';
  if (type === 'error') t.style.background = '#dc2626';
  t.textContent = msg;
  container.appendChild(t);
  setTimeout(() => t.remove(), 3000);
}
