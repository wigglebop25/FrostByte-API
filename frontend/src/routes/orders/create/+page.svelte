<script lang="ts">
    /**
     * Customer Order Creation Page
     *
     * Provides a product browsing interface with cart management for
     * placing new orders. Restricted to the Customer role. Submits
     * orders to the /orders API endpoint.
     */
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import { goto } from "$app/navigation";
    import type { Product } from "$lib/types";
    import { ShoppingCart, Minus, Plus, ArrowLeft } from "lucide-svelte";
    import ProductCard from "$lib/components/ProductCard.svelte";
    import { getProductImageUrl, handleImageError } from "$lib/utils/image";

    let products: Product[] = [];
    let loading = true;
    let searchQuery = "";
    let submitting = false;

    type CartItem = Product & { quantity: number };
    let cart: CartItem[] = [];

    onMount(async () => {
        /** Enforces Customer-only access; redirects other roles to dashboard */
        const userStr = localStorage.getItem("user");
        if (userStr) {
            const user = JSON.parse(userStr);
            const roles = user.roles?.map((r: any) => r.name.toLowerCase()) || [];
            if (!roles.includes('customer')) {
                goto("/");
                return;
            }
        }

        try {
            const res = await api.get("/products");
            products = res.data;
        } catch (e) {
            console.error("Failed to load menu", e);
        } finally {
            loading = false;
        }
    });

    function addToCart(product: Product) {
        const existing = cart.find(item => item.product_id === product.product_id);
        if (existing) {
            existing.quantity++;
            cart = [...cart];
        } else {
            cart = [...cart, { ...product, quantity: 1 }];
        }
    }

    function updateQuantity(productId: number, change: number) {
        const index = cart.findIndex(item => item.product_id === productId);
        if (index === -1) return;

        const item = cart[index];
        const newQty = item.quantity + change;

        if (newQty <= 0) {
            cart.splice(index, 1);
        } else {
            item.quantity = newQty;
        }
        cart = [...cart];
    }

    async function placeOrder() {
        if (cart.length === 0) return;
        submitting = true;

        const payload = {
            products: cart.map(item => ({
                product_id: item.product_id,
                quantity: item.quantity
            }))
        };

        try {
            await api.post("/orders", payload);
            goto("/orders");
        } catch (e: any) {
            alert("Order failed: " + (e.response?.data?.error || "Unknown error"));
            submitting = false;
        }
    }

    $: filteredProducts = products.filter(p => p.name.toLowerCase().includes(searchQuery.toLowerCase()));
    $: cartTotal = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
</script>

