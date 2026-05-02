"""Генерирует фрагмент SQL для 003: тот же FNV-1a + пулы, что в seeds/seed.go."""
import json
import os

FNV_OFFSET = 2166136261
FNV_PRIME = 16777619

def fnv1a32(s: str) -> int:
    h = FNV_OFFSET
    for b in s.encode("utf-8"):
        h ^= b
        h = (h * FNV_PRIME) & 0xFFFFFFFF
    return h

POOLS = {
    "pasta": [
        "https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1563379926898-05f4575a45d8?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1473093295043-cdd812d0e601?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1551183053-bf91a1d81141?w=1200&auto=format&fit=crop",
    ],
    "meat": [
        "https://images.unsplash.com/photo-1544025162-d76694265947?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1608039829572-08f78c7c5e63?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1529692236671-f94f979ee425?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1547592180-2a62e887fd72?w=1200&auto=format&fit=crop",
    ],
    "vegetarian": [
        "https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1540420773420-3366772f4999?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1490645935967-10de6ba17061?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=1200&auto=format&fit=crop",
    ],
    "breakfast": [
        "https://images.unsplash.com/photo-1533089860892-a7c6f0a88666?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1525351484163-7529144344d8?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1493770348161-369560ae357d?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1484723091739-30a097e8f929?w=1200&auto=format&fit=crop",
    ],
    "dessert": [
        "https://images.unsplash.com/photo-1563729784474-d77dbb933a9e?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1551782450-a2132b4ba21d?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1563805042-7684c019e1cb?w=1200&auto=format&fit=crop",
    ],
    "soup": [
        "https://images.unsplash.com/photo-1547592166-23ac45744acd?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1578247329974-62f1fbf53b5d?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1604908177522-8b0ee18b8338?w=1200&auto=format&fit=crop",
        "https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=1200&auto=format&fit=crop",
    ],
}


def url_for(title: str, category: str) -> str:
    cat = category.strip().lower()
    pool = POOLS.get(cat) or POOLS["vegetarian"]
    idx = fnv1a32(title.strip().lower()) % len(pool)
    return pool[idx]


def main() -> None:
    root = os.path.dirname(os.path.abspath(__file__))
    jpath = os.path.join(root, "..", "seeds", "recipes.json")
    with open(jpath, encoding="utf-8") as f:
        recipes = json.load(f)
    vals = []
    for r in recipes:
        t = r["title"].replace("'", "''")
        u = url_for(r["title"], r["category"]).replace("'", "''")
        vals.append(f"    ('{t}', '{u}')")
    body = ",\n".join(vals)
    sql = f"""-- Персональный image_url по английскому title (как в seeds/recipes.json).
-- Замените любой URL на свой. Перегенерация: py -3 gen_003_images.py > fragment.sql

UPDATE recipes AS r SET image_url = v.url
FROM (VALUES
{body}
) AS v(title, url)
WHERE r.title = v.title;
"""
    out = os.path.join(root, "003_recipe_stock_images.sql")
    with open(out, "w", encoding="utf-8") as f:
        f.write(sql)
    print(f"Wrote {out} ({len(recipes)} recipes)")


if __name__ == "__main__":
    main()
