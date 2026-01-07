# Branch Protection Ruleset for FrostByte-API

This document outlines the recommended branch protection rules for the FrostByte-API repository. These rules help maintain code quality, security, and collaboration standards.

## Recommended Ruleset Configuration

### Essential Rules (Highly Recommended)

#### 1. **Require a pull request before merging** ✅
**Status:** ENABLE

Require all commits be made to a non-target branch and submitted via a pull request before they can be merged.

**Configuration:**
- Required approving reviews: 1
- Dismiss stale pull request approvals when new commits are pushed: Yes
- Require review from Code Owners: Optional (if CODEOWNERS file exists)

**Rationale:** Ensures code review happens before changes are merged, improving code quality and knowledge sharing.

---

#### 2. **Require status checks to pass** ✅
**Status:** ENABLE

Choose which status checks must pass before the ref is updated. When enabled, commits must first be pushed to another ref where the checks pass.

**Configuration:**
- Require branches to be up to date before merging: Yes
- Status checks to require: (Add as you set up CI/CD)
  - Build
  - Tests
  - Linting

**Rationale:** Prevents broken code from being merged into protected branches.

---

#### 3. **Block force pushes** ✅
**Status:** ENABLE

Prevent users with push access from force pushing to refs.

**Rationale:** Protects commit history and prevents accidental data loss on main branches.

---

#### 4. **Restrict deletions** ✅
**Status:** ENABLE

Only allow users with bypass permission to delete matching refs.

**Rationale:** Prevents accidental deletion of important branches like main/master.

---

### Recommended Rules (Best Practices)

#### 5. **Require linear history** ⚠️
**Status:** ENABLE (Recommended)

Prevent merge commits from being pushed to matching refs.

**Configuration:**
- Use "Squash and merge" or "Rebase and merge" strategies

**Rationale:** Maintains a clean, linear git history that's easier to understand and debug.

---

#### 6. **Require signed commits** 🔐
**Status:** ENABLE (Recommended for security)

Commits pushed to matching refs must have verified signatures.

**Rationale:** Ensures commit authenticity and prevents impersonation.

---

#### 7. **Automatically request Copilot code review** 🤖
**Status:** ENABLE (Recommended)

Request Copilot code review for new pull requests automatically if the author has access to Copilot code review.

**Rationale:** Provides additional automated code review feedback to catch potential issues early.

---

### Optional Rules (Configure Based on Project Needs)

#### 8. **Restrict updates** ⚠️
**Status:** OPTIONAL

Only allow users with bypass permission to update matching refs.

**Use Case:** For extremely sensitive branches or release branches.

**Rationale:** Most projects don't need this; it's very restrictive.

---

#### 9. **Restrict creations** ⚠️
**Status:** OPTIONAL

Only allow users with bypass permission to create matching refs.

**Use Case:** For controlling release branch creation or preventing branch sprawl.

**Rationale:** Usually not needed for most development workflows.

---

#### 10. **Require deployments to succeed** 🚀
**Status:** OPTIONAL

Choose which environments must be successfully deployed to before refs can be pushed into a ref that matches this rule.

**Configuration:**
- Staging environment (if applicable)
- Test environment (if applicable)

**Use Case:** For production branches where you want to ensure deployment succeeds in staging first.

**Rationale:** Useful for deployment pipelines but requires deployment infrastructure.

---

#### 11. **Require code scanning results** 🔍
**Status:** ENABLE when CI/CD is set up

Choose which tools must provide code scanning results before the reference is updated.

**Configuration:**
- CodeQL (GitHub's native security scanning)
- Other SAST tools as needed

**Rationale:** Helps identify security vulnerabilities before merging code.

---

#### 12. **Require code quality results** 📊
**Status:** OPTIONAL

Choose which severity levels of code quality results should block pull request merges.

**Configuration:**
- Block on: High and Critical severity issues
- Allow: Medium and Low severity issues (as warnings)

**Rationale:** Maintains code quality standards but can be overly restrictive initially.

---

#### 13. **Manage static analysis tools in Copilot code review** 🔬
**Status:** ENABLE (if available)

Copilot code review will include findings from the selected static analysis tools in its review comments.

**Rationale:** Integrates static analysis findings into code review process.

---

## Recommended Configuration for FrostByte-API

For a typical API project, here's the recommended minimal ruleset:

### For `main` branch:

**MUST ENABLE:**
1. ✅ Require a pull request before merging
2. ✅ Require status checks to pass (once CI/CD is configured)
3. ✅ Block force pushes
4. ✅ Restrict deletions

**SHOULD ENABLE:**
5. ⚠️ Require linear history
6. 🔐 Require signed commits
7. 🤖 Automatically request Copilot code review

**CONFIGURE LATER:**
8. 🔍 Require code scanning results (after setting up CodeQL/security scanning)
9. 📊 Require code quality results (after setting up linters)

### For `develop` or `staging` branches (if used):

**MUST ENABLE:**
1. ✅ Require a pull request before merging
2. ✅ Block force pushes

**OPTIONAL:**
3. Require status checks to pass
4. Automatically request Copilot code review

---

## How to Apply These Rules

1. **Navigate to Repository Settings:**
   - Go to your repository on GitHub
   - Click "Settings" → "Rules" → "Rulesets"

2. **Create New Ruleset:**
   - Click "New ruleset" → "New branch ruleset"
   - Name it (e.g., "main-branch-protection")

3. **Set Target Branches:**
   - Use pattern: `main` or `master` (depending on your default branch)
   - Or use `**/*` for all branches with different rules

4. **Configure Rules:**
   - Enable the rules recommended above
   - Start with the "MUST ENABLE" rules first
   - Add "SHOULD ENABLE" rules as your workflow matures

5. **Set Bypass Permissions:**
   - Configure who can bypass these rules (typically: repository administrators)
   - Be conservative with bypass permissions

---

## Migration Path

If you're adding these rules to an existing repository:

1. **Phase 1:** Start with basic protection
   - Require pull requests
   - Block force pushes
   - Restrict deletions

2. **Phase 2:** Add automated checks
   - Require status checks (after CI/CD setup)
   - Add Copilot code review

3. **Phase 3:** Enhance security
   - Require signed commits
   - Add code scanning
   - Add code quality checks

4. **Phase 4:** Optimize workflow
   - Require linear history
   - Fine-tune status check requirements
   - Add deployment requirements

---

## Notes

- Branch protection rules only apply to branches that match the specified patterns
- Users with admin access can typically bypass some restrictions unless specifically configured
- Start with fewer rules and add more as your team's workflow stabilizes
- Communicate rule changes to your team before implementing them
- Review and adjust rules periodically based on team feedback

---

## Additional Resources

- [GitHub Branch Protection Documentation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets)
- [About Rulesets](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets)