<div class="h-[calc(100dvh-7rem)] sm:h-[calc(100dvh-8rem)] lg:h-[calc(100dvh-9rem)] flex flex-col lg:flex-row gap-3 lg:gap-4 overflow-hidden">
    <!-- Product Browsing Section -->
    <div class="flex-[3] lg:flex-1 flex flex-col min-w-0 min-h-0">
        <!-- Header -->
        <div class="flex items-center gap-3 mb-3 shrink-0">
            <a href="/orders" class="p-2 hover:bg-[var(--glass-bg)] rounded-xl transition-colors text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] border border-transparent hover:border-[var(--glass-border)]">
                <ArrowLeft size={18} />
            </a>
            <div class="flex-1 max-w-md">
                <input 
                    bind:value={searchQuery}
                    type="text" 
                    placeholder="Search menu..." 
                    class="glass-input w-full px-4 py-2 text-sm"
                />
            </div>
        </div>

        <!-- Product Grid -->
        <div class="flex-1 overflow-y-auto pr-1 scrollbar-thin min-h-0">
            {#if loading}
                <div class="text-center py-10 text-[var(--color-text-secondary)] text-sm">
                    <div class="w-7 h-7 border-2 border-[var(--color-primary)]/30 border-t-[var(--color-primary)] rounded-full animate-spin mx-auto mb-2"></div>
                    Loading menu...
                </div>
            {:else}
                <div class="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-2.5 sm:gap-3 pb-2">
                    {#each filteredProducts as product}
                        <ProductCard product={product} onAdd={addToCart} />
                    {/each}
                </div>
            {/if}
        </div>
    </div>

    <!-- Cart Summary & Checkout Panel -->
    <div class="w-full flex-[2] lg:flex-none lg:w-80 xl:w-96 glass-card-solid flex flex-col min-h-0 lg:h-full shadow-2xl">
        <div class="px-3 py-2.5 lg:px-4 lg:py-3 border-b border-[var(--glass-border-subtle)] bg-[var(--color-primary)]/5 shrink-0">
            <h2 class="text-sm lg:text-base font-bold flex items-center gap-2 text-[var(--color-text-primary)]">
                <ShoppingCart class="text-[var(--color-primary)]" size={18} />
                Current Order
                <span class="text-xs font-semibold text-[var(--color-text-secondary)] ml-auto">{cart.reduce((a, b) => a + b.quantity, 0)} items</span>
            </h2>
        </div>

        <!-- Cart Items (scrollable region between header and sticky footer) -->
        <div class="flex-1 overflow-y-auto p-2.5 lg:p-3 space-y-2 scrollbar-thin min-h-0">
            {#if cart.length === 0}
                <div class="h-full flex flex-col items-center justify-center text-[var(--color-text-secondary)] space-y-2 opacity-40 py-4">
                    <ShoppingCart size={28} />
                    <p class="font-semibold text-xs">Cart is empty</p>
                </div>
            {:else}
                {#each cart as item}
                    <div class="flex items-center gap-2.5 bg-[var(--color-primary)]/5 p-2 sm:p-2.5 rounded-xl border border-[var(--glass-border-subtle)] hover:border-[var(--color-primary)]/30 transition-all duration-200">
                        <img 
                            src={getProductImageUrl(item)} 
                            alt={item.name} 
                            class="w-9 h-9 sm:w-10 sm:h-10 rounded-lg object-cover bg-[var(--color-bg-surface)]" 
                            on:error={handleImageError}
                        />
                        <div class="flex-1 min-w-0">
                            <h4 class="font-bold truncate text-[var(--color-text-primary)] text-xs sm:text-sm">{item.name}</h4>
                            <p class="text-xs text-[var(--color-primary)] font-bold">${(item.price * item.quantity).toFixed(2)}</p>
                        </div>
                        <div class="flex items-center gap-1 bg-[var(--glass-bg)] rounded-lg p-0.5 border border-[var(--glass-border)]">
                            <button on:click={() => updateQuantity(item.product_id, -1)} class="p-1 hover:bg-[var(--color-bg-main)] rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">
                                <Minus size={12} />
                            </button>
                            <span class="w-5 text-center text-xs font-bold text-[var(--color-text-primary)]">{item.quantity}</span>
                            <button on:click={() => updateQuantity(item.product_id, 1)} class="p-1 hover:bg-[var(--color-bg-main)] rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">
                                <Plus size={12} />
                            </button>
                        </div>
                    </div>
                {/each}
            {/if}
        </div>

        <!-- Order Total & Submit — always visible at cart bottom -->
        <div class="px-3 py-2.5 lg:px-4 lg:py-3 border-t border-[var(--glass-border-subtle)] bg-[var(--color-bg-surface)]/80 backdrop-blur-sm shrink-0">
            <div class="flex justify-between items-center mb-2">
                <span class="text-[var(--color-text-secondary)] font-semibold text-xs">Total</span>
                <span class="text-lg lg:text-xl font-extrabold text-[var(--color-primary)]">${cartTotal.toFixed(2)}</span>
            </div>
            <button 
                on:click={placeOrder}
                disabled={cart.length === 0 || submitting}
                class="w-full glass-button py-2.5 lg:py-3 text-sm disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
                {#if submitting}
                    <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    Processing...
                {:else}
                    Confirm Order
                {/if}
            </button>
        </div>
    </div>
</div>