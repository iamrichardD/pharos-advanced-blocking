import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  site: 'https://iamrichardd.com',
  base: '/pharos-advanced-blocking',
  integrations: [mdx()],
  vite: {
    plugins: [tailwindcss()],
  },
});
