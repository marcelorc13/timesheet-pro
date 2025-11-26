# Deploying Timesheet Pro to Render.com

This guide will walk you through deploying your timesheet-pro application to Render.com.

## Prerequisites

- A [Render.com](https://render.com) account (free tier available)
- Your code pushed to a Git repository (GitHub, GitLab, or Bitbucket)
- Basic understanding of environment variables

---

## Deployment Options

You have two options to deploy on Render:

### Option 1: Blueprint (Automated) - Recommended

This uses the `render.yaml` file for automatic setup.

### Option 2: Manual Setup

Set up each service individually through the Render dashboard.

---

## Option 1: Blueprint Deployment (Recommended)

### Step 1: Push Your Code

Make sure your code (including `Dockerfile` and `render.yaml`) is pushed to GitHub:

```bash
git add Dockerfile render.yaml
git commit -m "Add Docker and Render configuration"
git push origin main
```

### Step 2: Create New Blueprint

1. Go to [Render Dashboard](https://dashboard.render.com/)
2. Click **"New +"** → **"Blueprint"**
3. Connect your Git repository
4. Render will automatically detect `render.yaml` and show you the services it will create:
   - **Web Service**: `timesheet-pro`
   - **PostgreSQL Database**: `timesheet-pro-db`
5. Click **"Apply"**

### Step 3: Configure Environment Variables

The blueprint will automatically configure most environment variables, but you may need to add:

| Variable | Value | Description |
|----------|-------|-------------|
| `POSTGRES_URL` | Auto-generated | Database connection string (auto-linked) |
| `JWT_SECRET` | Auto-generated | Secret for JWT tokens |
| `GIN_MODE` | `release` | Gin framework mode |
| `PORT` | `8080` | Application port (Render auto-sets this) |

> [!NOTE]
> Render automatically injects the `PORT` environment variable. Your application should read from `os.Getenv("PORT")` if needed.

### Step 4: Run Database Migrations

After deployment, you'll need to run migrations:

1. Go to your web service in the Render dashboard
2. Click **"Shell"** tab
3. Run migrations (adjust based on your migration tool):

```bash
# If using goose
goose -dir ./migrations postgres "${POSTGRES_URL}" up
```

> [!IMPORTANT]
> Make sure your migrations directory is included in the Docker image. You may need to update the Dockerfile to copy migrations.

---

## Option 2: Manual Deployment

### Step 1: Create PostgreSQL Database

1. Go to [Render Dashboard](https://dashboard.render.com/)
2. Click **"New +"** → **"PostgreSQL"**
3. Configure:
   - **Name**: `timesheet-pro-db`
   - **Database**: `timesheet_pro`
   - **User**: `timesheet_user`
   - **Plan**: Free
4. Click **"Create Database"**
5. **Save the connection string** (Internal Database URL)

### Step 2: Create Web Service

1. Click **"New +"** → **"Web Service"**
2. Connect your Git repository
3. Configure:
   - **Name**: `timesheet-pro`
   - **Environment**: `Docker`
   - **Plan**: Free
   - **Branch**: `main` (or your default branch)

### Step 3: Add Environment Variables

In the web service settings, add these environment variables:

| Key | Value |
|-----|-------|
| `POSTGRES_URL` | Paste the Internal Database URL from Step 1 |
| `JWT_SECRET` | Generate a random secure string |
| `GIN_MODE` | `release` |

> [!TIP]
> To generate a secure JWT_SECRET:
> ```bash
> openssl rand -base64 32
> ```

### Step 4: Deploy

1. Click **"Create Web Service"**
2. Render will:
   - Clone your repository
   - Build the Docker image using your `Dockerfile`
   - Deploy the application
3. Wait for deployment to complete (check the logs)

### Step 5: Run Migrations

1. Go to your web service
2. Click **"Shell"** tab
3. Run your migration commands

---

## Updating Application Code in main.go

Your application needs to listen on the port that Render provides via the `PORT` environment variable.

Update `cmd/api/main.go` to read the port:

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/marcelorc13/timesheet-pro/internal/repository"
    "github.com/marcelorc13/timesheet-pro/internal/server"
    "github.com/marcelorc13/timesheet-pro/internal/server/api"
    "github.com/marcelorc13/timesheet-pro/internal/server/views"
    service "github.com/marcelorc13/timesheet-pro/internal/services"
)

func main() {
    _ = godotenv.Load()
    connString := os.Getenv("POSTGRES_URL")

    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    ctx := context.Background()

    db := repository.NewPool(ctx, connString)

    if err := db.Ping(ctx); err != nil {
        panic(err)
    }

    r := gin.Default()

    router := server.NewRouter(r)

    ur := repository.NewUserRepository(db)
    us := service.NewUserService(*ur)
    uh := api.NewUserHandler(*us)

    or := repository.NewOrganizationRepository(db)
    os := service.NewOrganizationService(*or, *ur)
    oh := api.NewOrganizationHandler(*os)

    // Timesheet setup
    tr := repository.NewTimesheetRepository(db)
    ts := service.NewTimesheetService(tr, or)
    th := api.NewTimesheetHandler(ts)

    // View handlers
    ovh := views.NewOrganizationViewHandler(*os, *us)
    tvh := views.NewTimesheetViewHandler(ts, os)
    pvh := views.NewProfileViewHandler(us)

    router.APIRoutes(*uh, *oh, *th)
    router.ViewsRoutes(*ovh, *tvh, *pvh, or)

    // Start server on dynamic port
    r.Run(fmt.Sprintf(":%s", port))
}
```

> [!WARNING]
> You'll need to update the `router.Start()` method or replace it with `r.Run()` as shown above to properly use the PORT environment variable.

---

## Including Migrations in Docker

If you need to copy migrations into the Docker image, update your `Dockerfile`:

```dockerfile
# Copy the entire project (including migrations)
COPY . .
```

The current Dockerfile already does this with `COPY . .`, so your migrations should be included.

---

## Post-Deployment

### Verify Deployment

1. Check the build logs for any errors
2. Visit your app URL: `https://timesheet-pro.onrender.com`
3. Check the application logs in the Render dashboard

### Monitor Your Application

- **Logs**: Available in the Render dashboard under your web service
- **Metrics**: CPU, memory usage, and request metrics
- **Alerts**: Set up in Render dashboard for downtime notifications

### Free Tier Limitations

> [!CAUTION]
> Render's free tier has some limitations:
> - Web services spin down after 15 minutes of inactivity
> - First request after spin-down may take 30-60 seconds
> - PostgreSQL free tier has limited storage (1GB)
> - Databases are deleted after 90 days on free tier

---

## Troubleshooting

### Application Won't Start

**Check logs** in the Render dashboard for specific errors.

Common issues:
- Missing environment variables
- Database connection failures
- Port binding issues

### Database Connection Issues

Verify:
1. `POSTGRES_URL` is set correctly
2. Database is running and accessible
3. Connection string format is correct

### Build Failures

Check:
1. Dockerfile syntax
2. Go module dependencies
3. Build logs for specific errors

---

## Next Steps

- Set up custom domain in Render dashboard
- Configure health checks
- Set up monitoring and alerting
- Consider upgrading to paid plans for production workloads

## Additional Resources

- [Render Documentation](https://render.com/docs)
- [Render PostgreSQL Guide](https://render.com/docs/databases)
- [Render Docker Guide](https://render.com/docs/docker)
