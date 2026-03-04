<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { Category } from "$lib/types";
    import { Plus, Search, Trash2, Edit, X, Tag } from "lucide-svelte";

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

<div class="space-y-6">
    <div class="flex flex-col md:flex-row gap-4 justify-between items-start md:items-center">
        <div>
            <h1 class="text-3xl font-bold tracking-tight text-accent font-serif">Categories</h1>
            <p class="text-text-secondary mt-1">Organize your menu items.</p>
        </div>
        <button on:click={openAddModal} class="flex items-center gap-2 bg-primary hover:bg-primary-dark px-5 py-2.5 rounded-xl font-bold transition-all text-text-inverse shadow-lg shadow-primary/20 active:scale-95">
            <Plus size={18} />
            <span>Add Category</span>
        </button>
    </div>

    <div class="glass-card p-4">
        <div class="relative max-w-md">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary" size={18} />
            <input bind:value={searchQuery} placeholder="Search categories..." class="w-full bg-surface/50 border border-glass-border rounded-xl py-2.5 pl-10 pr-4 text-text-primary focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all" />
        </div>
    </div>

    {#if loading}
        <div class="text-center py-12 text-text-secondary">Loading...</div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each filteredCategories as cat}
                <div class="glass-card p-6 bg-surface/40 hover:shadow-lg transition-all duration-300 group relative">
                    <div class="flex justify-between items-start mb-4">
                        <div class="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center text-primary border border-primary/20"><Tag size={24} /></div>
                        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button on:click={() => openEditModal(cat)} class="p-2 hover:bg-primary/10 text-text-secondary hover:text-primary rounded-lg transition-colors"><Edit size={18} /></button>
                            <button on:click={() => deleteCategory(cat.category_id)} class="p-2 hover:bg-status-error/10 text-text-secondary hover:text-status-error rounded-lg transition-colors"><Trash2 size={18} /></button>
                        </div>
                    </div>
                    <h3 class="text-xl font-bold text-text-primary mb-2">{cat.name}</h3>
                    <p class="text-text-secondary text-sm">{cat.description || "No description."}</p>
                </div>
            {/each}
        </div>
    {/if}
</div>

{#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm" role="dialog" aria-modal="true" tabindex="-1">
        <button type="button" class="absolute inset-0 w-full h-full cursor-default focus:outline-none" on:click={closeModal} aria-label="Close modal"></button>
        <div class="glass-card w-full max-w-md shadow-2xl overflow-hidden bg-surface relative z-10">
            <div class="p-6 border-b border-glass-border flex justify-between items-center bg-surface/80">
                <h2 class="text-xl font-bold text-text-primary font-serif">{isEditing ? 'Edit Category' : 'New Category'}</h2>
                <button on:click={closeModal} class="text-text-secondary hover:text-primary"><X size={24} /></button>
            </div>
            <form on:submit|preventDefault={handleSubmit} class="p-6 space-y-5 bg-surface/50">
                {#if modalError}<div class="p-3 bg-status-error/10 border border-status-error/20 text-status-error rounded-xl text-sm font-medium">{modalError}</div>{/if}
                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-bold text-text-primary mb-1.5 ml-1" for="name">Name</label>
                        <input bind:value={formData.name} id="name" required class="w-full bg-main/50 border border-glass-border rounded-xl py-2.5 px-4 text-text-primary focus:border-primary outline-none transition-all" />
                    </div>
                    <div>
                        <label class="block text-sm font-bold text-text-primary mb-1.5 ml-1" for="desc">Description</label>
                        <textarea bind:value={formData.description} id="desc" rows="3" class="w-full bg-main/50 border border-glass-border rounded-xl py-2.5 px-4 text-text-primary focus:border-primary outline-none transition-all"></textarea>
                    </div>
                </div>
                <div class="pt-4 flex gap-3">
                    <button type="button" on:click={closeModal} class="flex-1 py-3 rounded-xl border border-glass-border hover:bg-main text-text-secondary font-bold">Cancel</button>
                    <button type="submit" disabled={submitting} class="flex-1 py-3 rounded-xl bg-primary hover:bg-primary-dark text-text-inverse font-bold transition-all disabled:opacity-50">Save Category</button>
                </div>
            </form>
        </div>
    </div>
{/if}