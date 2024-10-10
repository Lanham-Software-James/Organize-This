import { purgeCss } from 'vite-plugin-tailwind-purgecss';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { configDefaults, coverageConfigDefaults } from 'vitest/config'

export default defineConfig({
	plugins: [sveltekit(), purgeCss()],
	test: {
		exclude: [
			...configDefaults.exclude,
			'playwright-tests/*'
		],
		coverage: {
			provider: 'istanbul',
			exclude: ['*.config*', ...coverageConfigDefaults.exclude]
		},
		globals: true,
		environment: 'jsdom',
	}
});
