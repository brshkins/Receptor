"""Собирает seeds/recipe_stock_images.json из 003_recipe_stock_images.sql (title -> url)."""
import json
import pathlib
import re

root = pathlib.Path(__file__).resolve().parent
sql_path = root / "003_recipe_stock_images.sql"
out_path = root.parent / "seeds" / "recipe_stock_images.json"

text = sql_path.read_text(encoding="utf-8")
pat = re.compile(r"\(\s*'((?:[^']|'')*?)'\s*,\s*'(https:[^']+)'\s*\)")
out: dict[str, str] = {}
for m in pat.finditer(text):
    title = m.group(1).replace("''", "'")
    out[title] = m.group(2)

out_path.write_text(json.dumps(out, ensure_ascii=False, indent=2), encoding="utf-8")
print(f"wrote {len(out)} entries -> {out_path}")
