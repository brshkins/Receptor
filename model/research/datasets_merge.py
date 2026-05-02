import shutil
from pathlib import Path
import yaml

FINAL_CLASSES = [
    "apple",
    "banana",
    "orange",
    "tomato",
    "cucumber",
    "potato",
    "onion",
    "carrot",
    "cabbage",
    "mushroom",
    "garlic",
    "pepper",
    "milk",
    "cheese",
    "yogurt",
    "butter",
    "eggs",
    "bread",
    "meat",
    "fish",
    "juice",
    "water",
    "instant_noodle"
]

FINAL_CLASS_TO_ID = {name: i for i, name in enumerate(FINAL_CLASSES)}

DATASETS = [
    {
        "path": "Food-Imgae---YOLO-1",
        "prefix": "foodimg",
        "mapping": {
            "Apple": "apple",
            "Apple -Green-": "apple",
            "Apple -Red-": "apple",
            "Banana": "banana",
            "Bread": "bread",
            "Butter": "butter",
            "Cheese": "cheese",
            "Cucumber": "cucumber",
            "Egg": "eggs",
            "Eggs": "eggs",
            "Fish": "fish",
            "Juice": "juice",
            "Meat": "meat",
            "Meat -Red-": "meat",
            "Milk": "milk",
            "Noodles": "instant_noodle",
            "Orange": "orange",
            "Oranges": "orange",
            "Potatoes -Package-": "potato",
            "Tomatoes": "tomato",
            "Water": "water",
            "Yogurt": "yogurt"
        }
    },
    {
        "path": "grocery-items-1",
        "prefix": "gitems",
        "mapping": {
            "apple": "apple",
            "instant_noodle": "instant_noodle",
            "juice": "juice",
            "orange": "orange"
        }
    },
    {
        "path": "Grocery-Store-1",
        "prefix": "gstore",
        "mapping": {
            "Carrot": "carrot",
            "Cucumber": "cucumber",
            "Floury potato": "potato",
            "Solid potato": "potato",
            "Sweet potato": "potato",
            "Onion": "onion",
            "Tomato": "tomato",
            "Tomato-": "tomato",
            "Cabbage": "cabbage",
            "Mushroom": "mushroom",
            "Garlic": "garlic",
            "Capsicum": "pepper"
        }
    }
]

OUTPUT_DIR = Path("dataset")


def load_names(dataset_path: Path) -> list[str]:
    with open(dataset_path / "data.yaml", "r", encoding="utf-8") as f:
        data = yaml.safe_load(f)

    names = data["names"]
    if isinstance(names, dict):
        return [names[i] for i in sorted(names.keys())]
    return names


def ensure_dirs():
    for split in ["train", "valid", "test"]:
        (OUTPUT_DIR / split / "images").mkdir(parents=True, exist_ok=True)
        (OUTPUT_DIR / split / "labels").mkdir(parents=True, exist_ok=True)


def process_split(dataset_path: Path, split: str, class_names: list[str], mapping: dict[str, str], prefix: str):
    images_dir = dataset_path / split / "images"
    labels_dir = dataset_path / split / "labels"

    if not images_dir.exists() or not labels_dir.exists():
        return

    for image_path in images_dir.iterdir():
        if image_path.suffix.lower() not in [".jpg", ".jpeg", ".png", ".bmp", ".webp"]:
            continue

        label_path = labels_dir / f"{image_path.stem}.txt"
        if not label_path.exists():
            continue

        converted_lines = []

        with open(label_path, "r", encoding="utf-8") as f:
            for line in f:
                parts = line.strip().split()
                if len(parts) < 5:
                    continue

                old_class_id = int(parts[0])
                bbox = parts[1:]

                if old_class_id >= len(class_names):
                    continue

                old_class_name = class_names[old_class_id]

                if old_class_name not in mapping:
                    continue

                final_class_name = mapping[old_class_name]
                new_class_id = FINAL_CLASS_TO_ID[final_class_name]

                converted_lines.append(f"{new_class_id} {' '.join(bbox)}")

        if not converted_lines:
            continue

        new_image_name = f"{prefix}_{image_path.name}"
        new_label_name = f"{prefix}_{image_path.stem}.txt"

        shutil.copy2(image_path, OUTPUT_DIR / split / "images" / new_image_name)

        with open(OUTPUT_DIR / split / "labels" / new_label_name, "w", encoding="utf-8") as f:
            f.write("\n".join(converted_lines) + "\n")


def write_yaml():
    data = {
        "path": str(OUTPUT_DIR.resolve()),
        "train": "train/images",
        "val": "valid/images",
        "test": "test/images",
        "names": {i: name for i, name in enumerate(FINAL_CLASSES)}
    }
    with open(OUTPUT_DIR / "data.yaml", "w", encoding="utf-8") as f:
        yaml.dump(data, f, allow_unicode=True, sort_keys=False)


def main():
    ensure_dirs()

    for ds in DATASETS:
        dataset_path = Path(ds["path"])
        class_names = load_names(dataset_path)

        for split in ["train", "valid", "test"]:
            process_split(
                dataset_path=dataset_path,
                split=split,
                class_names=class_names,
                mapping=ds["mapping"],
                prefix=ds["prefix"]
            )

    write_yaml()
    print("Готово: merged_dataset создан")


if __name__ == "__main__":
    main()