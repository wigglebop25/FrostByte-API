<script lang="ts">
    import type { Product } from "$lib/types";
    import { ShoppingCart, Tag } from "lucide-svelte";
    
    export let product: Product;
    export let onAdd: (p: Product) => void = () => {};

    const formatPrice = (price: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(price);
    };

    function getImageUrl(p: Product): string {
        const azureBase = "https://frostbytedata.blob.core.windows.net/products/";
        
        // 1. If we have a full URL (already Azure or other external), use it as is.
        if (p.product_image_uri && p.product_image_uri.startsWith('http')) {
            return p.product_image_uri;
        }

        // 2. If it's a local path (like '/images/pork-gyoza.jpg'), convert to Azure URL.
        if (p.product_image_uri && p.product_image_uri.includes('/')) {
            const filename = p.product_image_uri.split('/').pop()?.split('.')[0].replace(/-/g, '_') + '.png';
            return `${azureBase}${filename}`;
        }

        // 3. Fallback: Construct filename from product name if URI is missing or just a string.
        const nameToSlug = (p.product_image_uri || p.name).toLowerCase().replace(/[-\s]/g, '_').split('.')[0] + '.png';
        return `${azureBase}${nameToSlug}`;
    }

    function handleImageError(e: Event) {
        const target = e.currentTarget as HTMLImageElement;
        target.src = "/images/itadaki_logo.png";
    }
</script>

<div class="glass-card group h-full flex flex-col overflow-hidden hover:scale-[1.02] transition-all duration-300">
    <div class="relative aspect-video overflow-hidden bg-surface/50">
        <img 
            src={getImageUrl(product)} 
            alt={product.name} 
            class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
            on:error={handleImageError}
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"></div>
        <div class="absolute bottom-3 left-3 flex gap-2">
            {#if product.categories}
                {#each product.categories.slice(0, 1) as cat}
                    <span class="px-2 py-1 bg-white/20 backdrop-blur-md border border-white/30 rounded text-[10px] font-bold text-white uppercase tracking-wider">
                        {typeof cat === 'string' ? cat : cat.name}
                    </span>
                {/each}
            {/if}
        </div>
    </div>

    <div class="p-5 flex flex-col flex-1">
        <div class="flex justify-between items-start mb-2">
            <h3 class="font-bold text-lg text-text-primary leading-tight line-clamp-1">{product.name}</h3>
            <span class="text-primary font-bold">{formatPrice(product.price)}</span>
        </div>
        
        <p class="text-text-secondary text-sm line-clamp-2 mb-6 flex-1 italic">
            {product.description}
        </p>
        
        <button 
            on:click={() => onAdd(product)}
            class="w-full py-3 bg-primary hover:bg-primary-dark text-text-inverse rounded-xl font-bold flex items-center justify-center gap-2 transition-all shadow-lg shadow-primary/20 active:scale-95"
        >
            <ShoppingCart size={18} />
            <span>Add to Order</span>
        </button>
    </div>
</div>
