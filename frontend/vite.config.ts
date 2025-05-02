import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig(({ }) => {

  return {
    server: {
      host: "::",
      port: 5173, // Explicitly set the port to 5173
      proxy: {
        // Proxy API requests to the backend server
        '/api': {
          target: 'http://localhost:8080', // Your backend server address
          changeOrigin: true, // Needed for virtual hosted sites
          // rewrite: (path) => path.replace(/^\/api/, ''), // Uncomment if you need to remove /api prefix
        },
      }
    },
    plugins: [
      react(),
    ].filter(Boolean),
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
  };
});
