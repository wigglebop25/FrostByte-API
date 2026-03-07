<script lang="ts">
    /**
     * Admin Dashboard Page
     *
     * Displays real-time analytics including revenue, order counts,
     * and distribution metrics. Refreshes automatically via WebSocket
     * events. Restricted to Admin and Cashier roles.
     */
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { ws } from '$lib/stores/websocket';
    import { env } from '$env/dynamic/public';
    import api from '$lib/utils/api';
    import RevenueChart from '$lib/components/RevenueChart.svelte';
    import { TrendingUp, Users, ShoppingBag, Activity, CheckCircle, Package, DollarSign, Clock } from 'lucide-svelte';

    interface Stats {
        total_revenue: number;
        total_orders: number;
        pending_orders: number;
        accepted_orders: number;
        cooking_orders: number;
        ready_orders: number;
        completed_orders: number;
        cancelled_orders: number;
        average_order_value: number;
        daily_revenue?: Record<string, number>;
    }

    let stats: Stats = {
        total_revenue: 0,
        total_orders: 0,
        pending_orders: 0,
        accepted_orders: 0,
        cooking_orders: 0,
        ready_orders: 0,
        completed_orders: 0,
        cancelled_orders: 0,
        average_order_value: 0
    };
    
    let activeUsers = 0;
    let lastWsTs = 0;
    let revenueData: number[] = [0, 0, 0, 0, 0, 0, 0];
    let revenueLabels: string[] = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
    let greeting = 'Good day';
    let username = '';

    onMount(async () => {
        const token = localStorage.getItem('token');
        const userStr = localStorage.getItem('user');
        
        if (!token || !userStr) {
            goto("/login");
            return;
        }

        const user = JSON.parse(userStr);
        username = user.username || 'User';
        const roles = user.roles?.map((r: any) => r.name.toLowerCase()) || [];
        
        if (roles.includes('customer')) {
            goto("/orders/create");
            return;
        }

        const hour = new Date().getHours();
        if (hour < 12) greeting = 'Good morning';
        else if (hour < 17) greeting = 'Good afternoon';
        else greeting = 'Good evening';

        fetchAnalytics();
    });

    async function fetchAnalytics() {
        try {
            const res = await api.get('/analytics/dashboard');
            stats = res.data;
            
            if (stats.daily_revenue) {
                const dailyRevenue = stats.daily_revenue;
                const dates = Object.keys(dailyRevenue).sort();
                revenueData = dates.map(d => dailyRevenue[d]);
                revenueLabels = dates.map(d => {
                    const date = new Date(d);
                    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
                });
            }
        } catch (e) {
            console.error("Analytics error:", e);
        }
    }

    $: if ($ws.data && $ws.data._ts !== lastWsTs) {
        lastWsTs = $ws.data._ts;
        handleMessage($ws.data);
    }

    function handleMessage(msg: any) {
        if (msg.event === 'NEW_ORDER' || msg.event === 'ORDER_UPDATED' || msg.order_id) {
            fetchAnalytics();
        }
        
        if (msg.event === 'USER_JOINED') {
            activeUsers++;
        }
    }

    $: cards = [
        { label: 'Total Revenue', value: '$' + (stats.total_revenue || 0).toLocaleString(), icon: DollarSign, gradient: 'from-emerald-500 to-teal-600', bg: 'bg-emerald-500/10' },
        { label: 'Total Orders', value: (stats.total_orders || 0).toString(), icon: Package, gradient: 'from-blue-500 to-indigo-600', bg: 'bg-blue-500/10' },
        { label: 'Pending', value: (stats.pending_orders || 0).toString(), icon: Clock, gradient: 'from-amber-500 to-orange-600', bg: 'bg-amber-500/10' },
        { label: 'Completed', value: (stats.completed_orders || 0).toString(), icon: CheckCircle, gradient: 'from-violet-500 to-purple-600', bg: 'bg-violet-500/10' },
    ];
</script>

