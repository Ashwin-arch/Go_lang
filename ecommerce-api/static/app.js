// Gin E-commerce Storefront Application Logic

let currentUserID = 1;
let currentUserRole = 'customer';
let allProducts = [];
let categories = [];

document.addEventListener('DOMContentLoaded', () => {
  loadCategories();
  loadProducts();
  updateCartBadge();
});

// Load Categories
async function loadCategories() {
  try {
    const res = await fetch('/api/v1/categories');
    const data = await res.json();
    categories = data.categories || [];

    const tabs = document.getElementById('category-tabs');
    const adminSelect = document.getElementById('admin-p-category');
    
    categories.forEach(cat => {
      const btn = document.createElement('button');
      btn.className = 'cat-btn';
      btn.textContent = cat.name;
      btn.onclick = (e) => filterCategory(cat.id, btn);
      tabs.appendChild(btn);

      if (adminSelect) {
        const opt = document.createElement('option');
        opt.value = cat.id;
        opt.textContent = cat.name;
        adminSelect.appendChild(opt);
      }
    });
  } catch (err) {
    console.error('Failed to load categories', err);
  }
}

// Load Products
async function loadProducts(catID = 'all', query = '') {
  try {
    let url = '/api/v1/products';
    const params = new URLSearchParams();
    if (catID !== 'all') params.append('categoryId', catID);
    if (query) params.append('q', query);
    if (params.toString()) url += '?' + params.toString();

    const res = await fetch(url);
    const data = await res.json();
    allProducts = data.products || [];
    renderProducts(allProducts);
  } catch (err) {
    console.error('Failed to load products', err);
  }
}

// Render Products Grid
function renderProducts(products) {
  const grid = document.getElementById('product-grid');
  grid.innerHTML = '';

  if (products.length === 0) {
    grid.innerHTML = `<p style="grid-column:1/-1; text-align:center; color:#94a3b8;">No products found.</p>`;
    return;
  }

  products.forEach(p => {
    const card = document.createElement('div');
    card.className = 'product-card';
    card.innerHTML = `
      <div>
        <span class="p-category">${p.category ? p.category.name : 'Item'}</span>
        <h3 class="p-title">${p.name}</h3>
        <p class="p-desc">${p.description}</p>
      </div>
      <div class="p-footer">
        <span class="p-price">$${p.price.toFixed(2)}</span>
        <button class="btn btn-primary" onclick="addToCart(${p.id})">+ Cart</button>
      </div>
    `;
    grid.appendChild(card);
  });
}

// Add Item to Cart
async function addToCart(productID) {
  try {
    const res = await fetch('/api/v1/cart/items', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-User-ID': currentUserID.toString()
      },
      body: JSON.stringify({ productId: productID, quantity: 1 })
    });
    const data = await res.json();
    if (res.ok) {
      showToast(data.message || 'Added to cart!');
      updateCartBadge();
    } else {
      showToast(data.error || 'Failed to add item', 'error');
    }
  } catch (err) {
    showToast('Network error', 'error');
  }
}

// Update Cart Badge
async function updateCartBadge() {
  try {
    const res = await fetch('/api/v1/cart', {
      headers: { 'X-User-ID': currentUserID.toString() }
    });
    const data = await res.json();
    const count = (data.cart || []).reduce((acc, item) => acc + item.quantity, 0);
    document.getElementById('cart-badge-count').textContent = count;
  } catch (err) {
    console.error(err);
  }
}

// Category Filter & Search
function filterCategory(catID, btn) {
  document.querySelectorAll('.cat-btn').forEach(b => b.classList.remove('active'));
  btn.classList.add('active');
  loadProducts(catID, document.getElementById('search-input').value);
}

let searchTimer;
function debounceSearch() {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    const activeBtn = document.querySelector('.cat-btn.active');
    const catID = activeBtn ? activeBtn.getAttribute('data-id') || 'all' : 'all';
    loadProducts(catID, document.getElementById('search-input').value);
  }, 300);
}

// Modals
function toggleCartModal() {
  const modal = document.getElementById('cart-modal');
  modal.classList.toggle('hidden');
  if (!modal.classList.contains('hidden')) {
    loadCartDetails();
  }
}

async function loadCartDetails() {
  try {
    const res = await fetch('/api/v1/cart', {
      headers: { 'X-User-ID': currentUserID.toString() }
    });
    const data = await res.json();
    const list = document.getElementById('cart-items-list');
    list.innerHTML = '';

    (data.cart || []).forEach(item => {
      const row = document.createElement('div');
      row.style.cssText = 'display:flex; justify-content:space-between; margin-bottom:12px; padding:10px; background:rgba(255,255,255,0.05); border-radius:8px;';
      row.innerHTML = `
        <div>
          <strong>${item.product ? item.product.name : 'Product'}</strong>
          <div style="font-size:12px; color:#94a3b8;">Qty: ${item.quantity} × $${item.product ? item.product.price : 0}</div>
        </div>
        <strong>$${((item.product ? item.product.price : 0) * item.quantity).toFixed(2)}</strong>
      `;
      list.appendChild(row);
    });

    document.getElementById('cart-subtotal').textContent = `$${(data.subtotal || 0).toFixed(2)}`;
  } catch (err) {
    console.error(err);
  }
}

// Checkout
async function handleCheckout() {
  try {
    const res = await fetch('/api/v1/orders/checkout', {
      method: 'POST',
      headers: { 'X-User-ID': currentUserID.toString() }
    });
    const data = await res.json();
    if (res.ok) {
      showToast('🎉 Checkout completed! Order placed.');
      toggleCartModal();
      updateCartBadge();
    } else {
      showToast(data.error || 'Checkout failed', 'error');
    }
  } catch (err) {
    showToast('Checkout failed', 'error');
  }
}

// Auth & Role Switcher
function toggleAuthModal() {
  document.getElementById('auth-modal').classList.toggle('hidden');
}

function selectRole(role) {
  currentUserRole = role;
  currentUserID = role === 'admin' ? 2 : 1;
  document.getElementById('role-customer').classList.toggle('active', role === 'customer');
  document.getElementById('role-admin').classList.toggle('active', role === 'admin');
}

function handleLoginSubmit(e) {
  e.preventDefault();
  document.getElementById('nav-user-name').textContent = currentUserRole === 'admin' ? 'Store Admin (ID: 2)' : 'Customer (ID: 1)';
  document.getElementById('admin-panel').classList.toggle('hidden', currentUserRole !== 'admin');
  toggleAuthModal();
  updateCartBadge();
  showToast(`Switched context to ${currentUserRole.toUpperCase()}`);
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
