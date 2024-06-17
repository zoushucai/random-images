import os
import pandas as pd
from PIL import Image
import json

# 设置图片文件夹路径
IMAGE_BASE_FOLDER = 'images'

# 初始化列表来存储文件信息
data = []

# 遍历图片文件夹及其子文件夹
for subdir, dirs, files in os.walk(IMAGE_BASE_FOLDER):
    for file in files:
        filepath = os.path.join(subdir, file)
        try:
            with Image.open(filepath) as img:
                width, height = img.size
            # 获取子文件夹路径
            subfolder = os.path.relpath(subdir, IMAGE_BASE_FOLDER)
            # 获取文件后缀
            suffix = os.path.splitext(file)[1].lower()
            # 将信息添加到数据列表中
            data.append([subfolder, file, suffix, width, height])
        except Exception as e:
            print(f"无法处理文件 {filepath}: {e}")

# 将数据列表转换为 DataFrame
df = pd.DataFrame(data, columns=['sub', 'file', 'suf', 'width', 'height'])

# 打印 DataFrame
print(df)

# 将 DataFrame 转换为字典
data_dict = df.to_dict(orient='records')

# 将字典保存为 JSON 文件
with open('images_info.json', 'w', encoding='utf-8') as f:
    json.dump(data_dict, f, ensure_ascii=False, indent=4)
