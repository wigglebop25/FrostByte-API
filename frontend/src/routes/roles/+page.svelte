<script lang="ts">
    /**
     * Role & Permission Management Page
     *
     * Provides CRUD operations for RBAC roles and their associated
     * permissions. Admin-only access. Roles control navigation
     * visibility and API endpoint authorization.
     */
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { Role } from "$lib/types";
    import { Plus, Trash2, Edit, X, Shield, Lock, ShieldCheck } from "lucide-svelte";

    let roles: Role[] = [];
    let loading = true;
    let searchQuery = "";
    
    let showModal = false;
    let isEditing = false;
    let modalError = "";
    let submitting = false;

    let formData = { name: "", description: "", permissions: "" };
    let editingId: number | null = null;

    async function fetchRoles() {
        loading = true;
        try {
            const res = await api.get("/roles");
            roles = res.data;
        } catch (e) {
            console.error("Failed to fetch roles", e);
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        fetchRoles();
    });

    function openAddModal() {
        isEditing = false;
        editingId = null;
        formData = { name: "", description: "", permissions: "" };
        modalError = "";
        showModal = true;
    }

    function openEditModal(role: any) {
        isEditing = true;
        editingId = role.role_id;
        formData = { name: role.name, description: role.description, permissions: role.permissions || "" };
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
                await api.post(`/roles/update/${editingId}`, { name: formData.name, description: formData.description });
                if (formData.permissions) {
                    await api.post(`/roles/${editingId}/set_permission`, { permission: formData.permissions });
                }
            } else {
                await api.post("/roles/create", { name: formData.name, description: formData.description });
            }
            showModal = false;
            fetchRoles();
        } catch (e: any) {
            console.error("Save error", e);
            modalError = e.response?.data?.error || "Failed to save role.";
        } finally {
            submitting = false;
        }
    }

    async function deleteRole(id: number) {
        if (!confirm("Are you sure?")) return;
        try {
            await api.delete(`/roles/${id}`);
            fetchRoles();
        } catch (e: any) {
            alert("Delete failed: " + (e.response?.data?.error || "Unknown error"));
        }
    }

    $: filteredRoles = roles.filter(r => r.name.toLowerCase().includes(searchQuery.toLowerCase()));
</script>

<div class="space-y-5">
    <div class="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div>
            <h1 class="text-2xl font-extrabold tracking-tight text-[var(--color-text-primary)]">Roles & Permissions</h1>
            <p class="text-[var(--color-text-secondary)] text-sm mt-0.5">Manage access control levels.</p>
        </div>
        <button on:click={openAddModal} class="glass-button flex items-center gap-2 text-sm">
            <Plus size={16} />
            Add Role
        </button>
    </div>

    <!-- Search -->
    <div class="max-w-md">
        <input bind:value={searchQuery} placeholder="Search roles..." class="glass-input w-full px-4 text-sm" />
    </div>

    {#if loading}
        <div class="text-center py-16 text-[var(--color-text-secondary)] text-sm">
            <div class="w-8 h-8 border-2 border-[var(--color-primary)]/30 border-t-[var(--color-primary)] rounded-full animate-spin mx-auto mb-3"></div>
            Loading...
        </div>
    {:else if filteredRoles.length === 0}
        <div class="text-center py-16 text-[var(--color-text-secondary)]">
            <ShieldCheck size={40} class="mx-auto mb-3 opacity-30" />
            <p class="text-sm">No roles found.</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each filteredRoles as role}
                <div class="glass-card p-5 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 group">
                    <div class="flex justify-between items-start mb-3">
                        <div class="w-10 h-10 rounded-xl bg-[var(--color-primary)]/10 flex items-center justify-center text-[var(--color-primary)]">
                            <Shield size={20} />
                        </div>
                        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button on:click={() => openEditModal(role)} class="p-2 hover:bg-[var(--color-primary)]/10 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] rounded-lg transition-colors">
                                <Edit size={16} />
                            </button>
                            <button on:click={() => deleteRole(role.role_id)} class="p-2 hover:bg-[var(--status-error)]/10 text-[var(--color-text-secondary)] hover:text-[var(--status-error)] rounded-lg transition-colors">
                                <Trash2 size={16} />
                            </button>
                        </div>
                    </div>
                    <h3 class="text-base font-bold text-[var(--color-text-primary)] mb-1">{role.name}</h3>
                    <p class="text-[var(--color-text-secondary)] text-sm mb-3 opacity-80">{role.description || "No description."}</p>
                    <div class="flex items-center gap-2 text-xs font-mono text-[var(--color-primary)] bg-[var(--color-primary)]/5 px-2.5 py-1 rounded-lg border border-[var(--color-primary)]/10 w-fit">
                        <Lock size={11} /> {role.permissions || "NONE"}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- Role Form Modal — Create / Edit with Permission Assignment -->
{#if showModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-md" role="dialog" aria-modal="true" tabindex="-1">
        <button type="button" class="absolute inset-0 w-full h-full cursor-default focus:outline-none" on:click={closeModal} aria-label="Close modal"></button>
        <div class="glass-card-solid w-full max-w-md shadow-2xl overflow-hidden relative z-10">
            <div class="p-5 border-b border-[var(--glass-border-subtle)] flex justify-between items-center">
                <h2 class="text-lg font-bold text-[var(--color-text-primary)]">{isEditing ? 'Edit Role' : 'New Role'}</h2>
                <button on:click={closeModal} class="p-2 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] rounded-xl hover:bg-[var(--glass-bg)] transition-colors"><X size={20} /></button>
            </div>
            <form on:submit|preventDefault={handleSubmit} class="p-5 space-y-4">
                {#if modalError}<div class="p-3 bg-[var(--status-error)]/10 border border-[var(--status-error)]/20 text-[var(--status-error)] rounded-xl text-sm font-medium">{modalError}</div>{/if}
                <div class="space-y-3">
                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="name">Role Name</label>
                        <input bind:value={formData.name} id="name" required class="glass-input w-full text-sm" />
                    </div>
                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="desc">Description</label>
                        <textarea bind:value={formData.description} id="desc" rows="2" class="glass-input w-full text-sm"></textarea>
                    </div>
                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="perm">Permissions</label>
                        <input bind:value={formData.permissions} id="perm" placeholder="READ,WRITE,ADMIN" class="glass-input w-full text-sm" />
                    </div>
                </div>
                <div class="pt-2 flex gap-2.5">
                    <button type="button" on:click={closeModal} class="flex-1 py-2.5 rounded-xl bg-[var(--glass-bg)] border border-[var(--glass-border)] text-[var(--color-text-secondary)] font-semibold text-sm hover:bg-[var(--color-bg-surface)] transition-colors">Cancel</button>
                    <button type="submit" disabled={submitting} class="flex-1 glass-button text-sm disabled:opacity-50">Save Role</button>
                </div>
            </form>
        </div>
    </div>
{/if}