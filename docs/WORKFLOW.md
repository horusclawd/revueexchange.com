# Development Workflow

## Sprint Process

### Starting a Sprint
1. Create and switch to a new branch: `git checkout -b sprint-X`

### During Sprint
2. Commit each feature when complete with a descriptive message
3. Push branch regularly: `git push -u origin sprint-X`

### Completing a Sprint
3. **Code Review**: Compare UI API contract with the API, ensure they match
4. **Security Review**: Check for common vulnerabilities
5. **Fix Issues**: Address anything found in reviews
6. **Notify**: Tell the user the sprint is done
7. **Merge**: After user approval ("ok"), merge into main:

```bash
git checkout main
git merge sprint-X
git push origin main
```

## Notes
- No PRs for now - direct merge after approval
- Commit messages should describe what was done, not just "fix" or "update"
- Keep Terraform infrastructure up to date as features are added
