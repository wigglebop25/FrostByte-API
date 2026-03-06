<script lang="ts">
    import api from "$lib/utils/api";
    import { goto } from "$app/navigation";
    import { auth } from "$lib/stores/auth"; // Use the new store
    import { Lock, User, ArrowRight, Eye, EyeOff } from "lucide-svelte";
    import { onMount } from "svelte";

    let username = "";
    let password = "";
    let showPassword = false;
    let error = "";
    let loading = false;

    onMount(() => {
        auth.init();
        if (localStorage.getItem("token")) {
            goto("/");
        }
    });

    async function handleLogin() {
        loading = true;
        error = "";
        try {
            const res = await api.post("/auth/login", { username, password });
            
            if (res.data.token) {
                auth.login(res.data.token, res.data.user);
                
                const roles = res.data.user.roles?.map((r: any) => r.name.toLowerCase()) || [];
                if (roles.includes('customer')) {
                    goto("/orders/create");
                } else {
                    goto("/");
                }
            }
        } catch (e: any) {
            console.error("Login Error Full:", e);
            if (e.code === "ERR_NETWORK") {
                error = "Network Error: Unable to reach server. Check internet or CORS.";
            } else {
                error = e.response?.data?.error || "Invalid username or password";
            }
        } finally {
            loading = false;
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center p-4 bg-main relative overflow-hidden transition-colors duration-300">
    <!-- Background Patterns -->
    <div class="absolute top-0 right-0 w-[500px] h-[500px] bg-primary/5 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2 pointer-events-none"></div>
    <div class="absolute bottom-0 left-0 w-[500px] h-[500px] bg-status-info/5 rounded-full blur-3xl translate-y-1/2 -translate-x-1/2 pointer-events-none"></div>

    <div class="w-full max-w-md glass-card p-10 relative z-10 border-glass-border shadow-xl">
        <div class="text-center mb-10">
            <!-- Logo -->
            <div class="flex justify-center mb-6">
                <div class="relative w-24 h-24 rounded-full bg-white p-1 shadow-lg shadow-primary/20">
                    <img src="/images/itadaki_logo.png" alt="Itadaki Logo" class="w-full h-full object-contain rounded-full" />
                </div>
            </div>
            
            <h1 class="text-3xl font-bold text-accent mb-2 font-serif tracking-tight">Welcome Back</h1>
            <p class="text-text-secondary">Sign in to manage your restaurant</p>
        </div>

        <form on:submit|preventDefault={handleLogin} class="space-y-6">
            {#if error}
                <div class="bg-status-error/10 border border-status-error/20 text-status-error p-4 rounded-xl text-sm text-center font-medium animate-pulse">
                    {error}
                </div>
            {/if}

            <div class="space-y-2">
                <label class="text-sm font-bold text-text-primary ml-1" for="username">Username</label>
                <div class="relative group">
                    <User class="absolute left-4 top-1/2 -translate-y-1/2 text-text-secondary group-focus-within:text-primary transition-colors" size={20} />
                    <input bind:value={username} type="text" id="username" required placeholder="Enter your username"
                        class="w-full bg-surface/50 border border-glass-border rounded-xl py-3.5 pl-12 pr-4 text-text-primary placeholder:text-text-secondary/50 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all" />
                </div>
            </div>

            <div class="space-y-2">
                <label class="text-sm font-bold text-text-primary ml-1" for="password">Password</label>
                <div class="relative group">
                    <Lock class="absolute left-4 top-1/2 -translate-y-1/2 text-text-secondary group-focus-within:text-primary transition-colors" size={20} />
                    <input bind:value={password} type={showPassword ? 'text' : 'password'} id="password" required placeholder="********"
                        class="w-full bg-surface/50 border border-glass-border rounded-xl py-3.5 pl-12 pr-12 text-text-primary placeholder:text-text-secondary/50 focus:border-primary focus:ring-1 focus:ring-primary outline-none transition-all" />
                    <button type="button" on:click={() => showPassword = !showPassword} class="absolute right-4 top-1/2 -translate-y-1/2 text-text-secondary hover:text-primary transition-colors focus:outline-none cursor-pointer">
                        {#if showPassword}
                            <EyeOff size={20} />
                        {:else}
                            <Eye size={20} />
                        {/if}
                    </button>
                </div>
            </div>

            <button type="submit" disabled={loading}
                class="w-full bg-primary hover:bg-primary-dark text-text-inverse font-bold py-4 rounded-xl transition-all shadow-lg shadow-primary/20 disabled:opacity-50 flex items-center justify-center gap-2 group mt-4 active:scale-95">
                {#if loading}
                    <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    <span>Authenticating...</span>
                {:else}
                    <span>Sign In</span>
                    <ArrowRight size={20} class="group-hover:translate-x-1 transition-transform" />
                {/if}
            </button>
        </form>
        
        <div class="mt-8 text-center text-sm text-text-secondary opacity-60">
            Copyright 2026 Itadaki Restaurant System
        </div>
    </div>
</div>