# GitHub Pages Configuration Checklist

## ✅ Repository Settings to Verify

1. **GitHub Pages Source:**
   - Go to your repo: Settings → Pages
   - Source should be set to "GitHub Actions"
   - Not "Deploy from a branch"

2. **Required Secrets (Settings → Secrets and variables → Actions):**
   ```
   AZURE_CREDENTIALS          # Azure service principal JSON
   AZURE_OPENAI_ENDPOINT      # https://your-resource.openai.azure.com/
   AZURE_OPENAI_API_KEY       # Your Azure OpenAI key
   AZURE_OPENAI_DEPLOYMENT_NAME # gpt-4 (or your model name)
   ```

3. **Repository Permissions:**
   - Settings → Actions → General
   - Workflow permissions: "Read and write permissions" ✅
   - Allow GitHub Actions to create and approve pull requests ✅

## 🚀 How to Enable GitHub Pages

### Step 1: Configure GitHub Pages Source
1. Go to your repository on GitHub
2. Click **Settings** tab
3. Scroll down to **Pages** section
4. Under **Source**, select **"GitHub Actions"**
5. Click **Save**

### Step 2: Add Required Secrets
1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Click **"New repository secret"** for each:

**AZURE_CREDENTIALS** (Azure Service Principal):
```json
{
  "clientId": "your-client-id",
  "clientSecret": "your-client-secret", 
  "subscriptionId": "your-subscription-id",
  "tenantId": "your-tenant-id"
}
```

**AZURE_OPENAI_ENDPOINT**:
```
https://your-openai-resource.openai.azure.com/
```

**AZURE_OPENAI_API_KEY**:
```
your_azure_openai_api_key_here
```

**AZURE_OPENAI_DEPLOYMENT_NAME**:
```
gpt-4
```

### Step 3: Push to Main Branch
Once configured, push any commit to the `main` branch to trigger deployment:

```bash
git add .
git commit -m "Configure GitHub Pages deployment"
git push origin main
```

### Step 4: Check Deployment Status
- Go to **Actions** tab in your repo
- Watch the "Deploy Reading Tracker with azd" workflow
- When complete, your app will be live at:
  `https://yourusername.github.io/reading-app`

## 🔧 Testing GitHub Pages

After the workflow completes:

1. **Check deployment URL**: Actions → Latest workflow run → deploy-frontend job → Deploy to GitHub Pages step
2. **Verify app loads**: Visit your GitHub Pages URL
3. **Test functionality**: 
   - Dashboard should show your 26/90 reading progress
   - All stats should load correctly
   - AI chat should work (if Azure backend deployed)

## 🚨 Troubleshooting

**If deployment fails:**
- Check Actions tab for error messages
- Verify all secrets are correctly set
- Ensure GitHub Pages source is set to "GitHub Actions"
- Make sure repository is public (or you have GitHub Pro/Team for private Pages)

**If page loads but doesn't work:**
- Check browser console for errors
- Verify `books.js` file is loading correctly
- Test with hardcoded data first

## 📱 Expected Result

Once configured, your app will:
- ✅ Auto-deploy on every push to main
- ✅ Show 26/90 reading progress 
- ✅ Display all your book statistics
- ✅ Connect to Azure Functions backend
- ✅ Work on mobile and desktop
- ✅ Update automatically when you add books

Your live URL will be: `https://yourusername.github.io/reading-app`
