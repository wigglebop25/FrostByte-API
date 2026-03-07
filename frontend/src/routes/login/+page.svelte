<script lang="ts">
    /**
     * Authentication Login Page
     *
     * Handles user credential submission via the /auth/login API endpoint.
     * Redirects authenticated users based on role: Customers to order
     * creation, Admin/Cashier to the dashboard.
     */
    import api from "$lib/utils/api";
    import { goto } from "$app/navigation";
    import { auth } from "$lib/stores/auth";
    import { theme } from "$lib/stores/theme";
    import { Lock, User, ArrowRight, Eye, EyeOff } from "lucide-svelte";
    import { onMount } from "svelte";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";

    let username = "";
    let password = "";
    let showPassword = false;
    let error = "";
    let loading = false;

    onMount(() => {
        auth.init();
        theme.init();
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
                auth.login(res.data.token, res.data.user, res.data.refresh_token);
                
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

<div class="min-h-dvh flex items-center justify-center p-4 bg-[var(--color-bg-main)] relative overflow-hidden">
    <!-- Decorative Animated Mesh Background -->
    <div class="mesh-bg">
        <div class="mesh-orb w-[300px] sm:w-[450px] lg:w-[600px] h-[300px] sm:h-[450px] lg:h-[600px] top-[-15%] right-[-10%]" style="background: var(--mesh-1);"></div>
        <div class="mesh-orb w-[250px] sm:w-[400px] lg:w-[500px] h-[250px] sm:h-[400px] lg:h-[500px] bottom-[-10%] left-[-5%]" style="background: var(--mesh-2); animation-delay: -8s;"></div>
        <div class="mesh-orb w-[200px] sm:w-[280px] lg:w-[350px] h-[200px] sm:h-[280px] lg:h-[350px] top-[30%] left-[20%]" style="background: var(--mesh-3); animation-delay: -15s;"></div>
    </div>

    <!-- Theme Toggle -->
    <div class="absolute top-6 right-6 z-20">
        <ThemeToggle compact={true} />
    </div>

    <!-- Authentication Card -->
    <div class="w-full max-w-md glass-card p-6 sm:p-10 relative z-10 shadow-xl">
        <!-- Branding Section -->
        <div class="text-center mb-6 sm:mb-10">
            <div class="flex justify-center mb-4 sm:mb-6">
                <div class="relative w-16 h-16 sm:w-20 sm:h-20 rounded-3xl bg-[var(--color-bg-surface)] p-2 shadow-lg">
                    <img src="/images/itadaki_logo.png" alt="Itadaki Logo" class="w-full h-full object-contain rounded-2xl" />
                </div>
            </div>
            
            <h1 class="text-2xl sm:text-3xl font-extrabold text-[var(--color-text-primary)] mb-2 tracking-tight">Welcome Back</h1>
            <p class="text-[var(--color-text-secondary)] text-sm">Sign in to manage your restaurant</p>
        </div>

        <!-- Credential Input Form -->
        <form on:submit|preventDefault={handleLogin} class="space-y-5">
            {#if error}
                <div class="bg-[var(--status-error)]/10 border border-[var(--status-error)]/20 text-[var(--status-error)] p-3.5 rounded-xl text-sm text-center font-medium">
                    {error}
                </div>
            {/if}

            <div class="space-y-1.5">
                <label class="text-sm font-semibold text-[var(--color-text-primary)] ml-1 flex items-center gap-1.5" for="username">
                    Username
                    <User size={15} class="text-[var(--color-text-secondary)]" />
                </label>
                <div class="relative group">
                    <input bind:value={username} type="text" id="username" required placeholder="Enter your username"
                        class="glass-input w-full px-4" />
                </div>
            </div>

            <div class="space-y-1.5">
                <label class="text-sm font-semibold text-[var(--color-text-primary)] ml-1 flex items-center gap-1.5" for="password">
                    Password
                    <Lock size={15} class="text-[var(--color-text-secondary)]" />
                </label>
                <div class="relative group">
                    <input bind:value={password} type={showPassword ? 'text' : 'password'} id="password" required placeholder="••••••••"
                        class="glass-input w-full pl-4 pr-12" />
                    <button type="button" on:click={() => showPassword = !showPassword} class="absolute right-3.5 top-1/2 -translate-y-1/2 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors focus:outline-none cursor-pointer">
                        {#if showPassword}
                            <EyeOff size={18} />
                        {:else}
                            <Eye size={18} />
                        {/if}
                    </button>
                </div>
            </div>

            <button type="submit" disabled={loading}
                class="glass-button w-full py-3.5 flex items-center justify-center gap-2 group mt-2 disabled:opacity-50">
                {#if loading}
                    <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    <span>Authenticating...</span>
                {:else}
                    <span>Sign In</span>
                    <ArrowRight size={18} class="group-hover:translate-x-1 transition-transform" />
                {/if}
            </button>
        </form>
        
        <div class="mt-8 text-center text-xs text-[var(--color-text-secondary)] opacity-50">
            © 2026 Itadaki Restaurant System
        </div>
    </div>
</div>