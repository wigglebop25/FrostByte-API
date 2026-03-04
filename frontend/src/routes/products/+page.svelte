<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { Product, CreateProductPayload } from "$lib/types";
    import { Plus, Search, X, Tag } from "lucide-svelte";
    import ProductCard from "$lib/components/ProductCard.svelte";

    let products: Product[] = [];
    let loading = true;
    let searchQuery = "";
    
    // RBAC
    let isAdmin = false;
    
    // Modal State
    let showModal = false;
    let isEditing = false;
    let modalError = "";
    let submitting = false;

    // Form State
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
        // Check Role
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
        // Only Admins can edit
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

<div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row gap-4 justify-between items-start md:items-center">
        <div>
            <h1 class="text-3xl font-bold tracking-tight text-accent font-serif">Products</h1>
            <p class="text-text-secondary mt-1">Manage your menu items and inventory.</p>
        </div>
        {#if isAdmin}
            <button on:click={openAddModal} class="flex items-center gap-2 bg-primary hover:bg-primary-dark px-4 py-2 rounded-xl font-bold transition-all text-text-inverse shadow-lg shadow-primary/20 active:scale-95">
                <Plus size={18} />
                <span>Add Product</span>
            </button>
        {/if}
    </div>

    <!-- Filters -->
    <div class="glass-card p-4">
        <div class="relative max-w-md">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary" size={18} />
            <input 
                bind:value={searchQuery}
                type="text" 
                placeholder="Search products..." 
                class="w-full bg-surface/50 border border-glass-border rounded-xl py-2.5 pl-10 pr-4 text-text-primary focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all"
            />
        </div>
    </div>

    <!-- Product Grid -->
    {#if loading}
        <div class="text-center py-12 text-text-secondary">Loading products...</div>
    {:else if filteredProducts.length === 0}
        <div class="text-center py-12 text-text-secondary">No products found.</div>
    {:else}
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {#each filteredProducts as product}
                <!-- Card is only clickable for Admins -->
                <button 
                    type="button"
                    class="h-full w-full text-left {isAdmin ? 'cursor-pointer' : 'cursor-default'}" 
                    on:click={() => handleCardClick(product)}
                    disabled={!isAdmin}
                >
                    <!-- We use the same card, but pass a no-op for the button if needed, or hide actions -->
                    <ProductCard product={product} onAdd={() => isAdmin ? handleCardClick(product) : null} />
                </button>
            {/each}
        </div>
    {/if}
</div>

<!-- Add/Edit Modal (Only rendered if showModal is true, which is guarded by isAdmin) -->
{#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm transition-all">
        <div class="glass-card w-full max-w-lg shadow-2xl overflow-hidden bg-surface relative">
            <div class="p-6 border-b border-glass-border flex justify-between items-center bg-surface/80">
                <h2 class="text-xl font-bold text-text-primary">{isEditing ? 'Edit Product' : 'Add New Product'}</h2>
                <button on:click={closeModal} class="text-text-secondary hover:text-primary transition-colors">
                    <X size={24} />
                </button>
            </div>
            
            <form on:submit|preventDefault={handleSubmit} class="p-6 space-y-4 bg-surface/50">
                {#if modalError}
                    <div class="p-3 bg-status-error/10 border border-status-error/20 text-status-error rounded-lg text-sm font-medium">
                        {modalError}
                    </div>
                {/if}

                <div class="grid grid-cols-2 gap-4">
                    <div class="col-span-2">
                        <label class="block text-sm font-bold text-text-primary mb-1" for="name">Product Name</label>
                        <input bind:value={formData.name} id="name" required class="w-full bg-main border border-glass-border rounded-xl px-4 py-2 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors" />
                    </div>

                    <div class="col-span-2 sm:col-span-1">
                        <label class="block text-sm font-bold text-text-primary mb-1" for="price">Price ($)</label>
                        <input bind:value={formData.price} id="price" type="number" step="0.01" required class="w-full bg-main border border-glass-border rounded-xl px-4 py-2 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors" />
                    </div>

                    <div class="col-span-2">
                        <label class="block text-sm font-bold text-text-primary mb-1" for="desc">Description</label>
                        <textarea bind:value={formData.description} id="desc" rows="3" class="w-full bg-main border border-glass-border rounded-xl px-4 py-2 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors"></textarea>
                    </div>

                    <div class="col-span-2">
                        <label class="block text-sm font-bold text-text-primary mb-1" for="img">Image URL</label>
                        <input bind:value={formData.product_image_uri} id="img" placeholder="https://example.com/image.png" class="w-full bg-main border border-glass-border rounded-xl px-4 py-2 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors" />
                    </div>

                    <div class="col-span-2">
                        <label class="block text-sm font-bold text-text-primary mb-1" for="cats">Categories</label>
                        <div class="relative">
                            <Tag class="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary" size={16} />
                            <input bind:value={categoryInput} id="cats" placeholder="Sushi, Drinks (comma separated)" class="w-full bg-main border border-glass-border rounded-xl pl-10 pr-4 py-2 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-colors" />
                        </div>
                    </div>
                </div>

                <div class="pt-6 flex gap-3">
                    {#if isEditing}
                        <button type="button" on:click={deleteProduct} class="px-4 py-2 rounded-xl border border-status-error/30 text-status-error hover:bg-status-error/10 transition-colors font-bold">Delete</button>
                    {/if}
                    <div class="flex-1"></div>
                    <button type="button" on:click={closeModal} class="px-4 py-2 rounded-xl border border-glass-border hover:bg-surface transition-colors font-bold text-text-secondary">Cancel</button>
                    <button type="submit" disabled={submitting} class="px-6 py-2 rounded-xl bg-primary hover:bg-primary-dark text-text-inverse font-bold transition-colors disabled:opacity-50 shadow-lg shadow-primary/20">
                        {submitting ? 'Saving...' : 'Save Product'}
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}