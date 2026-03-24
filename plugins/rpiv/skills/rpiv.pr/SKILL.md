---
name: rpiv.pr
description: Generate comprehensive PR descriptions following repository templates
disable-model-invocation: true
---

# Generate PR Description

You are tasked with generating a comprehensive pull request description following the repository's
standard template.

**Tracker**: At the start of this workflow, update `docs/.tracker.json` — set the task status to
`review`.

## Steps to follow:

1. **Read the PR description template:**
   - Check if `.github/PULL_REQUEST_TEMPLATE.md` exists in the repository root
   - If not found, check `.github/PULL_REQUEST_TEMPLATE/` directory for template files
   - If no template exists, use a sensible default structure: Summary, Changes, How to verify, Notes
   - Read the template carefully to understand all sections and requirements

2. **Identify the PR to describe:**
   - Check if the current branch has an associated PR:
     `gh pr view --json url,number,title,state 2>/dev/null`
   - If no PR exists for the current branch, or if on main/master, list open PRs:
     `gh pr list --limit 10 --json number,title,headRefName,author`
   - Ask the user which PR they want to describe

3. **Gather comprehensive PR information:**
   - Get the full PR diff: `gh pr diff {number}`
   - If you get an error about no default remote repository, instruct the user to run
     `gh repo set-default` and select the appropriate repository
   - Get commit history: `gh pr view {number} --json commits`
   - Review the base branch: `gh pr view {number} --json baseRefName`
   - Get PR metadata: `gh pr view {number} --json url,title,number,state`

4. **Analyze the changes thoroughly:** (ultrathink about the code changes, their architectural
   implications, and potential impacts)
   - Read through the entire diff carefully
   - For context, read any files that are referenced but not shown in the diff
   - Understand the purpose and impact of each change
   - Identify user-facing changes vs internal implementation details
   - Look for breaking changes or migration requirements

5. **Handle verification requirements:**
   - Look for any checklist items in the "How to verify it" section of the template
   - For each verification step:
     - If it's a command you can run (like `make check test`, `npm test`, etc.), run it
     - If it passes, mark the checkbox as checked: `- [x]`
     - If it fails, keep it unchecked and note what failed: `- [ ]` with explanation
     - If it requires manual testing (UI interactions, external services), leave unchecked and note
       for user
   - Document any verification steps you couldn't complete

6. **Generate the description:**
   - If a template was found, fill out each section from the template thoroughly:
     - Answer each question/section based on your analysis
     - Be specific about problems solved and changes made
     - Focus on user impact where relevant
     - Include technical details in appropriate sections
     - Write a concise changelog entry if the template calls for one
   - If no template was found, use the default structure:
     - **Summary**: What this PR does and why
     - **Changes**: Bulleted list of significant changes
     - **How to verify**: Steps to test the changes
     - **Notes**: Breaking changes, migration steps, or other callouts
   - Ensure all checklist items are addressed (checked or explained)

7. **Present the description to the user:**
   - Show the user the generated description
   - Ask if they'd like any changes before updating the PR

8. **Update the PR:**
   - Write the description to a temporary file
   - Update the PR description directly:
     `gh pr edit {number} --body-file /tmp/pr_description_{number}.md`
   - Clean up the temporary file
   - Confirm the update was successful
   - If any verification steps remain unchecked, remind the user to complete them before merging

## Important

- Always discover the local PR template — this works across different repositories.
- Run verification commands when possible. Clearly note which steps need manual testing.
- Include breaking changes or migration notes prominently.
