---
name: fix-ascii-diagram
description:
  Validate and fix ASCII box diagrams in markdown files. Use when the user asks to check, fix, or
  align ASCII art boxes/diagrams.
argument-hint: 'path/to/file.md'
---

# ASCII Box Diagram Spec

You are a specialist for validating and fixing pure-ASCII box diagrams. When invoked, read the
target file (passed as `$ARGUMENTS`, or ask the user), find all ASCII diagrams inside fenced code
blocks, and validate/fix them against the spec below.

## Terminology

- **Box**: a rectangular region bounded by border characters (`+`, `-`, `|`)
- **Border**: the perimeter characters of a box
- **Corner**: `+` character at every intersection of horizontal and vertical borders
- **Horizontal border**: contiguous `-` characters between two corners on the same row
- **Vertical border**: contiguous `|` characters between two corners on the same column
- **Content**: text inside a box, excluding borders and padding
- **Padding**: exactly 1 space between content and each vertical border

## Rules

1. **Corners** — Every point where a horizontal border meets a vertical border must be a `+`.
2. **Horizontal contiguity** — Every column between two corners on the same row must be `-` (or `+`
   if another border intersects).
3. **Vertical contiguity** — Every row between two corners on the same column must be `|` (or `+` if
   another border intersects).
4. **Padding** — There must be exactly 1 space between content and each vertical border (`|`). That
   is, the character immediately after a left `|` and immediately before a right `|` must be a
   space.
5. **Box width** — Determined by the longest content line in that box:
   `width = max(len(content)) + 4` (1 border + 1 pad + content + 1 pad + 1 border). Shorter content
   lines are right-padded with spaces to fill the box.
6. **Box height** — Determined by the number of content rows plus 2 (top border + bottom border). If
   the box is subdivided by horizontal rules, each `+---...---+` row acts as both a bottom and top
   border for adjacent sections.

## Validation Procedure

For each diagram found:

1. **Locate corners**: find all `+` characters and record their `(row, col)` positions.
2. **Identify boxes**: group corners into rectangles — a box exists when four corners form an
   axis-aligned rectangle and the edges between them are valid borders.
3. **Check horizontal borders**: for each pair of corners on the same row that form a box edge,
   verify every column between them is `-` or `+`.
4. **Check vertical borders**: for each pair of corners on the same column that form a box edge,
   verify every row between them is `|` or `+`.
5. **Check padding**: for each content row inside a box, verify `line[left_border_col + 1] == ' '`
   and `line[right_border_col - 1] == ' '`.
6. **Check width consistency**: measure the longest content string inside the box, verify box inner
   width equals `max(len(content)) + 2` (1 pad each side).
7. **Check all lines within a box have equal length** between borders (right-padded with spaces).

## Fixing

When fixing, apply these transformations:~

1. **Determine correct width** for each box from its longest content line.
2. **Rewrite borders** — rebuild `+` corners and `-` fills to the correct width.
3. **Rewrite content lines** — left-pad with 1 space, right-pad with spaces to fill, close with 1
   space + `|`.
4. **Preserve non-box text** between boxes (labels, arrows, annotations) — adjust spacing to
   maintain visual alignment with the resized boxes.

## Programmatic Check (use with Bash)

To audit column alignment, run:

```bash
awk 'NR>=START && NR<=END {
  printf "Line %2d (len %3d): ", NR, length($0)
  for (i=1; i<=length($0); i++) {
    c = substr($0, i, 1)
    if (c == "|" || c == "+") printf "%d(%s) ", i-1, c
  }
  print ""
}' FILE
```

Verify that each border column contains `|` or `+` on **every** row between its topmost and
bottommost corner. If any row is missing a border character at an expected column, the diagram is
broken.

## Output

After fixing, re-run the programmatic check and show the user the before/after column alignment to
confirm correctness.
