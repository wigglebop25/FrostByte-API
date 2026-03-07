<script lang="ts">
    /**
     * Product Catalog Management Page
     *
     * Provides CRUD operations for restaurant menu products. Restricted
     * to Admin role. Supports product image URLs and category assignment.
     */
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { Product, CreateProductPayload } from "$lib/types";
    import { Plus, Search, X, Tag, Grid3x3 } from "lucide-svelte";
    import ProductCard from "$lib/components/ProductCard.svelte";

    let products: Product[] = [];
    let loading = true;
    let searchQuery = "";
    
    let isAdmin = false;
    
    let showModal = false;
    let isEditing = false;
    let modalError = "";
    let submitting = false;

    let formData: CreateProductPayload = {
        name: "",
        description: "",
        price: 0,
        product_image_uri: "",
        categories: []
    };
    let categoryInput = ""; 
    let editingId: number | null = null;

    async function fetchProducts() {
        loading = true;
        try {
            const res = await api.get("/products");
            products = res.data;
        } catch (e) {
            console.error("Failed to fetch products", e);
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        const userStr = localStorage.getItem("user");
        if (userStr) {
            const user = JSON.parse(userStr);
            const roles = user.roles?.map((r: any) => r.name) || [];
            isAdmin = roles.includes('Admin');
        }
        fetchProducts();
    });

    function openAddModal() {
        if (!isAdmin) return;
        isEditing = false;
        editingId = null;
        formData = { name: "", description: "", price: 0, product_image_uri: "", categories: [] };
        categoryInput = "";
        modalError = "";
        showModal = true;
    }

    function handleCardClick(product: Product) {
        if (!isAdmin) return;
        
        isEditing = true;
        editingId = product.product_id;
        formData = {
            name: product.name,
            description: product.description,
            price: product.price,
            product_image_uri: product.product_image_uri,
            categories: product.categories?.map(c => typeof c === 'string' ? c : c.name) || []
        };
        categoryInput = formData.categories.join(", ");
        modalError = "";
        showModal = true;
    }

    function closeModal() {
        showModal = false;
    }

    async function handleSubmit() {
        if (!isAdmin) return;
        submitting = true;
        modalError = "";
        formData.categories = categoryInput.split(",").map(s => s.trim()).filter(s => s.length > 0);

        try {
            if (isEditing && editingId) {
                await api.put(`/products/${editingId}`, formData);
            } else {
                await api.post("/products", formData);
            }
            showModal = false;
            fetchProducts();
        } catch (e: any) {
            console.error("Save error", e);
            modalError = e.response?.data?.error || "Failed to save product.";
        } finally {
            submitting = false;
        }
    }

    async function deleteProduct() {
        if (!isAdmin || !editingId || !confirm("Are you sure you want to delete this product?")) return;
        try {
            await api.delete(`/products/${editingId}`);
            closeModal();
            fetchProducts();
        } catch (e: any) {
            alert("Failed to delete: " + (e.response?.data?.error || "Unknown error"));
        }
    }

    $: filteredProducts = products.filter(p => 
        p.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
        p.description.toLowerCase().includes(searchQuery.toLowerCase())
    );
</script>

