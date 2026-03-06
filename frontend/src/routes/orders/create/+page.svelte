<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import { goto } from "$app/navigation";
    import type { Product } from "$lib/types";
    import { Search, ShoppingCart, Minus, Plus, ArrowLeft } from "lucide-svelte";
    import ProductCard from "$lib/components/ProductCard.svelte";

    let products: Product[] = [];
    let loading = true;
    let searchQuery = "";
    let submitting = false;

    type CartItem = Product & { quantity: number };
    let cart: CartItem[] = [];

    onMount(async () => {
        // Access Control: Only Customers can create orders
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

    function getImageUrl(p: Product): string {
        return p.product_image_uri || "/images/itadaki_logo.png";
    }

    function handleImageError(e: Event) {
        const target = e.currentTarget as HTMLImageElement;
        target.src = "/images/itadaki_logo.png";
    }

    $: filteredProducts = products.filter(p => p.name.toLowerCase().includes(searchQuery.toLowerCase()));
    $: cartTotal = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
</script>

<div class="h-[calc(100vh-2rem)] flex gap-6 overflow-hidden">
    <!-- Left: Product Grid -->
    <div class="flex-1 flex flex-col min-w-0">
        <!-- Header -->
        <div class="flex items-center gap-4 mb-6">
            <a href="/orders" class="p-3 hover:bg-surface rounded-xl transition-colors text-text-secondary hover:text-primary border border-transparent hover:border-glass-border">
                <ArrowLeft size={24} />
            </a>
            <div class="relative flex-1 max-w-md">
                <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary" size={18} />
                <input 
                    bind:value={searchQuery}
                    type="text" 
                    placeholder="Search menu..." 
                    class="w-full bg-surface/50 border border-glass-border rounded-xl py-3 pl-10 pr-4 text-text-primary focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all shadow-sm"
                />
            </div>
        </div>

        <!-- Grid -->
        <div class="flex-1 overflow-y-auto pr-2">
            {#if loading}
                <div class="text-center py-12 text-text-secondary">Loading menu...</div>
            {:else}
                <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6 pb-20">
                    {#each filteredProducts as product}
                        <ProductCard product={product} onAdd={addToCart} />
                    {/each}
                </div>
            {/if}
        </div>
    </div>

    <!-- Right: Cart Summary -->
    <div class="w-96 glass-card flex flex-col shrink-0 h-full shadow-2xl bg-surface/80 backdrop-blur-xl border-l border-glass-border">
        <div class="p-6 border-b border-glass-border bg-white/50">
            <h2 class="text-xl font-bold flex items-center gap-2 text-text-primary">
                <ShoppingCart class="text-primary" />
                Current Order
            </h2>
            <p class="text-sm text-text-secondary mt-1 font-medium">
                {cart.reduce((a, b) => a + b.quantity, 0)} items selected
            </p>
        </div>

        <!-- Cart Items -->
        <div class="flex-1 overflow-y-auto p-4 space-y-3">
            {#if cart.length === 0}
                <div class="h-full flex flex-col items-center justify-center text-text-secondary space-y-4 opacity-50">
                    <ShoppingCart size={48} />
                    <p class="font-medium">Cart is empty</p>
                </div>
            {:else}
                {#each cart as item}
                    <div class="flex items-center gap-3 bg-main/50 p-3 rounded-xl border border-glass-border hover:border-primary/30 transition-colors">
                        <img 
                            src={getImageUrl(item)} 
                            alt={item.name} 
                            class="w-12 h-12 rounded-lg object-cover bg-surface" 
                            on:error={handleImageError}
                        />
                        <div class="flex-1 min-w-0">
                            <h4 class="font-bold truncate text-text-primary text-sm">{item.name}</h4>
                            <p class="text-xs text-primary font-bold">${(item.price * item.quantity).toFixed(2)}</p>
                        </div>
                        <div class="flex items-center gap-2 bg-surface rounded-lg p-1 border border-glass-border">
                            <button on:click={() => updateQuantity(item.product_id, -1)} class="p-1 hover:bg-main rounded text-text-secondary hover:text-primary transition-colors">
                                <Minus size={14} />
                            </button>
                            <span class="w-5 text-center text-sm font-bold text-text-primary">{item.quantity}</span>
                            <button on:click={() => updateQuantity(item.product_id, 1)} class="p-1 hover:bg-main rounded text-text-secondary hover:text-primary transition-colors">
                                <Plus size={14} />
                            </button>
                        </div>
                    </div>
                {/each}
            {/if}
        </div>

        <!-- Footer / Checkout -->
        <div class="p-6 border-t border-glass-border bg-surface">
            <div class="flex justify-between items-end mb-6">
                <span class="text-text-secondary font-medium">Total Amount</span>
                <span class="text-3xl font-bold text-primary">${cartTotal.toFixed(2)}</span>
            </div>
            <button 
                on:click={placeOrder}
                disabled={cart.length === 0 || submitting}
                class="w-full bg-primary hover:bg-primary-dark text-text-inverse font-bold py-4 rounded-xl transition-all shadow-lg shadow-primary/20 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 active:scale-95"
            >
                {#if submitting}
                    <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    Processing...
                {:else}
                    Confirm Order
                {/if}
            </button>
        </div>
    </div>
</div>