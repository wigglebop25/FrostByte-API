<script lang="ts">
    import { onMount } from 'svelte';
    import { ws } from '$lib/stores/websocket';
    import { env } from '$env/dynamic/public';
    import api from '$lib/utils/api';
    import RevenueChart from '$lib/components/RevenueChart.svelte';
    import { TrendingUp, Users, ShoppingBag, Activity, CheckCircle, Package } from 'lucide-svelte';

    let stats = {
        total_revenue: 0,
        total_orders: 0,
        pending_orders: 0,
        completed_orders: 0,
        cancelled_orders: 0,
        average_order_value: 0
    };
    
    let activeUsers = 0;
    let revenueData: number[] = [0, 0, 0, 0, 0, 0, 0];
    let revenueLabels: string[] = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

    onMount(async () => {
        const token = localStorage.getItem('token');
        const userStr = localStorage.getItem('user');
        
        if (!token || !userStr) {
            goto("/login");
            return;
        }

        const user = JSON.parse(userStr);
        const roles = user.roles?.map((r: any) => r.name.toLowerCase()) || [];
        
        if (roles.includes('customer')) {
            goto("/orders/create");
            return;
        }

        const wsUrl = env.PUBLIC_WS_URL || 'ws://localhost:8080/ws';
        ws.connect(wsUrl, token);
        fetchAnalytics();
    });

    async function fetchAnalytics() {
        try {
            const res = await api.get('/analytics/dashboard');
            stats = res.data;
            
            if (stats.daily_revenue) {
                const dates = Object.keys(stats.daily_revenue).sort();
                revenueData = dates.map(d => stats.daily_revenue[d]);
                revenueLabels = dates.map(d => {
                    const date = new Date(d);
                    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
                });
            }
        } catch (e) {
            console.error("Analytics error:", e);
        }
    }

    $: if ($ws.data) {
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
        { label: 'Total Revenue', value: '$' + (stats.total_revenue || 0).toLocaleString(), icon: TrendingUp, color: 'text-emerald-600', bg: 'bg-emerald-100' },
        { label: 'Total Orders', value: (stats.total_orders || 0).toString(), icon: Package, color: 'text-blue-600', bg: 'bg-blue-100' },
        { label: 'Pending Orders', value: (stats.pending_orders || 0).toString(), icon: ShoppingBag, color: 'text-orange-600', bg: 'bg-orange-100' },
        { label: 'Completed', value: (stats.completed_orders || 0).toString(), icon: CheckCircle, color: 'text-purple-600', bg: 'bg-purple-100' },
    ];
</script>

<div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
        <div>
            <h1 class="text-4xl font-bold tracking-tight text-accent font-serif">Dashboard</h1>
            <p class="text-text-secondary mt-1">Real-time performance metrics from Azure Backend.</p>
        </div>
        
        <div class="glass-card flex items-center gap-2 px-4 py-2 text-sm shadow-sm bg-surface/50">
            <span class="relative flex h-2.5 w-2.5">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 {$ws.connected ? 'bg-status-success' : 'bg-status-error'}"></span>
                <span class="relative inline-flex rounded-full h-2.5 w-2.5 {$ws.connected ? 'bg-status-success' : 'bg-status-error'}"></span>
            </span>
            <span class="font-bold {$ws.connected ? 'text-status-success' : 'text-status-error'}">
                {$ws.connected ? 'LIVE FEED' : 'OFFLINE'}
            </span>
        </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {#each cards as card}
            <div class="glass-card p-6 hover:shadow-lg transition-all duration-300 bg-surface/40 border-glass-border">
                <div class="flex items-center justify-between mb-4">
                    <span class="text-text-secondary text-[10px] font-bold uppercase tracking-widest">{card.label}</span>
                    <div class="p-2 rounded-xl {card.bg} shadow-sm">
                        <svelte:component this={card.icon} size={18} class={card.color} />
                    </div>
                </div>
                <div class="text-3xl font-black text-text-primary tracking-tight">{card.value}</div>
            </div>
        {/each}
    </div>

    <!-- Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-2 glass-card p-8 bg-surface/40 border-glass-border">
            <div class="flex items-center justify-between mb-8">
                <h3 class="text-xl font-bold text-text-primary font-serif">Revenue Performance</h3>
                <div class="text-xs font-bold text-text-secondary bg-main px-3 py-1 rounded-full border border-glass-border">
                    Last 7 Days
                </div>
            </div>
            <RevenueChart data={revenueData} labels={revenueLabels} />
        </div>

        <!-- System Health / Quick Stats -->
        <div class="glass-card p-8 bg-surface/40 border-glass-border flex flex-col">
            <h3 class="text-xl font-bold text-text-primary mb-8 font-serif">Order Distribution</h3>
            
            <div class="space-y-6 flex-1">
                <div class="space-y-2">
                    <div class="flex justify-between text-xs font-bold">
                        <span class="text-text-secondary">COMPLETED</span>
                        <span class="text-text-primary">{Math.round((stats.completed_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-main rounded-full h-2 overflow-hidden border border-glass-border">
                        <div class="bg-status-success h-full transition-all duration-1000" style="width: {(stats.completed_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>

                <div class="space-y-2">
                    <div class="flex justify-between text-xs font-bold">
                        <span class="text-text-secondary">PENDING</span>
                        <span class="text-text-primary">{Math.round((stats.pending_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-main rounded-full h-2 overflow-hidden border border-glass-border">
                        <div class="bg-status-warning h-full transition-all duration-1000" style="width: {(stats.pending_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>

                <div class="space-y-2">
                    <div class="flex justify-between text-xs font-bold">
                        <span class="text-text-secondary">CANCELLED</span>
                        <span class="text-text-primary">{Math.round((stats.cancelled_orders / (stats.total_orders || 1)) * 100)}%</span>
                    </div>
                    <div class="w-full bg-main rounded-full h-2 overflow-hidden border border-glass-border">
                        <div class="bg-status-error h-full transition-all duration-1000" style="width: {(stats.cancelled_orders / (stats.total_orders || 1)) * 100}%"></div>
                    </div>
                </div>
            </div>

            <div class="mt-8 pt-6 border-t border-glass-border text-center">
                <div class="text-sm font-bold text-text-secondary">AVG ORDER VALUE</div>
                <div class="text-2xl font-black text-primary">${(stats.average_order_value || 0).toFixed(2)}</div>
            </div>
        </div>
    </div>
</div>