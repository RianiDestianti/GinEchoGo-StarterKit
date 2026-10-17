module.exports = {
  apps: [
    {
      name: "gin-skeleton",
      script: "bin/app",
      interpreter: "none",
      watch: false,
      autorestart: true,
      max_memory_restart: "500M"
    }
  ]
}
