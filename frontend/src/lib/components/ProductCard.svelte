<script lang="ts">
    /**
     * Product Display Card Component
     *
     * Renders a menu product with image, price badge, category pills,
     * and an add-to-cart action. Used across catalog and order pages.
     */
    import type { Product } from "$lib/types";
    import { ShoppingCart } from "lucide-svelte";
    import { getProductImageUrl, handleImageError } from "$lib/utils/image";
    
    export let product: Product;
    export let onAdd: (p: Product) => void = () => {};

    const formatPrice = (price: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(price);
    };
</script>

<div class="glass-card group h-full flex flex-col overflow-hidden hover:shadow-lg hover:-translate-y-1 transition-all duration-300">
    <!-- Product Image & Overlay Badges -->
    <div class="relative aspect-[4/3] overflow-hidden bg-[var(--color-bg-surface)]">
        <img 
            src={getProductImageUrl(product)} 
            alt={product.name} 
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
            on:error={handleImageError}
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent"></div>
        
        <!-- Price badge -->
        <div class="absolute top-3 right-3 px-3 py-1.5 rounded-xl bg-white/90 dark:bg-black/70 backdrop-blur-md text-sm font-bold text-[var(--color-primary)] shadow-lg">
            {formatPrice(product.price)}
        </div>

        <!-- Category pills -->
        {#if product.categories && product.categories.length > 0}
            <div class="absolute bottom-3 left-3 flex gap-1.5">
                {#each product.categories.slice(0, 2) as cat}
                    <span class="px-2.5 py-1 bg-white/20 backdrop-blur-md border border-white/30 rounded-lg text-[10px] font-bold text-white uppercase tracking-wider">
                        {typeof cat === 'string' ? cat : cat.name}
                    </span>
                {/each}
            </div>
        {/if}
    </div>

    <!-- Product Details & Add-to-Cart Action -->
    <div class="p-4 flex flex-col flex-1">
        <h3 class="font-bold text-base text-[var(--color-text-primary)] leading-tight line-clamp-1 mb-1">{product.name}</h3>
        <p class="text-[var(--color-text-secondary)] text-sm line-clamp-2 mb-4 flex-1 opacity-80">
            {product.description}
        </p>
        
        <button 
            on:click|stopPropagation={() => onAdd(product)}
            class="w-full py-2.5 bg-[var(--color-primary)] hover:bg-[var(--color-primary-dark)] text-[var(--color-text-inverse)] rounded-xl font-bold flex items-center justify-center gap-2 transition-all shadow-md hover:shadow-lg active:scale-[0.97] text-sm"
        >
            <ShoppingCart size={16} />
            <span>Add to Order</span>
        </button>
    </div>
</div>
