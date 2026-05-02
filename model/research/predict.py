from ultralytics import YOLO

model = YOLO("runs/detect/train2/weights/best.pt")

results = model("test.jpg", conf=0.1)  # ↓ снизили порог

for r in results:
    names = r.names

    products = []
    for cls_id, conf in zip(r.boxes.cls.tolist(), r.boxes.conf.tolist()):
        if conf < 0.1:
            continue

        product = names[int(cls_id)]

        if product not in products:
            products.append(product)

    print(products)