<div class="space-y-6">
    <!-- Greeting Header -->
    <div class="flex items-end justify-between">
        <div>
            <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight text-[var(--color-text-primary)]">
                {greeting}, <span class="text-[var(--color-primary)]">{username}</span>
            </h1>
            <p class="text-[var(--color-text-secondary)] mt-1 text-sm">Here's what's happening with your restaurant today.</p>
        </div>
        <div class="text-right hidden sm:block">
            <div class="text-xs font-semibold text-[var(--color-text-secondary)]">
                {new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
            </div>
        </div>
    </div>

    <!-- KPI Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {#each cards as card}
            <div class="glass-card p-5 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 group">
                <div class="flex items-center justify-between mb-3">
                    <span class="text-[var(--color-text-secondary)] text-xs font-semibold uppercase tracking-wider">{card.label}</span>
                    <div class="w-9 h-9 rounded-xl bg-gradient-to-br {card.gradient} flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                        <svelte:component this={card.icon} size={16} class="text-white" />
                    </div>
                </div>
                <div class="text-2xl font-extrabold text-[var(--color-text-primary)] tracking-tight">{card.value}</div>
            </div>
        {/each}
    </div>

    <!-- Analytics Charts Row -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Revenue Performance Chart — Live Updates via WebSocket -->
        <div class="lg:col-span-2 glass-card p-6">
            <div class="flex items-center justify-between mb-5">
                <div>
                    <h3 class="text-lg font-bold text-[var(--color-text-primary)]">Revenue Performance</h3>
                    <p class="text-xs text-[var(--color-text-secondary)] mt-0.5">Last 7 days</p>
                </div>
                <div class="flex items-center gap-2 text-xs font-semibold px-3 py-1.5 rounded-xl bg-[var(--status-info)]/10 text-[var(--status-info)] border border-[var(--status-info)]/20">
                    <span class="w-1.5 h-1.5 rounded-full bg-[var(--status-info)] animate-pulse"></span>
                    Live
                </div>
            </div>
            <RevenueChart data={revenueData} labels={revenueLabels} />
        </div>

        <!-- Order Status Distribution Breakdown -->
        <div class="glass-card p-6 flex flex-col">
            <h3 class="text-lg font-bold text-[var(--color-text-primary)] mb-5">Order Distribution</h3>
            
            <div class="space-y-4 flex-1">
                <!-- Completed -->
                <div class="space-y-1.5">
                    <div class="flex justify-between text-xs font-semibold">
                        <span class="text-[var(--color-text-secondary)]">Completed</span>
                        <span class="text-emerald-500">{Math.round((stats.completed_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-[var(--glass-bg)] rounded-full h-2 overflow-hidden border border-[var(--glass-border-subtle)]">
                        <div class="bg-gradient-to-r from-emerald-400 to-emerald-600 h-full rounded-full transition-all duration-1000 ease-out" style="width: {(stats.completed_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>

                <!-- Ready -->
                <div class="space-y-1.5">
                    <div class="flex justify-between text-xs font-semibold">
                        <span class="text-[var(--color-text-secondary)]">Ready</span>
                        <span class="text-blue-500">{Math.round((stats.ready_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-[var(--glass-bg)] rounded-full h-2 overflow-hidden border border-[var(--glass-border-subtle)]">
                        <div class="bg-gradient-to-r from-blue-400 to-blue-600 h-full rounded-full transition-all duration-1000 ease-out" style="width: {(stats.ready_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>

                <!-- Cooking -->
                <div class="space-y-1.5">
                    <div class="flex justify-between text-xs font-semibold">
                        <span class="text-[var(--color-text-secondary)]">Cooking</span>
                        <span class="text-orange-500">{Math.round((stats.cooking_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-[var(--glass-bg)] rounded-full h-2 overflow-hidden border border-[var(--glass-border-subtle)]">
                        <div class="bg-gradient-to-r from-orange-400 to-orange-600 h-full rounded-full transition-all duration-1000 ease-out" style="width: {(stats.cooking_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>

                <!-- Accepted -->
                <div class="space-y-1.5">
                    <div class="flex justify-between text-xs font-semibold">
                        <span class="text-[var(--color-text-secondary)]">Accepted</span>
                        <span class="text-indigo-500">{Math.round((stats.accepted_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-[var(--glass-bg)] rounded-full h-2 overflow-hidden border border-[var(--glass-border-subtle)]">
                        <div class="bg-gradient-to-r from-indigo-400 to-indigo-600 h-full rounded-full transition-all duration-1000 ease-out" style="width: {(stats.accepted_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>

                <!-- Pending -->
                <div class="space-y-1.5">
                    <div class="flex justify-between text-xs font-semibold">
                        <span class="text-[var(--color-text-secondary)]">Pending</span>
                        <span class="text-amber-500">{Math.round((stats.pending_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-[var(--glass-bg)] rounded-full h-2 overflow-hidden border border-[var(--glass-border-subtle)]">
                        <div class="bg-gradient-to-r from-amber-400 to-amber-600 h-full rounded-full transition-all duration-1000 ease-out" style="width: {(stats.pending_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>

                <!-- Cancelled -->
                <div class="space-y-1.5">
                    <div class="flex justify-between text-xs font-semibold">
                        <span class="text-[var(--color-text-secondary)]">Cancelled</span>
                        <span class="text-red-500">{Math.round((stats.cancelled_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-[var(--glass-bg)] rounded-full h-2 overflow-hidden border border-[var(--glass-border-subtle)]">
                        <div class="bg-gradient-to-r from-red-400 to-red-600 h-full rounded-full transition-all duration-1000 ease-out" style="width: {(stats.cancelled_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>
            </div>

            <!-- Avg Order Value -->
            <div class="mt-6 pt-5 border-t border-[var(--glass-border-subtle)]">
                <div class="flex items-center justify-between">
                    <span class="text-xs font-semibold text-[var(--color-text-secondary)] uppercase tracking-wider">Avg Order Value</span>
                    <span class="text-xl font-extrabold text-[var(--color-primary)]">${(stats.average_order_value || 0).toFixed(2)}</span>
                </div>
            </div>
        </div>
    </div>
</div>