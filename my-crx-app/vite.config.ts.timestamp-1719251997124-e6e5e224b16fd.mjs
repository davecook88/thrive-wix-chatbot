// vite.config.ts
import { defineConfig } from "file:///Users/personal/wix-thrive-placement-chatbot/my-crx-app/node_modules/vite/dist/node/index.js";
import { crx } from "file:///Users/personal/wix-thrive-placement-chatbot/my-crx-app/node_modules/@crxjs/vite-plugin/dist/index.mjs";
import { svelte } from "file:///Users/personal/wix-thrive-placement-chatbot/my-crx-app/node_modules/@sveltejs/vite-plugin-svelte/src/index.js";
import path from "path";
import sveltePreprocess from "file:///Users/personal/wix-thrive-placement-chatbot/my-crx-app/node_modules/svelte-preprocess/dist/index.js";

// src/manifest.ts
import { defineManifest } from "file:///Users/personal/wix-thrive-placement-chatbot/my-crx-app/node_modules/@crxjs/vite-plugin/dist/index.mjs";

// package.json
var package_default = {
  name: "my-crx-app",
  displayName: "my-crx-app",
  version: "0.0.0",
  author: "**",
  description: "",
  type: "module",
  license: "MIT",
  keywords: [
    "chrome-extension",
    "svelte",
    "vite",
    "create-chrome-ext"
  ],
  engines: {
    node: ">=14.18.0"
  },
  scripts: {
    dev: "vite",
    build: "vite build",
    preview: "vite preview",
    fmt: "prettier --write '**/*.{svelte,ts,json,css,scss,md}'",
    zip: "npm run build && node src/zip.js"
  },
  devDependencies: {
    "@crxjs/vite-plugin": "^2.0.0-beta.19",
    "@sveltejs/vite-plugin-svelte": "2.4.6",
    "@types/chrome": "^0.0.246",
    gulp: "^4.0.2",
    "gulp-zip": "^6.0.0",
    prettier: "^3.0.3",
    "prettier-plugin-svelte": "^3.0.3",
    svelte: "^4.2.1",
    "svelte-preprocess": "^5.0.4",
    tslib: "^2.6.2",
    typescript: "^5.2.2",
    vite: "^4.4.11"
  }
};

// src/manifest.ts
var isDev = process.env.NODE_ENV == "development";
var manifest_default = defineManifest({
  name: `${package_default.displayName || package_default.name}${isDev ? ` \u27A1\uFE0F Dev` : ""}`,
  description: package_default.description,
  version: package_default.version,
  manifest_version: 3,
  icons: {
    16: "img/logo-16.png",
    32: "img/logo-34.png",
    48: "img/logo-48.png",
    128: "img/logo-128.png"
  },
  action: {
    default_popup: "popup.html",
    default_icon: "img/logo-48.png"
  },
  options_page: "options.html",
  devtools_page: "devtools.html",
  background: {
    service_worker: "src/background/index.ts",
    type: "module"
  },
  content_scripts: [
    {
      matches: ["http://*/*", "https://*/*"],
      js: ["src/contentScript/index.ts"]
    }
  ],
  side_panel: {
    default_path: "sidepanel.html"
  },
  web_accessible_resources: [
    {
      resources: ["img/logo-16.png", "img/logo-34.png", "img/logo-48.png", "img/logo-128.png"],
      matches: []
    }
  ],
  permissions: ["sidePanel", "storage"],
  chrome_url_overrides: {
    newtab: "newtab.html"
  }
});

