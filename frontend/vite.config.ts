import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'https://frostbyte-api.southeastasia.cloudapp.azure.com',
				changeOrigin: true,
				secure: false, // Accept self-signed certs if needed
				rewrite: (path) => path.replace(/^\/api/, '/api/v1')
			},
            '/ws': {
                target: 'wss://frostbyte-api.southeastasia.cloudapp.azure.com',
                changeOrigin: true,
                ws: true,
                secure: false
            }
		}
	}
});