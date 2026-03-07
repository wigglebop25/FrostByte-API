<script lang="ts">
    /**
     * Revenue Chart Component
     *
     * Renders an ApexCharts area chart for revenue data. Automatically
     * adapts colors and styling to the current light/dark theme.
     */
    import { onMount } from "svelte";
    import { theme } from "$lib/stores/theme";
    import { browser } from "$app/environment";
    
    export let data: number[] = [];
    export let labels: string[] = [];
    let chartNode: HTMLElement;
    let chart: any;

    function getResolvedTheme(): 'light' | 'dark' {
        if (!browser) return 'light';
        return document.documentElement.classList.contains('dark') ? 'dark' : 'light';
    }

    function getChartOptions(mode: 'light' | 'dark') {
        return {
            chart: {
                type: "area",
                height: 300,
                toolbar: { show: false },
                background: "transparent",
                animations: { enabled: true, easing: 'easeinout', speed: 800 },
                fontFamily: 'Inter, system-ui, sans-serif'
            },
            series: [{ name: "Revenue", data: data }],
            colors: [mode === 'dark' ? '#FF6B6B' : '#E63946'],
            xaxis: {
                categories: labels,
                labels: { 
                    style: { 
                        colors: mode === 'dark' ? '#9CA3AF' : '#6B7280',
                        fontWeight: 600, 
                        fontSize: '11px' 
                    } 
                },
                axisBorder: { show: false },
                axisTicks: { show: false }
            },
            yaxis: {
                labels: { 
                    style: { 
                        colors: mode === 'dark' ? '#9CA3AF' : '#6B7280',
                        fontWeight: 600 
                    },
                    formatter: (val: number) => `$${val.toLocaleString()}`
                }
            },
            grid: {
                borderColor: mode === 'dark' ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.06)',
                strokeDashArray: 4
            },
            theme: { mode },
            stroke: { curve: "smooth", width: 3 },
            fill: {
                type: "gradient",
                gradient: {
                    shadeIntensity: 1,
                    opacityFrom: 0.35,
                    opacityTo: 0.02,
                    stops: [0, 90, 100]
                }
            },
            dataLabels: { enabled: false },
            markers: { 
                size: 4, 
                strokeWidth: 2, 
                strokeColors: mode === 'dark' ? '#16161D' : '#FFFFFF',
                hover: { sizeOffset: 3 }
            },
            tooltip: {
                theme: mode,
                style: { fontSize: '13px' },
                y: { formatter: (val: number) => `$${val.toLocaleString()}` }
            }
        };
    }

    onMount(() => {
        (async () => {
            const module = await import("apexcharts");
            const ApexCharts = module.default;
            const mode = getResolvedTheme();

            if (chartNode) {
                chart = new ApexCharts(chartNode, getChartOptions(mode));
                chart.render();
            }
        })();

        return () => {
            if (chart) chart.destroy();
        };
    });

    /** Handles reactive updates when chart data or labels change */
    $: if (chart && data && labels) {
        chart.updateOptions({
            xaxis: { categories: labels },
            series: [{ data: data }]
        });
    }

    /** Re-renders chart styling when the application theme changes */
    $: if (chart && $theme) {
        const mode = getResolvedTheme();
        chart.updateOptions(getChartOptions(mode));
    }
</script>

<div bind:this={chartNode}></div>
