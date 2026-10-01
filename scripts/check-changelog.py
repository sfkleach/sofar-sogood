#!/usr/bin/env python3
"""Check that the top section of CHANGELOG.md is ready for release.

Usage: check-changelog.py [TAG]

Exits with a non-zero status if:
- The first ## section heading contains the word 'unreleased', or
- A version (vX.Y.Z) appears in more than one ## heading, or
- TAG is given (e.g. v1.2.3 or v1.2.3-rc.1) and the first ## heading does not
  contain it. A pre-release suffix is ignored, so v1.2.3-rc.1 matches v1.2.3.

The version of this project comes from the git tag, so unlike a project with a
version in a manifest file there is nothing else to compare against.
"""
import re
import sys


def main() -> None:
    if len(sys.argv) > 2:
        print(f"Usage: {sys.argv[0]} [TAG]", file=sys.stderr)
        sys.exit(1)

    with open("CHANGELOG.md") as f:
        content = f.read()

    # Collect all level-2 headings.
    headings = [line for line in content.splitlines() if line.startswith("## ")]
    if not headings:
        print("FAIL: No ## heading found in CHANGELOG.md.")
        sys.exit(1)

    # Extract version strings (vX.Y.Z) from all headings and check for duplicates.
    seen: dict[str, str] = {}
    for heading in headings:
        m = re.search(r"v\d+\.\d+\.\d+", heading)
        if not m:
            continue
        ver = m.group(0)
        if ver in seen:
            print(f"FAIL: Version '{ver}' appears more than once in CHANGELOG.md ('{seen[ver]}' and '{heading}').")
            sys.exit(1)
        seen[ver] = heading

    heading = headings[0]
    if "unreleased" in heading.lower():
        print(f"FAIL: Top CHANGELOG section is '{heading}' - release is not ready.")
        sys.exit(1)

    if len(sys.argv) == 2:
        tag = sys.argv[1]
        # Strip any pre-release suffix (e.g. "v1.2.3-rc.1" -> "v1.2.3").
        base = tag.split("-", 1)[0]
        if base not in heading:
            print(f"FAIL: Tag '{tag}' (version '{base}') not found in CHANGELOG heading '{heading}'.")
            sys.exit(1)

    print(f"OK: Top CHANGELOG section is '{heading}'.")


if __name__ == "__main__":
    main()
