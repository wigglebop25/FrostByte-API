<script lang="ts">
    /**
     * User Administration Page
     *
     * Manages system users with CRUD operations and role assignment.
     * Restricted to Admin role. Passwords are hashed server-side
     * with Argon2 before storage.
     */
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { User, Role } from "$lib/types";
    import { Search, Plus, ShieldAlert, Trash2, Edit, X, User as UserIcon, Shield, UsersRound } from "lucide-svelte";

    let users: User[] = [];
    let roles: Role[] = [];
    let loading = true;
    let searchQuery = "";
    let error = "";

    let showModal = false;
    let isEditing = false;
    let modalError = "";
    let submitting = false;

    let formData = {
        username: "",
        password: "",
        role_name: "Customer"
    };
    let editingId: number | null = null;

    async function fetchData() {
        loading = true;
        error = "";
        try {
            const [usersRes, rolesRes] = await Promise.all([
                api.get("/users"),
                api.get("/roles")
            ]);
            users = usersRes.data;
            roles = rolesRes.data;
        } catch (e: any) {
            console.error("Failed to fetch user data", e);
            error = e.response?.data?.error || "Failed to load users. Are you an Admin?";
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        fetchData();
    });

    function openAddModal() {
        isEditing = false;
        editingId = null;
        formData = { username: "", password: "", role_name: "Customer" };
        modalError = "";
        showModal = true;
    }

    function openEditModal(user: User) {
        isEditing = true;
        editingId = user.user_id;
        formData = {
            username: user.username,
            password: "",
            role_name: user.roles?.[0]?.name || "Customer"
        };
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
                await api.put(`/users/${editingId}`, {
                    username: formData.username,
                    role_name: formData.role_name,
                    ...(formData.password ? { password: formData.password } : {})
                });
            } else {
                await api.post("/users", formData);
            }
            showModal = false;
            fetchData();
        } catch (e: any) {
            console.error("User save error", e);
            modalError = e.response?.data?.error || "Failed to save user.";
        } finally {
            submitting = false;
        }
    }

    async function deleteUser(id: number) {
        if (!confirm("Are you sure you want to delete this user?")) return;
        try {
            await api.delete(`/users/${id}`);
            fetchData();
        } catch (e: any) {
            alert("Failed to delete user: " + (e.response?.data?.error || "Unknown error"));
        }
    }

    $: filteredUsers = users.filter(user => 
        user.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.user_id.toString().includes(searchQuery)
    );

    function getRoleBadgeColor(roleName: string) {
        switch (roleName.toUpperCase()) {
            case 'ADMIN': return 'text-violet-600 dark:text-violet-400 bg-violet-500/10 border-violet-500/20';
            case 'CASHIER': return 'text-blue-600 dark:text-blue-400 bg-blue-500/10 border-blue-500/20';
            case 'CUSTOMER': return 'text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 border-emerald-500/20';
            default: return 'text-[var(--color-text-secondary)] bg-[var(--glass-bg)]';
        }
    }
</script>

