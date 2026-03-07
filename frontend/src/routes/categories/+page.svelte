<script lang="ts">
    /**
     * Category Management Page
     *
     * Provides CRUD operations for menu product categories. Admin-only
     * access. Categories are used to organize products in the catalog.
     */
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { Category } from "$lib/types";
    import { Plus, Search, Trash2, Edit, X, Tag, FolderOpen } from "lucide-svelte";

    let categories: Category[] = [];
    let loading = true;
    let searchQuery = "";
    let error = "";

    let showModal = false;
    let isEditing = false;
    let modalError = "";
    let submitting = false;

    let formData = { name: "", description: "" };
    let editingId: number | null = null;

    async function fetchCategories() {
        loading = true;
        try {
            const res = await api.get("/categories");
            categories = res.data;
        } catch (e: any) {
            console.error("Failed to fetch categories", e);
            error = "Failed to load categories.";
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        fetchCategories();
    });

    function openAddModal() {
        isEditing = false;
        editingId = null;
        formData = { name: "", description: "" };
        modalError = "";
        showModal = true;
    }

    function openEditModal(cat: Category) {
        isEditing = true;
        editingId = cat.category_id;
        formData = { name: cat.name, description: cat.description };
        modalError = "";
        showModal = true;
    }

    function closeModal() {
        showModal = false;
    }

    async function handleSubmit() {
        submitting = true;
        modalError = "";
        try {
            if (isEditing && editingId) {
                await api.put(`/categories/${editingId}`, formData);
            } else {
                await api.post("/categories", formData);
            }
            showModal = false;
            fetchCategories();
        } catch (e: any) {
            console.error("Save error", e);
            modalError = e.response?.data?.error || "Failed to save category.";
        } finally {
            submitting = false;
        }
    }

    async function deleteCategory(id: number) {
        if (!confirm("Are you sure?")) return;
        try {
            await api.delete(`/categories/${id}`);
            fetchCategories();
        } catch (e: any) {
            alert("Delete failed: " + (e.response?.data?.error || "Unknown error"));
        }
    }

    $: filteredCategories = categories.filter(c => c.name.toLowerCase().includes(searchQuery.toLowerCase()));
</script>

<div class="space-y-5">
    <div class="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div>
            <h1 class="text-2xl font-extrabold tracking-tight text-[var(--color-text-primary)]">Categories</h1>
            <p class="text-[var(--color-text-secondary)] text-sm mt-0.5">Organize your menu items into groups.</p>
        </div>
        <button on:click={openAddModal} class="glass-button flex items-center gap-2 text-sm">
            <Plus size={16} />
            Add Category
        </button>
    </div>

    <!-- Search -->
    <div class="relative max-w-md">
        <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--color-text-secondary)]" size={16} />
        <input bind:value={searchQuery} placeholder="Search categories..." class="glass-input w-full pl-10 pr-4 text-sm" />
    </div>

    {#if loading}
        <div class="text-center py-16 text-[var(--color-text-secondary)] text-sm">
            <div class="w-8 h-8 border-2 border-[var(--color-primary)]/30 border-t-[var(--color-primary)] rounded-full animate-spin mx-auto mb-3"></div>
            Loading...
        </div>
    {:else if filteredCategories.length === 0}
        <div class="text-center py-16 text-[var(--color-text-secondary)]">
            <FolderOpen size={40} class="mx-auto mb-3 opacity-30" />
            <p class="text-sm">No categories found.</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each filteredCategories as cat}
                <div class="glass-card p-5 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 group">
                    <div class="flex justify-between items-start mb-3">
                        <div class="w-10 h-10 rounded-xl bg-[var(--color-primary)]/10 flex items-center justify-center text-[var(--color-primary)]">
                            <Tag size={20} />
                        </div>
                        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button on:click={() => openEditModal(cat)} class="p-2 hover:bg-[var(--color-primary)]/10 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] rounded-lg transition-colors">
                                <Edit size={16} />
                            </button>
                            <button on:click={() => deleteCategory(cat.category_id)} class="p-2 hover:bg-[var(--status-error)]/10 text-[var(--color-text-secondary)] hover:text-[var(--status-error)] rounded-lg transition-colors">
                                <Trash2 size={16} />
                            </button>
                        </div>
                    </div>
                    <h3 class="text-base font-bold text-[var(--color-text-primary)] mb-1">{cat.name}</h3>
                    <p class="text-[var(--color-text-secondary)] text-sm opacity-80">{cat.description || "No description."}</p>
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- Category Form Modal — Create / Edit Operations -->
{#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-md" role="dialog" aria-modal="true" tabindex="-1">
        <button type="button" class="absolute inset-0 w-full h-full cursor-default focus:outline-none" on:click={closeModal} aria-label="Close modal"></button>
        <div class="glass-card-solid w-full max-w-md shadow-2xl overflow-hidden relative z-10">
            <div class="p-5 border-b border-[var(--glass-border-subtle)] flex justify-between items-center">
                <h2 class="text-lg font-bold text-[var(--color-text-primary)]">{isEditing ? 'Edit Category' : 'New Category'}</h2>
                <button on:click={closeModal} class="p-2 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] rounded-xl hover:bg-[var(--glass-bg)] transition-colors"><X size={20} /></button>
            </div>
            <form on:submit|preventDefault={handleSubmit} class="p-5 space-y-4">
                {#if modalError}<div class="p-3 bg-[var(--status-error)]/10 border border-[var(--status-error)]/20 text-[var(--status-error)] rounded-xl text-sm font-medium">{modalError}</div>{/if}
                <div class="space-y-3">
                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="name">Name</label>
                        <input bind:value={formData.name} id="name" required class="glass-input w-full text-sm" />
                    </div>
                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="desc">Description</label>
                        <textarea bind:value={formData.description} id="desc" rows="3" class="glass-input w-full text-sm"></textarea>
                    </div>
                </div>
                <div class="pt-2 flex gap-2.5">
                    <button type="button" on:click={closeModal} class="flex-1 py-2.5 rounded-xl bg-[var(--glass-bg)] border border-[var(--glass-border)] text-[var(--color-text-secondary)] font-semibold text-sm hover:bg-[var(--color-bg-surface)] transition-colors">Cancel</button>
                    <button type="submit" disabled={submitting} class="flex-1 glass-button text-sm disabled:opacity-50">Save Category</button>
                </div>
            </form>
        </div>
    </div>
{/if}