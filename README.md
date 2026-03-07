⚠n This is work in pogress

# mulonda

> *Mulonda &mdash; from Tonga (Malawi), meaning "a guard."*

A small  tool that intercepts destructive shell commands and asks you before execution.


## The Problem

The terminal doesn't protect you from yourself.

```bash
rm -rf ./        # wrong directory. everything gone.
mv config.yml /etc/nginx/  # silently overwrites the original
cp -r src/ backup/         # clobbers your backup, no warning
```

No trash bin. No undo. Just gone.


## What Mulonda Does

Mulonda sits between you and your shell's most dangerous commands. Before anything destructive runs, it stops and asks:

```
mulonda: rm -rf ./dist
   Proceed? [y/N] 
```

It works via shell aliases &mdash; no daemon, no background process, no kernel magic. Just a fast Go binary that intercepts, prompts, then passes through to the real command if you confirm.

Safe commands and non-interactive contexts (scripts, CI pipelines) are passed through silently with zero overhead.


## Watchlist

Mulonda ships with somewhat sensible defaults (`rm`, `mv`, `cp` and more to be added soon). Users will be able to add or remove some rules a YAML watchlist:

```yaml
rules:
  - command: chmod
    args_match: ["777", "a+rwx"]
    warning: "Setting world-writable permissions"

  - command: rm
    flags_contain: ["-rf"]
    warning: "Recursive force delete"
```


## Usage

```bash
mulonda install          # inject aliases into your shell
mulonda add "chmod"      # add a command to your watchlist
mulonda remove "chmod"   # remove a command
mulonda list             # show active rules
mulonda uninstall        # remove all aliases
```

## Future Plans

- **v2 &mdash; eBPF mode:** intercept every command system-wide via kernel-level `execve()` hooking. No aliases needed. Works in any shell or context.
- Community watchlist presets (Docker, Kubernetes, database tooling)
- Dry-run mode &mdash; show what would be intercepted without prompting
- Audit log &mdash; keep a record of intercepted commands
- If it's cp or overwriting a file, show them diffs? or ask them if they meant to rewrite file `x` to `y` 


Please remember be reminded, this is still in development, you might sometimes not get the desired outcome.