<div class="space-y-5">
    <div class="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
        <div>
            <h1 class="text-2xl font-extrabold tracking-tight text-[var(--color-text-primary)]">User Management</h1>
            <p class="text-[var(--color-text-secondary)] text-sm mt-0.5">Manage system access, roles, and permissions.</p>
        </div>
        <button on:click={openAddModal} class="glass-button flex items-center gap-2 text-sm">
            <Plus size={16} />
            Add User
        </button>
    </div>

    {#if error}
        <div class="bg-[var(--status-error)]/10 border border-[var(--status-error)]/20 text-[var(--status-error)] p-3.5 rounded-xl flex items-center gap-3 text-sm">
            <ShieldAlert size={18} />
            <span class="font-medium">{error}</span>
        </div>
    {/if}

    <!-- Search -->
    <div class="relative max-w-md">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-text-secondary)]" size={16} />
        <input bind:value={searchQuery} type="text" placeholder="Search by username or ID..." class="glass-input w-full pl-12 pr-4 text-sm" />
    </div>

    {#if loading}
        <div class="text-center py-16 text-[var(--color-text-secondary)] text-sm">
            <div class="w-8 h-8 border-2 border-[var(--color-primary)]/30 border-t-[var(--color-primary)] rounded-full animate-spin mx-auto mb-3"></div>
            Loading users...
        </div>
    {:else if filteredUsers.length === 0}
        <div class="text-center py-16 text-[var(--color-text-secondary)]">
            <UsersRound size={40} class="mx-auto mb-3 opacity-30" />
            <p class="text-sm">No users found.</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {#each filteredUsers as user}
                <div class="glass-card p-5 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 group">
                    <div class="flex justify-between items-start mb-3">
                        <div class="w-11 h-11 rounded-2xl bg-[var(--color-primary)]/10 flex items-center justify-center text-[var(--color-primary)] text-sm font-bold">
                            {user.username[0]?.toUpperCase() || 'U'}
                        </div>
                        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button on:click={() => openEditModal(user)} class="p-2 hover:bg-[var(--color-primary)]/10 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] rounded-lg transition-colors">
                                <Edit size={15} />
                            </button>
                            <button on:click={() => deleteUser(user.user_id)} class="p-2 hover:bg-[var(--status-error)]/10 text-[var(--color-text-secondary)] hover:text-[var(--status-error)] rounded-lg transition-colors">
                                <Trash2 size={15} />
                            </button>
                        </div>
                    </div>
                    
                    <h3 class="text-base font-bold text-[var(--color-text-primary)] mb-0.5">{user.username}</h3>
                    <p class="text-[var(--color-text-secondary)] text-xs font-mono mb-3 opacity-60">UID: {user.user_id.toString().padStart(4, '0')}</p>

                    <div class="flex flex-wrap gap-1.5">
                        {#if user.roles && user.roles.length > 0}
                            {#each user.roles as role}
                                <span class="px-2.5 py-0.5 rounded-lg text-[10px] font-bold uppercase tracking-wider border {getRoleBadgeColor(role.name)}">
                                    {role.name}
                                </span>
                            {/each}
                        {:else}
                            <span class="text-[var(--color-text-secondary)] text-xs opacity-60">No roles</span>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- User Form Modal — Create / Edit Operations -->
{#if showModal}
    <div 
        role="button"
        tabindex="0"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-md w-full h-full cursor-default outline-none" 
        on:click={closeModal}
        on:keydown={(e) => e.key === 'Escape' && closeModal()}
    >
        <div 
            class="glass-card-solid w-full max-w-md shadow-2xl overflow-hidden relative cursor-auto text-left" 
            on:click|stopPropagation
            on:keydown|stopPropagation
            role="document"
            tabindex="-1"
        >
            <div class="p-5 border-b border-[var(--glass-border-subtle)] flex justify-between items-center">
                <h2 class="text-lg font-bold text-[var(--color-text-primary)]">
                    {isEditing ? 'Edit User' : 'Register New User'}
                </h2>
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

                <div class="space-y-3">
                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="username">Username</label>
                        <div class="relative">
                            <UserIcon class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-text-secondary)]" size={16} />
                            <input bind:value={formData.username} id="username" required class="glass-input w-full pl-12 text-sm" />
                        </div>
                    </div>

                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="password">
                            {isEditing ? 'New Password (Leave blank to keep)' : 'Password'}
                        </label>
                        <input bind:value={formData.password} id="password" type="password" required={!isEditing} class="glass-input w-full text-sm" />
                    </div>

                    <div>
                        <label class="block text-sm font-semibold text-[var(--color-text-primary)] mb-1.5" for="role">Assign Role</label>
                        <div class="relative">
                            <Shield class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--color-text-secondary)]" size={16} />
                            <select bind:value={formData.role_name} id="role" class="glass-input w-full pl-12 text-sm appearance-none cursor-pointer">
                                {#each roles as role}
                                    <option value={role.name}>{role.name}</option>
                                {/each}
                            </select>
                        </div>
                    </div>
                </div>

                <div class="pt-2 flex gap-2.5">
                    <button type="button" on:click={closeModal} class="flex-1 py-2.5 rounded-xl bg-[var(--glass-bg)] border border-[var(--glass-border)] text-[var(--color-text-secondary)] font-semibold text-sm hover:bg-[var(--color-bg-surface)] transition-colors">
                        Cancel
                    </button>
                    <button type="submit" disabled={submitting} class="flex-1 glass-button text-sm disabled:opacity-50">
                        {submitting ? 'Saving...' : 'Save User'}
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}
