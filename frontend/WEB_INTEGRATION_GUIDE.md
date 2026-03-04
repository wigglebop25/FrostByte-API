# Web Integration Guide: Svelte Implementation (Dark Mode Supported)

This guide demonstrates how to combine the **Product Data** and **Color Palette** into a functional Svelte component that automatically adapts to Light and Dark modes.

## 1. Global Styles (app.css)

Ensure your CSS variables are set to handle both themes.

```css
:root {
 /* --- Functional Colors (Shared) --- */
 --status-success: #10B981;
 --status-warning: #F59E0B;
 --status-info: #3B82F6;
 --status-error: #EF4444;

 /* --- Light Mode (Default) --- */
 --color-primary: #E63946;
 --color-primary-dark: #D32F2F;
 --color-accent: #9B1C1C;
 
 --color-bg-main: #FFF8E1;
 --color-bg-surface: #FFFFFF;
 
 --color-text-primary: #1D3557;
 --color-text-secondary: #455A64;
 --color-text-inverse: #F1FAEE;

 /* Glassmorphism - Light */
 --glass-bg: rgba(255, 255, 255, 0.2);
 --glass-border: 1px solid rgba(255, 255, 255, 0.27);
 --glass-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
 --glass-blur: blur(20px);
}

/* --- Dark Mode Overrides --- */
@media (prefers-color-scheme: dark) {
 :root {
 /* Brand Colors (Adjusted for contrast) */
 --color-primary: #FF8A80;
 --color-primary-dark: #E57373;
 --color-accent: #FFCDD2;

 /* Backgrounds */
 --color-bg-main: #121212;
 --color-bg-surface: #1E1E1E;

 /* Text */
 --color-text-primary: #E6E1E5;
 --color-text-secondary: #CAC4D0;
 --color-text-inverse: #1C1B1F;

 /* Glassmorphism - Dark */
 --glass-bg: rgba(30, 30, 30, 0.6);
 --glass-border: 1px solid rgba(255, 255, 255, 0.1);
 --glass-shadow: 0 4px 30px rgba(0, 0, 0, 0.5);
 }
}

body {
 background-color: var(--color-bg-main);
 color: var(--color-text-primary);
 transition: background-color 0.3s ease, color 0.3s ease;
}
```

## 2. Product Store (stores/products.js)

A simple Svelte store to fetch data.

```javascript
import { writable } from 'svelte/store';

export const products = writable([]);
export const loading = writable(false);
export const error = writable(null);

const API_BASE = 'https://frostbyte-api.southeastasia.cloudapp.azure.com/api/v1';

export const fetchProducts = async () => {
 loading.set(true);
 try {
 const response = await fetch(`${API_BASE}/products`);
 if (!response.ok) throw new Error('Failed to fetch products');
 const data = await response.json();
 
 // Map API image paths
 const mappedData = data.map(p => ({
 ...p,
 product_image_uri: p.product_image_uri.startsWith('/') 
 ? `/images${p.product_image_uri}` 
 : `/images/${p.product_image_uri}`
 }));
 
 products.set(mappedData);
 } catch (e) {
 error.set(e.message);
 } finally {
 loading.set(false);
 }
};
```

## 3. Glass Product Card Component (components/ProductCard.svelte)

A reusable component implementing the glassmorphism design.

```svelte
<script>
 export let product;
 
 const formatPrice = (price) => {
 return new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD'
 }).format(price);
 };
</script>

<div class="glass-card">
 <div class="image-container">
 <img src={product.product_image_uri} alt={product.name} loading="lazy" />
 </div>
 
 <div class="content">
 <div class="header">
 <h3>{product.name}</h3>
 <span class="price">{formatPrice(product.price)}</span>
 </div>
 
 <p class="description">{product.description}</p>
 
 <div class="actions">
 <button class="add-btn">Add to Order</button>
 </div>
 </div>
</div>

<style>
 .glass-card {
 background: var(--glass-bg);
 border: var(--glass-border);
 box-shadow: var(--glass-shadow);
 backdrop-filter: var(--glass-blur);
 -webkit-backdrop-filter: var(--glass-blur);
 border-radius: 16px;
 overflow: hidden;
 transition: transform 0.2s ease, background 0.3s ease;
 display: flex;
 flex-direction: column;
 height: 100%;
 }

 .glass-card:hover {
 transform: translateY(-4px);
 }

 .image-container {
 width: 100%;
 height: 180px;
 overflow: hidden;
 background-color: var(--color-bg-surface); /* Placeholder bg */
 }

 img {
 width: 100%;
 height: 100%;
 object-fit: cover;
 }

 .content {
 padding: 1.25rem;
 display: flex;
 flex-direction: column;
 flex-grow: 1;
 }

 .header {
 display: flex;
 justify-content: space-between;
 align-items: flex-start;
 margin-bottom: 0.5rem;
 }

 h3 {
 margin: 0;
 font-size: 1.1rem;
 color: var(--color-text-primary);
 font-weight: 600;
 }

 .price {
 font-weight: 700;
 color: var(--color-primary);
 background: rgba(var(--color-primary), 0.1);
 padding: 4px 8px;
 border-radius: 12px;
 font-size: 0.9rem;
 /* Fallback for opacity if not using rgb vars */
 border: 1px solid var(--color-primary); 
 background: transparent;
 }

 .description {
 font-size: 0.9rem;
 color: var(--color-text-secondary);
 line-height: 1.4;
 margin-bottom: 1.5rem;
 flex-grow: 1;
 }

 .add-btn {
 width: 100%;
 padding: 10px;
 border: none;
 border-radius: 8px;
 background-color: var(--color-primary);
 color: var(--color-text-inverse);
 font-weight: 600;
 cursor: pointer;
 transition: background-color 0.2s;
 }

 .add-btn:hover {
 background-color: var(--color-primary-dark);
 }
</style>
```

## 4. Usage in a Page (routes/+page.svelte)

```svelte
<script>
 import { onMount } from 'svelte';
 import { products, fetchProducts, loading, error } from '../stores/products';
 import ProductCard from '../components/ProductCard.svelte';

 onMount(() => {
 fetchProducts();
 });
</script>

<main>
 <h1>Our Menu</h1>
 
 {#if $loading}
 <div class="loading">Loading delicious items...</div>
 {:else if $error}
 <div class="error">Error: {$error}</div>
 {:else}
 <div class="grid">
 {#each $products as product (product.id)}
 <ProductCard {product} />
 {/each}
 </div>
 {/if}
</main>

<style>
 main {
 max-width: 1200px;
 margin: 0 auto;
 padding: 2rem;
 }
 
 h1 {
 color: var(--color-text-primary);
 margin-bottom: 2rem;
 text-align: center;
 }

 .grid {
 display: grid;
 grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
 gap: 2rem;
 }
 
 .loading, .error {
 text-align: center;
 color: var(--color-text-secondary);
 font-size: 1.2rem;
 margin-top: 3rem;
 }
</style>
```