<div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div>
            <h1 class="text-2xl font-extrabold tracking-tight text-[var(--color-text-primary)]">Products</h1>
            <p class="text-[var(--color-text-secondary)] text-sm mt-0.5">Manage your menu items and inventory.</p>
        </div>
        {#if isAdmin}
            <button on:click={openAddModal} class="glass-button flex items-center gap-2 text-sm">
                <Plus size={16} />
                Add Product
            </button>
        {/if}
    </div>

    <!-- Search -->
    <div class="relative max-w-md">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-text-secondary)]" size={16} />
        <input 
            bind:value={searchQuery}
            type="text" 
            placeholder="Search products..." 
            class="glass-input w-full pl-12 pr-4 text-sm"
        />
    </div>

    <!-- Product Grid -->
    {#if loading}
        <div class="text-center py-16 text-[var(--color-text-secondary)] text-sm">
            <div class="w-8 h-8 border-2 border-[var(--color-primary)]/30 border-t-[var(--color-primary)] rounded-full animate-spin mx-auto mb-3"></div>
            Loading products...
        </div>
    {:else if filteredProducts.length === 0}
        <div class="text-center py-16 text-[var(--color-text-secondary)]">
            <Grid3x3 size={40} class="mx-auto mb-3 opacity-30" />
            <p class="text-sm">No products found.</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
            {#each filteredProducts as product}
                <button 
                    type="button"
                    class="h-full w-full text-left {isAdmin ? 'cursor-pointer' : 'cursor-default'}" 
                    on:click={() => handleCardClick(product)}
                    disabled={!isAdmin}
                >
                    <ProductCard product={product} onAdd={() => isAdmin ? handleCardClick(product) : null} />
                </button>
            {/each}
        </div>
    {/if}
</div>

<!-- Product Form Modal — Create / Edit Operations -->
{#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-md">
        <button type="button" class="absolute inset-0 w-full h-full cursor-default focus:outline-none" on:click={closeModal} aria-label="Close modal"></button>
        <div class="glass-card-solid w-full max-w-lg shadow-2xl overflow-hidden relative z-10">
            <div class="p-5 border-b border-[var(--glass-border-subtle)] flex justify-between items-center">
                <h2 class="text-lg font-bold text-[var(--color-text-primary)]">{isEditing ? 'Edit Product' : 'Add New Product'}</h2>
                <button on:click={closeModal} class="p-2 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] rounded-xl hover:bg-[var(--glass-bg)] transition-colors">
                    <X size={20} />
                </button>
            </div>
            
            <form on:submit|preventDefault={handleSubmit} class="p-5 space-y-4">
                {#if modalError}
                    <div class="p-3 bg-[var(--status-error)]/10 border border-[var(--status-error)]/20 text-[var(--status-error)] rounded-xl text-sm font-medium">
                        {modalError}
                    </div>
                {/if}

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div class="col-span-1 sm:col-span-2">
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="name">Product Name</label>
                        <input bind:value={formData.name} id="name" required class="glass-input w-full text-sm" />
                    </div>

                    <div class="col-span-1">
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="price">Price ($)</label>
                        <input bind:value={formData.price} id="price" type="number" step="0.01" required class="glass-input w-full text-sm" />
                    </div>

                    <div class="col-span-1 sm:col-span-2">
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="desc">Description</label>
                        <textarea bind:value={formData.description} id="desc" rows="3" class="glass-input w-full text-sm"></textarea>
                    </div>

                    <div class="col-span-1 sm:col-span-2">
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="img">Image URL</label>
                        <input bind:value={formData.product_image_uri} id="img" placeholder="https://example.com/image.png" class="glass-input w-full text-sm" />
                    </div>

                    <div class="col-span-1 sm:col-span-2">
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="cats">Categories</label>
                        <div class="relative">
                            <Tag class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-text-secondary)]" size={15} />
                            <input bind:value={categoryInput} id="cats" placeholder="Sushi, Drinks (comma separated)" class="glass-input w-full pl-12 text-sm" />
                        </div>
                    </div>
                </div>

                <div class="pt-3 flex gap-2.5">
                    {#if isEditing}
                        <button type="button" on:click={deleteProduct} class="px-4 py-2.5 rounded-xl border border-[var(--status-error)]/30 text-[var(--status-error)] hover:bg-[var(--status-error)]/10 transition-colors font-semibold text-sm">Delete</button>
                    {/if}
                    <div class="flex-1"></div>
                    <button type="button" on:click={closeModal} class="px-4 py-2.5 rounded-xl bg-[var(--glass-bg)] border border-[var(--glass-border)] text-[var(--color-text-secondary)] font-semibold text-sm hover:bg-[var(--color-bg-surface)] transition-colors">Cancel</button>
                    <button type="submit" disabled={submitting} class="glass-button text-sm disabled:opacity-50">
                        {submitting ? 'Saving...' : 'Save Product'}
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}