"""
Собирает backend/seeds/recipe_i18n_ru.json: description_ru и steps_ru для каждого рецепта.
Требует: pip install deep-translator

Запуск из корня репозитория:
  py bot/build_recipe_i18n.py

Уже переведённые записи в JSON пропускаются (можно продолжить после обрыва).
"""

from __future__ import annotations

import json
import pathlib
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    try:
        sys.stdout.reconfigure(encoding="utf-8")
    except Exception:
        pass

try:
    from deep_translator import GoogleTranslator
except ImportError:
    print("Установите: py -m pip install deep-translator", file=sys.stderr)
    sys.exit(1)

ROOT = pathlib.Path(__file__).resolve().parent.parent
RECIPES = ROOT / "backend" / "seeds" / "recipes.json"
OUT = ROOT / "backend" / "seeds" / "recipe_i18n_ru.json"
SLEEP = 0.2


def main() -> None:
    recipes = json.loads(RECIPES.read_text(encoding="utf-8"))
    existing: dict = {}
    if OUT.exists():
        try:
            existing = json.loads(OUT.read_text(encoding="utf-8"))
        except json.JSONDecodeError:
            existing = {}

    tr = GoogleTranslator(source="en", target="ru")
    n = len(recipes)
    for i, r in enumerate(recipes):
        title = r["title"].strip()
        if title in existing and existing[title].get("description_ru") and existing[title].get("steps_ru"):
            print(f"[skip] {i + 1}/{n} {title!r}")
            continue
        print(f"[tr] {i + 1}/{n} {title!r}")
        desc_ru = tr.translate(r.get("description") or "")
        time.sleep(SLEEP)
        steps_ru = []
        for s in r.get("steps") or []:
            steps_ru.append(tr.translate(s or ""))
            time.sleep(SLEEP)
        existing[title] = {"description_ru": desc_ru, "steps_ru": steps_ru}
        OUT.write_text(json.dumps(existing, ensure_ascii=False, indent=2), encoding="utf-8")

    OUT.write_text(json.dumps(existing, ensure_ascii=False, indent=2), encoding="utf-8")
    print("wrote", OUT)


if __name__ == "__main__":
    main()
