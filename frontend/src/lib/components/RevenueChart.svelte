<script lang="ts">
    import { onMount } from "svelte";
    
    export let data: number[] = [];
    export let labels: string[] = [];
    let chartNode: HTMLElement;
    let chart: any;

    onMount(() => {
        (async () => {
            const module = await import("apexcharts");
            const ApexCharts = module.default;

            const options = {
                chart: {
                    type: "area",
                    height: 350,
                    toolbar: { show: false },
                    background: "transparent",
                    animations: { enabled: true, easing: 'easeinout', speed: 800 }
                },
                series: [{ name: "Revenue", data: data }],
                xaxis: {
                    categories: labels,
                    labels: { style: { colors: "var(--color-text-secondary)", fontWeight: 600, fontSize: '10px' } },
                    axisBorder: { show: false },
                    axisTicks: { show: false }
                },
                yaxis: {
                    labels: { 
                        style: { colors: "var(--color-text-secondary)", fontWeight: 600 },
                        formatter: (val: number) => `$${val.toLocaleString()}`
                    }
                },
                grid: {
                    borderColor: "var(--glass-border)",
                    strokeDashArray: 4
                },
                theme: { mode: "light" }, // Apex handles themes, we use vars mostly
                stroke: { curve: "smooth", width: 4, colors: ["var(--color-primary)"] },
                fill: {
                    type: "gradient",
                    gradient: {
                        shadeIntensity: 1,
                        opacityFrom: 0.4,
                        opacityTo: 0.05,
                        stops: [0, 90, 100]
                    }
                },
                dataLabels: { enabled: false },
                markers: { size: 5, colors: ["var(--color-primary)"], strokeWidth: 3, strokeColors: "#fff" }
            };

            if (chartNode) {
                chart = new ApexCharts(chartNode, options);
                chart.render();
            }
        })();

        return () => {
            if (chart) chart.destroy();
        };
    });

    // Reactive Updates
    $: if (chart && data && labels) {
        chart.updateOptions({
            xaxis: { categories: labels },
            series: [{ data: data }]
        });
    }
</script>

<div class="glass-card p-6 border-glass-border shadow-xl">
    <div class="flex items-center justify-between mb-6">
        <div>
            <h3 class="text-lg font-semibold text-text-primary">Weekly Revenue</h3>
            <p class="text-sm text-text-secondary">Live performance metrics</p>
        </div>
        <div class="px-3 py-1 rounded-full bg-status-info/10 text-status-info text-xs font-medium border border-status-info/20">
            Live Updates
        </div>
    </div>
    <div bind:this={chartNode}></div>
</div>
