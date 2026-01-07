# Branch Ruleset Quick Checklist

Use this checklist when configuring branch protection rules in GitHub Settings → Rules → Rulesets.

## For `main` Branch

### Essential Rules ✅
- [ ] **Require a pull request before merging**
  - [ ] Required approvals: 1
  - [ ] Dismiss stale reviews when new commits are pushed
  
- [ ] **Block force pushes**

- [ ] **Restrict deletions**

### Highly Recommended ⚠️
- [ ] **Require status checks to pass**
  - [ ] Require branches to be up to date before merging
  - [ ] Add specific status checks once CI/CD is configured
  
- [ ] **Require linear history**

- [ ] **Require signed commits**

- [ ] **Automatically request Copilot code review**

### Configure When Ready 🔧
- [ ] **Require code scanning results** (after setting up CodeQL)

- [ ] **Require code quality results** (after setting up linters)

- [ ] **Manage static analysis tools in Copilot code review**

### Usually Not Needed ⛔
- [ ] Restrict creations (only if you want to control who can create branches)
- [ ] Restrict updates (very restrictive, usually not needed)
- [ ] Require deployments to succeed (only if you have deployment pipelines)

---

## Recommended Minimal Configuration

For a new API repository, start with these 4 essential rules:

1. ✅ **Require a pull request before merging** (1 approval)
2. ✅ **Block force pushes**
3. ✅ **Restrict deletions**
4. ✅ **Automatically request Copilot code review**

Then add these as your CI/CD matures:

5. **Require status checks to pass** (once you have tests/build)
6. **Require linear history** (for clean git history)
7. **Require signed commits** (for security)

---

## Quick Start Commands

### 1. Create the Ruleset
```
Settings → Rules → Rulesets → New ruleset → New branch ruleset
```

### 2. Configure Basic Settings
- **Ruleset Name:** `main-branch-protection`
- **Enforcement status:** Active
- **Target branches:** Add target → Include default branch (or specify `main`)

### 3. Apply Rules
Enable the checkboxes for the rules listed above under "Essential Rules"

### 4. Save
Click "Create" to activate the ruleset

---

See [BRANCH_PROTECTION.md](../BRANCH_PROTECTION.md) for detailed explanations.
