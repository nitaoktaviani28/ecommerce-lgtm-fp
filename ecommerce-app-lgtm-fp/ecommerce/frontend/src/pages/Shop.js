import React, { useEffect, useState } from 'react';
import { trackPageView, trackAddToCart, pushLog, trackError } from '../lib/faro';

const CATEGORY_EMOJI = {
  electronics: '💻',
  fashion: '👕',
  books: '📚',
  food: '☕',
  outdoor: '🏕️',
};

export default function Shop({ cart, setCart }) {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('Semua');
  const [addedId, setAddedId] = useState(null);

  useEffect(() => {
    trackPageView('shop');
    pushLog('Shop page loaded', 'INFO');

    fetch('/api/products')
      .then((r) => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`);
        return r.json();
      })
      .then((data) => {
        setProducts(data);
        pushLog('Products loaded', 'INFO', { count: String(data.length) });
      })
      .catch((err) => {
        trackError(err.message, { source: 'fetch_products' });
        pushLog('Failed to load products', 'ERROR', { error: err.message });
      })
      .finally(() => setLoading(false));
  }, []);

  const categories = ['Semua', ...new Set(products.map((p) => p.category))];
  const filtered =
    filter === 'Semua' ? products : products.filter((p) => p.category === filter);

  const addToCart = (product) => {
    setCart((prev) => {
      const exists = prev.find((i) => i.id === product.id);
      if (exists) return prev.map((i) => i.id === product.id ? { ...i, qty: i.qty + 1 } : i);
      return [...prev, { ...product, qty: 1 }];
    });
    trackAddToCart(product);
    pushLog(`Added to cart: ${product.name}`, 'INFO', { product_id: String(product.id) });
    setAddedId(product.id);
    setTimeout(() => setAddedId(null), 1000);
  };

  if (loading) {
    return (
      <div className="shop-page">
        <div className="loading-state">
          <div className="spinner">⏳</div>
          <p>Memuat produk...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="shop-page">
      <div className="shop-header">
        <h1>Katalog Produk</h1>
        <p className="subtitle">Ecommerce FP — setiap interaksi dikirim ke LGTM Stack via Faro</p>
      </div>

      <div className="category-filter">
        {categories.map((c) => (
          <button
            key={c}
            className={filter === c ? 'cat-btn active' : 'cat-btn'}
            onClick={() => setFilter(c)}
          >
            {CATEGORY_EMOJI[c] || '🛍️'} {c}
          </button>
        ))}
      </div>

      <div className="product-grid">
        {filtered.map((p) => (
          <div key={p.id} className="product-card">
            <div className="product-emoji">{CATEGORY_EMOJI[p.category] || '🛍️'}</div>
            <div className="product-info">
              <h3>{p.name}</h3>
              <p className="product-desc">{p.description}</p>
              <span className="product-category">{p.category}</span>
            </div>
            <div className="product-footer">
              <div className="product-price">Rp {p.price.toLocaleString('id-ID')}</div>
              <div className="product-stock">Stok: {p.stock}</div>
              <button
                className={`btn-add ${addedId === p.id ? 'added' : ''}`}
                onClick={() => addToCart(p)}
                disabled={p.stock === 0}
              >
                {addedId === p.id ? '✓ Ditambahkan' : p.stock === 0 ? 'Habis' : '+ Keranjang'}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