// vite.config.ts
var __vite_injected_original_dirname = "/Users/personal/wix-thrive-placement-chatbot/my-crx-app";
var vite_config_default = defineConfig(({ mode }) => {
  const production = mode === "production";
  return {
    build: {
      emptyOutDir: true,
      outDir: "build",
      rollupOptions: {
        output: {
          chunkFileNames: "assets/chunk-[hash].js"
        }
      }
    },
    plugins: [
      crx({ manifest: manifest_default }),
      svelte({
        compilerOptions: {
          dev: !production
        },
        preprocess: sveltePreprocess()
      })
    ],
    resolve: {
      alias: {
        "@": path.resolve(__vite_injected_original_dirname, "src")
      }
    }
  };
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiLCAic3JjL21hbmlmZXN0LnRzIiwgInBhY2thZ2UuanNvbiJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIi9Vc2Vycy9wZXJzb25hbC93aXgtdGhyaXZlLXBsYWNlbWVudC1jaGF0Ym90L215LWNyeC1hcHBcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIi9Vc2Vycy9wZXJzb25hbC93aXgtdGhyaXZlLXBsYWNlbWVudC1jaGF0Ym90L215LWNyeC1hcHAvdml0ZS5jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL1VzZXJzL3BlcnNvbmFsL3dpeC10aHJpdmUtcGxhY2VtZW50LWNoYXRib3QvbXktY3J4LWFwcC92aXRlLmNvbmZpZy50c1wiO2ltcG9ydCB7IGRlZmluZUNvbmZpZyB9IGZyb20gJ3ZpdGUnXG5pbXBvcnQgeyBjcnggfSBmcm9tICdAY3J4anMvdml0ZS1wbHVnaW4nXG5pbXBvcnQgeyBzdmVsdGUgfSBmcm9tICdAc3ZlbHRlanMvdml0ZS1wbHVnaW4tc3ZlbHRlJ1xuaW1wb3J0IHBhdGggZnJvbSAncGF0aCdcbmltcG9ydCBzdmVsdGVQcmVwcm9jZXNzIGZyb20gJ3N2ZWx0ZS1wcmVwcm9jZXNzJ1xuaW1wb3J0IG1hbmlmZXN0IGZyb20gJy4vc3JjL21hbmlmZXN0J1xuXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoKHsgbW9kZSB9KSA9PiB7XG4gIGNvbnN0IHByb2R1Y3Rpb24gPSBtb2RlID09PSAncHJvZHVjdGlvbidcblxuICByZXR1cm4ge1xuICAgIGJ1aWxkOiB7XG4gICAgICBlbXB0eU91dERpcjogdHJ1ZSxcbiAgICAgIG91dERpcjogJ2J1aWxkJyxcbiAgICAgIHJvbGx1cE9wdGlvbnM6IHtcbiAgICAgICAgb3V0cHV0OiB7XG4gICAgICAgICAgY2h1bmtGaWxlTmFtZXM6ICdhc3NldHMvY2h1bmstW2hhc2hdLmpzJyxcbiAgICAgICAgfSxcbiAgICAgIH0sXG4gICAgfSxcbiAgICBwbHVnaW5zOiBbXG4gICAgICBjcngoeyBtYW5pZmVzdCB9KSxcbiAgICAgIHN2ZWx0ZSh7XG4gICAgICAgIGNvbXBpbGVyT3B0aW9uczoge1xuICAgICAgICAgIGRldjogIXByb2R1Y3Rpb24sXG4gICAgICAgIH0sXG4gICAgICAgIHByZXByb2Nlc3M6IHN2ZWx0ZVByZXByb2Nlc3MoKSxcbiAgICAgIH0pLFxuICAgIF0sXG4gICAgcmVzb2x2ZToge1xuICAgICAgYWxpYXM6IHtcbiAgICAgICAgJ0AnOiBwYXRoLnJlc29sdmUoX19kaXJuYW1lLCAnc3JjJyksXG4gICAgICB9LFxuICAgIH0sXG4gIH1cbn0pXG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIi9Vc2Vycy9wZXJzb25hbC93aXgtdGhyaXZlLXBsYWNlbWVudC1jaGF0Ym90L215LWNyeC1hcHAvc3JjXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvVXNlcnMvcGVyc29uYWwvd2l4LXRocml2ZS1wbGFjZW1lbnQtY2hhdGJvdC9teS1jcngtYXBwL3NyYy9tYW5pZmVzdC50c1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vVXNlcnMvcGVyc29uYWwvd2l4LXRocml2ZS1wbGFjZW1lbnQtY2hhdGJvdC9teS1jcngtYXBwL3NyYy9tYW5pZmVzdC50c1wiO2ltcG9ydCB7IGRlZmluZU1hbmlmZXN0IH0gZnJvbSAnQGNyeGpzL3ZpdGUtcGx1Z2luJ1xuaW1wb3J0IHBhY2thZ2VEYXRhIGZyb20gJy4uL3BhY2thZ2UuanNvbidcblxuLy9AdHMtaWdub3JlXG5jb25zdCBpc0RldiA9IHByb2Nlc3MuZW52Lk5PREVfRU5WID09ICdkZXZlbG9wbWVudCdcblxuZXhwb3J0IGRlZmF1bHQgZGVmaW5lTWFuaWZlc3Qoe1xuICBuYW1lOiBgJHtwYWNrYWdlRGF0YS5kaXNwbGF5TmFtZSB8fCBwYWNrYWdlRGF0YS5uYW1lfSR7aXNEZXYgPyBgIFx1MjdBMVx1RkUwRiBEZXZgIDogJyd9YCxcbiAgZGVzY3JpcHRpb246IHBhY2thZ2VEYXRhLmRlc2NyaXB0aW9uLFxuICB2ZXJzaW9uOiBwYWNrYWdlRGF0YS52ZXJzaW9uLFxuICBtYW5pZmVzdF92ZXJzaW9uOiAzLFxuICBpY29uczoge1xuICAgIDE2OiAnaW1nL2xvZ28tMTYucG5nJyxcbiAgICAzMjogJ2ltZy9sb2dvLTM0LnBuZycsXG4gICAgNDg6ICdpbWcvbG9nby00OC5wbmcnLFxuICAgIDEyODogJ2ltZy9sb2dvLTEyOC5wbmcnLFxuICB9LFxuICBhY3Rpb246IHtcbiAgICBkZWZhdWx0X3BvcHVwOiAncG9wdXAuaHRtbCcsXG4gICAgZGVmYXVsdF9pY29uOiAnaW1nL2xvZ28tNDgucG5nJyxcbiAgfSxcbiAgb3B0aW9uc19wYWdlOiAnb3B0aW9ucy5odG1sJyxcbiAgZGV2dG9vbHNfcGFnZTogJ2RldnRvb2xzLmh0bWwnLFxuICBiYWNrZ3JvdW5kOiB7XG4gICAgc2VydmljZV93b3JrZXI6ICdzcmMvYmFja2dyb3VuZC9pbmRleC50cycsXG4gICAgdHlwZTogJ21vZHVsZScsXG4gIH0sXG4gIGNvbnRlbnRfc2NyaXB0czogW1xuICAgIHtcbiAgICAgIG1hdGNoZXM6IFsnaHR0cDovLyovKicsICdodHRwczovLyovKiddLFxuICAgICAganM6IFsnc3JjL2NvbnRlbnRTY3JpcHQvaW5kZXgudHMnXSxcbiAgICB9LFxuICBdLFxuICBzaWRlX3BhbmVsOiB7XG4gICAgZGVmYXVsdF9wYXRoOiAnc2lkZXBhbmVsLmh0bWwnLFxuICB9LFxuICB3ZWJfYWNjZXNzaWJsZV9yZXNvdXJjZXM6IFtcbiAgICB7XG4gICAgICByZXNvdXJjZXM6IFsnaW1nL2xvZ28tMTYucG5nJywgJ2ltZy9sb2dvLTM0LnBuZycsICdpbWcvbG9nby00OC5wbmcnLCAnaW1nL2xvZ28tMTI4LnBuZyddLFxuICAgICAgbWF0Y2hlczogW10sXG4gICAgfSxcbiAgXSxcbiAgcGVybWlzc2lvbnM6IFsnc2lkZVBhbmVsJywgJ3N0b3JhZ2UnXSxcbiAgY2hyb21lX3VybF9vdmVycmlkZXM6IHtcbiAgICBuZXd0YWI6ICduZXd0YWIuaHRtbCcsXG4gIH0sXG59KVxuIiwgIntcbiAgXCJuYW1lXCI6IFwibXktY3J4LWFwcFwiLFxuICBcImRpc3BsYXlOYW1lXCI6IFwibXktY3J4LWFwcFwiLFxuICBcInZlcnNpb25cIjogXCIwLjAuMFwiLFxuICBcImF1dGhvclwiOiBcIioqXCIsXG4gIFwiZGVzY3JpcHRpb25cIjogXCJcIixcbiAgXCJ0eXBlXCI6IFwibW9kdWxlXCIsXG4gIFwibGljZW5zZVwiOiBcIk1JVFwiLFxuICBcImtleXdvcmRzXCI6IFtcbiAgICBcImNocm9tZS1leHRlbnNpb25cIixcbiAgICBcInN2ZWx0ZVwiLFxuICAgIFwidml0ZVwiLFxuICAgIFwiY3JlYXRlLWNocm9tZS1leHRcIlxuICBdLFxuICBcImVuZ2luZXNcIjoge1xuICAgIFwibm9kZVwiOiBcIj49MTQuMTguMFwiXG4gIH0sXG4gIFwic2NyaXB0c1wiOiB7XG4gICAgXCJkZXZcIjogXCJ2aXRlXCIsXG4gICAgXCJidWlsZFwiOiBcInZpdGUgYnVpbGRcIixcbiAgICBcInByZXZpZXdcIjogXCJ2aXRlIHByZXZpZXdcIixcbiAgICBcImZtdFwiOiBcInByZXR0aWVyIC0td3JpdGUgJyoqLyoue3N2ZWx0ZSx0cyxqc29uLGNzcyxzY3NzLG1kfSdcIixcbiAgICBcInppcFwiOiBcIm5wbSBydW4gYnVpbGQgJiYgbm9kZSBzcmMvemlwLmpzXCJcbiAgfSxcbiAgXCJkZXZEZXBlbmRlbmNpZXNcIjoge1xuICAgIFwiQGNyeGpzL3ZpdGUtcGx1Z2luXCI6IFwiXjIuMC4wLWJldGEuMTlcIixcbiAgICBcIkBzdmVsdGVqcy92aXRlLXBsdWdpbi1zdmVsdGVcIjogXCIyLjQuNlwiLFxuICAgIFwiQHR5cGVzL2Nocm9tZVwiOiBcIl4wLjAuMjQ2XCIsXG4gICAgXCJndWxwXCI6IFwiXjQuMC4yXCIsXG4gICAgXCJndWxwLXppcFwiOiBcIl42LjAuMFwiLFxuICAgIFwicHJldHRpZXJcIjogXCJeMy4wLjNcIixcbiAgICBcInByZXR0aWVyLXBsdWdpbi1zdmVsdGVcIjogXCJeMy4wLjNcIixcbiAgICBcInN2ZWx0ZVwiOiBcIl40LjIuMVwiLFxuICAgIFwic3ZlbHRlLXByZXByb2Nlc3NcIjogXCJeNS4wLjRcIixcbiAgICBcInRzbGliXCI6IFwiXjIuNi4yXCIsXG4gICAgXCJ0eXBlc2NyaXB0XCI6IFwiXjUuMi4yXCIsXG4gICAgXCJ2aXRlXCI6IFwiXjQuNC4xMVwiXG4gIH1cbn0iXSwKICAibWFwcGluZ3MiOiAiO0FBQXVWLFNBQVMsb0JBQW9CO0FBQ3BYLFNBQVMsV0FBVztBQUNwQixTQUFTLGNBQWM7QUFDdkIsT0FBTyxVQUFVO0FBQ2pCLE9BQU8sc0JBQXNCOzs7QUNKZ1UsU0FBUyxzQkFBc0I7OztBQ0E1WDtBQUFBLEVBQ0UsTUFBUTtBQUFBLEVBQ1IsYUFBZTtBQUFBLEVBQ2YsU0FBVztBQUFBLEVBQ1gsUUFBVTtBQUFBLEVBQ1YsYUFBZTtBQUFBLEVBQ2YsTUFBUTtBQUFBLEVBQ1IsU0FBVztBQUFBLEVBQ1gsVUFBWTtBQUFBLElBQ1Y7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxFQUNGO0FBQUEsRUFDQSxTQUFXO0FBQUEsSUFDVCxNQUFRO0FBQUEsRUFDVjtBQUFBLEVBQ0EsU0FBVztBQUFBLElBQ1QsS0FBTztBQUFBLElBQ1AsT0FBUztBQUFBLElBQ1QsU0FBVztBQUFBLElBQ1gsS0FBTztBQUFBLElBQ1AsS0FBTztBQUFBLEVBQ1Q7QUFBQSxFQUNBLGlCQUFtQjtBQUFBLElBQ2pCLHNCQUFzQjtBQUFBLElBQ3RCLGdDQUFnQztBQUFBLElBQ2hDLGlCQUFpQjtBQUFBLElBQ2pCLE1BQVE7QUFBQSxJQUNSLFlBQVk7QUFBQSxJQUNaLFVBQVk7QUFBQSxJQUNaLDBCQUEwQjtBQUFBLElBQzFCLFFBQVU7QUFBQSxJQUNWLHFCQUFxQjtBQUFBLElBQ3JCLE9BQVM7QUFBQSxJQUNULFlBQWM7QUFBQSxJQUNkLE1BQVE7QUFBQSxFQUNWO0FBQ0Y7OztBRGxDQSxJQUFNLFFBQVEsUUFBUSxJQUFJLFlBQVk7QUFFdEMsSUFBTyxtQkFBUSxlQUFlO0FBQUEsRUFDNUIsTUFBTSxHQUFHLGdCQUFZLGVBQWUsZ0JBQVksSUFBSSxHQUFHLFFBQVEsc0JBQVksRUFBRTtBQUFBLEVBQzdFLGFBQWEsZ0JBQVk7QUFBQSxFQUN6QixTQUFTLGdCQUFZO0FBQUEsRUFDckIsa0JBQWtCO0FBQUEsRUFDbEIsT0FBTztBQUFBLElBQ0wsSUFBSTtBQUFBLElBQ0osSUFBSTtBQUFBLElBQ0osSUFBSTtBQUFBLElBQ0osS0FBSztBQUFBLEVBQ1A7QUFBQSxFQUNBLFFBQVE7QUFBQSxJQUNOLGVBQWU7QUFBQSxJQUNmLGNBQWM7QUFBQSxFQUNoQjtBQUFBLEVBQ0EsY0FBYztBQUFBLEVBQ2QsZUFBZTtBQUFBLEVBQ2YsWUFBWTtBQUFBLElBQ1YsZ0JBQWdCO0FBQUEsSUFDaEIsTUFBTTtBQUFBLEVBQ1I7QUFBQSxFQUNBLGlCQUFpQjtBQUFBLElBQ2Y7QUFBQSxNQUNFLFNBQVMsQ0FBQyxjQUFjLGFBQWE7QUFBQSxNQUNyQyxJQUFJLENBQUMsNEJBQTRCO0FBQUEsSUFDbkM7QUFBQSxFQUNGO0FBQUEsRUFDQSxZQUFZO0FBQUEsSUFDVixjQUFjO0FBQUEsRUFDaEI7QUFBQSxFQUNBLDBCQUEwQjtBQUFBLElBQ3hCO0FBQUEsTUFDRSxXQUFXLENBQUMsbUJBQW1CLG1CQUFtQixtQkFBbUIsa0JBQWtCO0FBQUEsTUFDdkYsU0FBUyxDQUFDO0FBQUEsSUFDWjtBQUFBLEVBQ0Y7QUFBQSxFQUNBLGFBQWEsQ0FBQyxhQUFhLFNBQVM7QUFBQSxFQUNwQyxzQkFBc0I7QUFBQSxJQUNwQixRQUFRO0FBQUEsRUFDVjtBQUNGLENBQUM7OztBRDlDRCxJQUFNLG1DQUFtQztBQU96QyxJQUFPLHNCQUFRLGFBQWEsQ0FBQyxFQUFFLEtBQUssTUFBTTtBQUN4QyxRQUFNLGFBQWEsU0FBUztBQUU1QixTQUFPO0FBQUEsSUFDTCxPQUFPO0FBQUEsTUFDTCxhQUFhO0FBQUEsTUFDYixRQUFRO0FBQUEsTUFDUixlQUFlO0FBQUEsUUFDYixRQUFRO0FBQUEsVUFDTixnQkFBZ0I7QUFBQSxRQUNsQjtBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUEsSUFDQSxTQUFTO0FBQUEsTUFDUCxJQUFJLEVBQUUsMkJBQVMsQ0FBQztBQUFBLE1BQ2hCLE9BQU87QUFBQSxRQUNMLGlCQUFpQjtBQUFBLFVBQ2YsS0FBSyxDQUFDO0FBQUEsUUFDUjtBQUFBLFFBQ0EsWUFBWSxpQkFBaUI7QUFBQSxNQUMvQixDQUFDO0FBQUEsSUFDSDtBQUFBLElBQ0EsU0FBUztBQUFBLE1BQ1AsT0FBTztBQUFBLFFBQ0wsS0FBSyxLQUFLLFFBQVEsa0NBQVcsS0FBSztBQUFBLE1BQ3BDO0FBQUEsSUFDRjtBQUFBLEVBQ0Y7QUFDRixDQUFDOyIsCiAgIm5hbWVzIjogW10KfQo=
