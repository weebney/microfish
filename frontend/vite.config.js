// vite.config.js
import { defineConfig } from 'vite'
import { viteStaticCopy } from 'vite-plugin-static-copy'

export default defineConfig({
	plugins: [
		viteStaticCopy({
			targets: [
				{
					src: 'index.html',
					dest: ''
				}, {
					src: 'src/',
					dest: ''
				}, {
					src: 'wailsjs/',
					dest: ''
				}
			]
		})
	]
})