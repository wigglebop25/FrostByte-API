<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/utils/api";
    import type { User, Role } from "$lib/types";
    import { Search, Plus, ShieldAlert, Trash2, Edit, X, User as UserIcon, Shield } from "lucide-svelte";

    let users: User[] = [];
    let roles: Role[] = [];
    let loading = true;
    let searchQuery = "";
    let error = "";

    // Modal State
    let showModal = false;
    let isEditing = false;
    let modalError = "";
    let submitting = false;

    // Form State
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
            password: "", // Password not sent in edit unless changed
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
                // Update User
                await api.put(`/users/${editingId}`, {
                    username: formData.username,
                    role_name: formData.role_name,
                    ...(formData.password ? { password: formData.password } : {})
                });
            } else {
                // Create User
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
            case 'ADMIN': return 'text-purple-500 bg-purple-500/10 border-purple-500/20';
            case 'CASHIER': return 'text-status-info bg-status-info/10 border-status-info/20';
            case 'CUSTOMER': return 'text-status-success bg-status-success/10 border-status-success/20';
            default: return 'text-text-secondary bg-text-secondary/10 border-text-secondary/20';
        }
    }
</script>

<div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row gap-4 justify-between items-start md:items-center">
        <div>
            <h1 class="text-3xl font-bold tracking-tight text-accent font-serif">User Management</h1>
            <p class="text-text-secondary mt-1">Manage system access, roles, and permissions.</p>
        </div>
        <button on:click={openAddModal} class="flex items-center gap-2 bg-primary hover:bg-primary-dark px-5 py-2.5 rounded-xl font-bold transition-all text-text-inverse shadow-lg shadow-primary/20 active:scale-95">
            <Plus size={18} />
            <span>Add New User</span>
        </button>
    </div>

    <!-- Error State -->
    {#if error}
        <div class="bg-status-error/10 border border-status-error/20 text-status-error p-4 rounded-xl flex items-center gap-3">
            <ShieldAlert size={20} />
            <span class="font-medium">{error}</span>
        </div>
    {/if}

    <!-- Filters -->
    <div class="glass-card p-4">
        <div class="relative max-w-md">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-text-secondary" size={18} />
            <input 
                bind:value={searchQuery}
                type="text" 
                placeholder="Search by username or ID..." 
                class="w-full bg-surface/50 border border-glass-border rounded-xl py-2.5 pl-10 pr-4 text-text-primary focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all shadow-sm"
            />
        </div>
    </div>

    <!-- User Grid -->
    {#if loading}
        <div class="text-center py-12 text-text-secondary">Loading users...</div>
    {:else if filteredUsers.length === 0}
        <div class="text-center py-12 text-text-secondary">No users found.</div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {#each filteredUsers as user}
                <div class="glass-card p-6 hover:shadow-lg transition-all duration-300 bg-surface/40 group relative">
                    <div class="flex justify-between items-start mb-4">
                        <div class="w-14 h-14 rounded-2xl bg-primary/10 flex items-center justify-center text-primary shadow-inner">
                            <UserIcon size={28} />
                        </div>
                        <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button on:click={() => openEditModal(user)} class="p-2 hover:bg-primary/10 text-text-secondary hover:text-primary rounded-lg transition-colors">
                                <Edit size={18} />
                            </button>
                            <button on:click={() => deleteUser(user.user_id)} class="p-2 hover:bg-status-error/10 text-text-secondary hover:text-status-error rounded-lg transition-colors">
                                <Trash2 size={18} />
                            </button>
                        </div>
                    </div>
                    
                    <h3 class="text-xl font-bold text-text-primary mb-1">{user.username}</h3>
                    <p class="text-text-secondary text-xs font-mono mb-4">UID: {user.user_id.toString().padStart(4, '0')}</p>

                    <div class="flex flex-wrap gap-2">
                        {#if user.roles && user.roles.length > 0}
                            {#each user.roles as role}
                                <span class="px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider border {getRoleBadgeColor(role.name)}">
                                    {role.name}
                                </span>
                            {/each}
                        {:else}
                            <span class="text-text-secondary text-xs italic">No roles assigned</span>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- Add/Edit User Modal -->
{#if showModal}
    <div 
        role="button"
        tabindex="0"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm transition-all w-full h-full cursor-default outline-none" 
        on:click={closeModal}
        on:keydown={(e) => e.key === 'Escape' && closeModal()}
    >
        <div 
            class="glass-card w-full max-w-md shadow-2xl overflow-hidden bg-surface relative cursor-auto text-left" 
            on:click|stopPropagation
            on:keydown|stopPropagation
            role="document"
            tabindex="-1"
        >
            <div class="p-6 border-b border-glass-border flex justify-between items-center bg-surface/80">
                <h2 class="text-xl font-bold text-text-primary font-serif">
                    {isEditing ? 'Edit User' : 'Register New User'}
                </h2>
                <button on:click={closeModal} class="text-text-secondary hover:text-primary transition-colors p-2 hover:bg-main rounded-full">
                    <X size={24} />
                </button>
            </div>
            
            <form on:submit|preventDefault={handleSubmit} class="p-6 space-y-5 bg-surface/50">
                {#if modalError}
                    <div class="p-3 bg-status-error/10 border border-status-error/20 text-status-error rounded-xl text-sm font-medium animate-pulse">
                        {modalError}
                    </div>
                {/if}

                <div class="space-y-4">
                    <div>
                        <label class="block text-sm font-bold text-text-primary mb-1.5 ml-1" for="username">Username</label>
                        <div class="relative">
                            <UserIcon class="absolute left-3.5 top-1/2 -translate-y-1/2 text-text-secondary" size={18} />
                            <input bind:value={formData.username} id="username" required 
                                class="w-full bg-main/50 border border-glass-border rounded-xl py-2.5 pl-11 pr-4 text-text-primary focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all" />
                        </div>
                    </div>

                    <div>
                        <label class="block text-sm font-bold text-text-primary mb-1.5 ml-1" for="password">
                            {isEditing ? 'New Password (Leave blank to keep)' : 'Password'}
                        </label>
                        <div class="relative">
                            <Plus class="absolute left-3.5 top-1/2 -translate-y-1/2 text-text-secondary" size={18} />
                            <input bind:value={formData.password} id="password" type="password" required={!isEditing}
                                class="w-full bg-main/50 border border-glass-border rounded-xl py-2.5 pl-11 pr-4 text-text-primary focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all" />
                        </div>
                    </div>

                    <div>
                        <label class="block text-sm font-bold text-text-primary mb-1.5 ml-1" for="role">Assign Role</label>
                        <div class="relative">
                            <Shield class="absolute left-3.5 top-1/2 -translate-y-1/2 text-text-secondary" size={18} />
                            <select bind:value={formData.role_name} id="role" 
                                class="w-full bg-main/50 border border-glass-border rounded-xl py-2.5 pl-11 pr-4 text-text-primary focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all appearance-none">
                                {#each roles as role}
                                    <option value={role.name}>{role.name}</option>
                                {/each}
                            </select>
                        </div>
                    </div>
                </div>

                <div class="pt-4 flex gap-3">
                    <button type="button" on:click={closeModal} class="flex-1 py-3 rounded-xl border border-glass-border hover:bg-main text-text-secondary font-bold transition-all">
                        Cancel
                    </button>
                    <button type="submit" disabled={submitting} 
                        class="flex-1 py-3 rounded-xl bg-primary hover:bg-primary-dark text-text-inverse font-bold transition-all shadow-lg shadow-primary/20 disabled:opacity-50">
                        {submitting ? 'Saving...' : 'Save User'}
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}
